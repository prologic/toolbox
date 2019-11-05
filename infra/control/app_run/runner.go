package app_run

import (
	"flag"
	"github.com/GeertJohan/go.rice"
	"github.com/pkg/profile"
	"github.com/watermint/toolbox/catalogue"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/control/app_control_impl"
	"github.com/watermint/toolbox/infra/control/app_opt"
	"github.com/watermint/toolbox/infra/control/app_root"
	"github.com/watermint/toolbox/infra/control/app_run_impl"
	"github.com/watermint/toolbox/infra/network/app_diag"
	"github.com/watermint/toolbox/infra/network/app_network"
	"github.com/watermint/toolbox/infra/quality/qt_control_impl"
	"github.com/watermint/toolbox/infra/recpie/app_kitchen"
	"github.com/watermint/toolbox/infra/recpie/app_recipe_group"
	"github.com/watermint/toolbox/infra/recpie/app_vo_impl"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"github.com/watermint/toolbox/infra/ui/app_ui"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

func Run(args []string, bx, web *rice.Box) (found bool) {
	// Initialize resources
	mc := app_run_impl.NewContainer(bx)
	ui := app_ui.NewConsole(mc, qt_control_impl.NewMessageMock(), false)
	cat := catalogue.Catalogue()

	// Select recipe or group
	cmd, grp, rcp, rem, err := cat.Select(args)

	switch {
	case err != nil:
		//if grp != nil {
		//	grp.PrintUsage(ui)
		//} else {
		//	cat.PrintUsage(ui)
		//}
		//os.Exit(app_control.FailureInvalidCommand)
		return false

	case rcp == nil:
		grp.PrintUsage(ui)
		os.Exit(app_control.Success)
	}

	// Initialize recipe value object
	cmdPath := make([]string, 0)
	cmdPath = append(cmdPath, grp.Path...)
	if cmd != "" {
		cmdPath = append(cmdPath, cmd)
	}
	recipeName := strings.Join(cmdPath, " ")

	vo := rcp.Requirement()
	f := flag.NewFlagSet(recipeName, flag.ContinueOnError)
	com := app_opt.NewDefaultCommonOpts()

	cvc := app_vo_impl.NewValueContainer(com)
	cvc.MakeFlagSet(f, ui)

	vc := app_vo_impl.NewValueContainer(vo)
	vc.MakeFlagSet(f, ui)

	err = f.Parse(rem)
	rem2 := f.Args()
	if err != nil || (len(rem2) > 0 && rem2[0] == "help") {
		grp.PrintRecipeUsage(ui, rcp, f)
		os.Exit(app_control.FailureInvalidCommandFlags)
	}
	vc.Apply(vo)
	cvc.Apply(com)

	// Apply common flags

	// - Quiet
	if com.Quiet {
		ui = app_ui.NewQuiet(mc)
	}

	// Up
	so := make([]app_control.UpOpt, 0)
	if com.Workspace != "" {
		so = append(so, app_control.WorkspacePath(com.Workspace))
	}
	if com.Debug {
		so = append(so, app_control.Debug())
	}
	if com.Secure {
		so = append(so, app_control.Secure())
	}
	so = append(so, app_control.Concurrency(com.Concurrency))
	so = append(so, app_control.RecipeName(recipeName))

	ctl := app_control_impl.NewSingle(ui, bx, web, mc, com.Quiet, catalogue.Recipes())
	err = ctl.Up(so...)
	if err != nil {
		os.Exit(app_control.FatalStartup)
	}
	defer ctl.Down()

	// - Quiet
	if qui, ok := ui.(*app_ui.Quiet); ok {
		qui.SetLogger(ctl.Log())
	}

	// Recover
	defer func() {
		err := recover()
		if err != nil {
			l := ctl.Log()
			l.Debug("Recovery from panic")
			l.Error(ctl.UI().Text("run.error.panic"),
				zap.Any("error", err),
			)
			l.Error(ctl.UI().Text("run.error.panic.instruction"),
				zap.String("JobPath", ctl.Workspace().Job()),
			)

			for depth := 0; ; depth++ {
				_, file, line, ok := runtime.Caller(depth)
				if !ok {
					break
				}
				ctl.Log().Debug("Trace",
					zap.Int("Depth", depth),
					zap.String("File", file),
					zap.Int("Line", line),
				)
			}
			ctl.Abort(app_control.Reason(app_control.FatalPanic))
		}
	}()

	// Trap signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT)
	go func() {
		for s := range sig {
			// fatal shutdown
			ctl.Log().Debug("Signal", zap.Any("signal", s))
			pc := make([]uintptr, 16)
			n := runtime.Callers(0, pc)
			pc = pc[:n]
			frames := runtime.CallersFrames(pc)
			for f := 0; ; f++ {
				frame, more := frames.Next()
				ctl.Log().Debug("Frame",
					zap.Int("Frame", f),
					zap.String("File", frame.File),
					zap.Int("Line", frame.Line),
					zap.String("Function", frame.Function),
				)
				if !more {
					break
				}
			}
			ui := ctl.UI()
			ui.Error("run.error.interrupted")
			ui.Error("run.error.interrupted.instruction", app_msg.P{"JobPath": ctl.Workspace().Job()})
			ctl.Abort(app_control.Reason(app_control.FatalInterrupted))

			// in case the controller didn't fire exit..
			os.Exit(app_control.FatalInterrupted)
		}
	}()

	// - Proxy config
	app_network.SetHttpProxy(com.Proxy, ctl)

	// App Header
	app_recipe_group.AppHeader(ui)

	// Diagnosis
	err = app_diag.Runtime(ctl)
	if err != nil {
		ctl.Abort(app_control.Reason(app_control.FatalRuntime))
	}
	if ctl.IsProduction() {
		err = app_diag.Network(ctl)
		if err != nil {
			ctl.Abort(app_control.Reason(app_control.FatalNetwork))
		}
	}

	// Apply profiler
	if com.Debug {
		defer profile.Start(profile.MemProfile).Stop()
	}

	// Run
	ctl.Log().Debug("Run recipe", zap.Any("vo", vo), zap.Any("common", com))
	k := app_kitchen.NewKitchen(ctl, vo)
	err = rcp.Exec(k)
	if err != nil {
		ctl.Log().Error("Recipe failed with an error", zap.Error(err))
		ui.Failure("run.error.recipe.failed", app_msg.P{"Error": err.Error()})
		os.Exit(app_control.FailureGeneral)
	}
	app_root.FlushSuccessShutdownHook()

	return true
}
