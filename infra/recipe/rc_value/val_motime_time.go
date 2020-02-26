package rc_value

import (
	"github.com/watermint/toolbox/domain/model/mo_essential"
	"github.com/watermint/toolbox/domain/model/mo_time"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"reflect"
)

func newValueMoTimeTime(name string) rc_recipe.Value {
	v := &ValueMoTimeTime{name: name}
	v.time = mo_time.Zero()
	v.isOptional = false
	return v
}

func newValueMoTimeTimeOptional(name string) rc_recipe.Value {
	v := &ValueMoTimeTime{name: name}
	v.time = mo_time.Zero()
	v.isOptional = true
	return v
}

type ValueMoTimeTime struct {
	name       string
	dateTime   string
	time       mo_time.Time
	isOptional bool
}

func (z *ValueMoTimeTime) ValueText() string {
	return z.dateTime
}

func (z *ValueMoTimeTime) Accept(t reflect.Type, v0 interface{}, name string) rc_recipe.Value {
	if t.Implements(reflect.TypeOf((*mo_time.TimeOptional)(nil)).Elem()) {
		return newValueMoTimeTimeOptional(name)
	}
	if t.Implements(reflect.TypeOf((*mo_time.Time)(nil)).Elem()) {
		return newValueMoTimeTime(name)
	}
	return nil
}

func (z *ValueMoTimeTime) Bind() interface{} {
	return &z.dateTime
}

func (z *ValueMoTimeTime) Init() (v interface{}) {
	return z.time
}

func (z *ValueMoTimeTime) ApplyPreset(v0 interface{}) {
	z.time = v0.(mo_time.Time)
	if !z.time.IsZero() {
		z.dateTime = z.time.String()
	}
}

func (z *ValueMoTimeTime) Apply() (v interface{}) {
	return z.time
}

func (z *ValueMoTimeTime) Debug() interface{} {
	return map[string]string{
		"dateTime": z.dateTime,
		"time":     z.time.String(),
	}
}

func (z *ValueMoTimeTime) SpinUp(ctl app_control.Control) (err error) {
	// argument was't given, but applied on preset or custom value
	if z.dateTime == "" && !z.time.IsZero() {
		return nil
	}

	// optional case
	if z.dateTime == "" && z.isOptional {
		// mark as unset
		if to, ok := z.time.(mo_essential.OptionalMutable); ok {
			to.Unset()
		}
		return nil
	}

	ti := z.time.(*mo_time.TimeImpl)
	if err = ti.UpdateTime(z.dateTime); err != nil {
		return err
	}
	return nil
}

func (z *ValueMoTimeTime) SpinDown(ctl app_control.Control) error {
	return nil
}
