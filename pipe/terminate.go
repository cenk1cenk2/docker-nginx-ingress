package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func TerminatePredicate(tl *TaskList[Pipe]) JobPredicate {
	return tl.Predicate(func(tl *TaskList[Pipe]) bool {
		tl.Log.Debugln("Registered terminate listener.")

		a := <-tl.Pipe.Terminator.ShouldTerminate

		tl.Log.Warnln("Running termination tasks...")

		return a
	})
}

func Terminate(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate").
		SetJobWrapper(func(job Job) Job {
			return tl.JobBackground(
				tl.JobIf(
					TerminatePredicate(tl), tl.GuardAlways(job),
				),
			)
		}).
		Set(func(t *Task[Pipe]) error {
			t.SetSubtask(
				tl.JobParallel(
					TerminateNginx(tl).Job(),
				),
			)

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			if err := t.RunSubtasks(); err != nil {
				return err
			}

			t.Log.Infoln("Graceful termination finished.")

			t.Pipe.Terminator.Terminated <- true

			close(t.Pipe.Terminator.ShouldTerminate)
			close(t.Pipe.Terminator.Terminated)

			return nil
		})
}

func TerminateNginx(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("terminate:nginx").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand("pkill", "nginx").
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
