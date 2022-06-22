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
				Set(func(c *Command[Pipe]) error {
					go func() {
						signal := <-t.Plumber.Terminator.ShouldTerminate

						c.Log.Debugf("Forwarding signal to process: %s", signal)

						if err := c.Command.Process.Signal(signal); err != nil {
							t.SendError(err)
						}

						t.Plumber.SendTerminated()
					}()

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
