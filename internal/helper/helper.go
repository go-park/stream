package helper

import "reflect"

type Empty struct{}

func RequireCanButNonNil(arg any) {
	if arg == nil {
		panic("unexpected nil pointer")
	}
}

func RequireNonNil(arg any) {
	RequireNonNilMsg(arg, "unexpected nil pointer")
}

func RequireNonNilMsg(arg any, msg string) {
	canNil, isNil := IsNil(arg)
	if canNil && isNil {
		if len(msg) == 0 {
			msg = "unexpected nil pointer"
		}
		panic(msg)
	}
}

func IsNil(arg any) (canNil, isNil bool) {
	rv := reflect.ValueOf(arg)
	defer func() {
		if err := recover(); err == nil {
			canNil = true
		}
	}()
	isNil = rv.IsNil()
	return
}
