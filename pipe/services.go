package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func RunNginx(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("nginx").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"nginx",
				"-g",
				"daemon off;",
			).
				EnableTerminator().
				SetOnTerminator(func(c *Command[Pipe]) error {
					t.Plumber.SendTerminated()

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
