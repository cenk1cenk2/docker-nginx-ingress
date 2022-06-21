package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type (
	Nginx struct {
		Configuration string
	}

	Terminator struct {
		ShouldTerminate chan bool
		Terminated      chan bool
	}

	Pipe struct {
		Ctx
		Terminator

		Nginx
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			tl.Pipe.Terminator.ShouldTerminate = make(chan bool, 1)
			tl.Pipe.Terminator.Terminated = make(chan bool, 1)

			return nil
		}).
		SetTasks(
			TL.JobParallel(
				Terminate(&TL).Job(),

				TL.JobSequence(
					Setup(&TL).Job(),

					TL.JobSequence(
						ReadTemplates(&TL).Job(),
						GenerateTemplates(&TL).Job(),
					),

					RunNginx(&TL).Job(),
				),
			),
		)
}
