package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func Services(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("services", "parent").
		SetJobWrapper(func(job Job) Job {
			return tl.JobParallel(
				RunNginx(tl).Job(),
			)
		})
}

func RunNginx(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("nginx").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"nginx",
				"-g",
				"daemon off;",
			).
				EnsureIsAlive().
				EnableTerminator().
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
