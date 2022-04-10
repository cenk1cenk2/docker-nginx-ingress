package main

import (
	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
)

type Ctx struct {
}

var Context Ctx

func VerifyVariables() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "verify"},
		Task: func(t *utils.Task) error {
			err := utils.ValidateAndSetDefaults(t.Metadata, &Pipe)

			if err != nil {
				return err
			}

			return nil
		},
	}
}
