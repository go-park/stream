package routine

import (
	"github.com/go-park/stream/internal/helper"
	"github.com/go-park/stream/support/function"
)

func Run(runner function.Runner) {
	helper.RequireCanButNonNil(runner)
	go runner.Run()
}

func RunArg[T any](arg T, runner func(T)) {
	helper.RequireCanButNonNil(runner)
	go runner(arg)
}
