package app_feature

import (
	"encoding/json"
	"github.com/watermint/toolbox/essentials/go/es_reflect"
	"github.com/watermint/toolbox/essentials/log/es_log"
	"github.com/watermint/toolbox/infra/app"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"os/user"
	"time"
)

type Feature interface {
	IsProduction() bool
	IsDebug() bool
	IsTest() bool
	IsTestWithMock() bool
	IsQuiet() bool
	IsSecure() bool
	IsLowMemory() bool
	IsAutoOpen() bool

	// UI format
	UIFormat() string

	// Concurrency configuration.
	Concurrency() int

	// Toolbox home path. Returns empty if a user doesn't specify the path.
	Home() string

	// Budget for memory usage
	BudgetMemory() string

	// Budget for storage usage
	BudgetStorage() string

	// Retrieve feature
	OptInGet(oi OptIn) (f OptIn, found bool)

	// Update opt-in feature
	OptInUpdate(oi OptIn) error

	// With test mode
	AsTest(useMock bool) Feature

	// With quiet mode, but this will not guarantee UI/log are converted into quiet mode.
	AsQuiet() Feature

	// Console log level
	ConsoleLogLevel() es_log.Level
}

type OptIn interface {
	// The timestamp of opt-in, in ISO8601 format.
	// Empty when the user is not yet agreed.
	OptInTimestamp() string

	// Name of the user who opt'ed in.
	OptInUser() string

	// True when this feature enabled.
	OptInIsEnabled() bool

	// Name of the feature.
	OptInName(v OptIn) string

	// Opt-in
	OptInCommit(enable bool) OptIn
}

func OptInFrom(v map[string]interface{}, oi OptIn) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, oi)
}

type OptInStatus struct {
	// The timestamp of opt-in, in ISO8601 format.
	Timestamp string `json:"timestamp"`

	// Name of the user who opt'ed in
	User string `json:"user"`

	// Opt-in status.
	Status bool `json:"status"`
}

func (z OptInStatus) OptInCommit(enable bool) OptIn {
	usr, _ := user.Current()

	switch {
	case usr.Name != "":
		z.User = usr.Name
	case usr.Username != "":
		z.User = usr.Username
	default:
		z.User = "unknown"
	}
	z.Status = enable
	z.Timestamp = time.Now().Format(time.RFC3339)
	return &z
}

func (z OptInStatus) OptInName(v OptIn) string {
	return es_reflect.Key(app.Pkg, v)
}

func (z OptInStatus) OptInTimestamp() string {
	return z.Timestamp
}

func (z OptInStatus) OptInUser() string {
	return z.User
}

func (z OptInStatus) OptInIsEnabled() bool {
	return z.Status
}

func OptInAgreement(v OptIn) app_msg.Message {
	return app_msg.ObjMessage(v, "agreement")
}

func OptInDisclaimer(v OptIn) app_msg.Message {
	return app_msg.ObjMessage(v, "disclaimer")
}

func OptInDescription(v OptIn) app_msg.Message {
	return app_msg.ObjMessage(v, "desc")
}
