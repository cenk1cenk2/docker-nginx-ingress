package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v3"
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
		SetTasks(
			TL.JobSequence(
				Setup(&TL).Job(),

				TL.JobSequence(
					ReadTemplates(&TL).Job(),
					GenerateTemplates(&TL).Job(),
				),

				RunNginx(&TL).Job(),
			),
		)
}
