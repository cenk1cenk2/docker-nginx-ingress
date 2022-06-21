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
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
