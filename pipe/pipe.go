package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Nginx struct {
		Configuration string
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
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				Tasks(tl).Job(),
				Services(tl).Job(),
			)
		})
}
