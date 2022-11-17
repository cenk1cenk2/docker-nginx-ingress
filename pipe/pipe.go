package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Nginx struct {
		Configuration
	}

	Pipe struct {
		Ctx

		Nginx
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				Tasks(tl).Job(),
				Services(tl).Job(),
			)
		})
}
