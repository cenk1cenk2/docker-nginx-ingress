package pipe

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
)

func Services(tl *TaskList) *Task {
	return tl.CreateTask("services", "parent").
		SetJobWrapper(func(_ Job, _ *Task) Job {
			return JobParallel(
				RunNginx(tl).Job(),
			)
		})
}

func RunNginx(tl *TaskList) *Task {
	return tl.CreateTask("nginx").
		Set(func(_ context.Context, t *Task) error {
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
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
