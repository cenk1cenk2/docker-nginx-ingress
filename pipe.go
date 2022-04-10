package main

import (
	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
)

type (
	Nginx struct {
		Configuration string
	}

	Plugin struct {
		Nginx Nginx
	}
)

var Pipe Plugin = Plugin{}

func (p Plugin) Exec() error {
	utils.AddTasks(
		[]utils.Task{VerifyVariables(), ReadTemplates(), GenerateTemplates(), StartNginx()},
	)

	utils.RunAllTasks(utils.DefaultRunAllTasksOptions)

	return nil
}
