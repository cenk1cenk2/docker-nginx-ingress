package pipe

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
)

type (
	Nginx struct {
		Configuration
	}

	Pipe struct {
		Nginx
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				Tasks(tl).Job(),
				Services(tl).Job(),
			)
		})
}
