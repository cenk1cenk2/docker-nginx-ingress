package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"

	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
	"github.com/flosch/pongo2/v5"
	"github.com/google/uuid"
)

type Ctx struct {
	NginxConfiguration Configuration
	Templates          struct {
		Server   string
		Upstream string
	}
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

			err = json.Unmarshal([]byte(Pipe.Nginx.Configuration), &Context.NginxConfiguration)

			if err != nil {
				t.Log.Fatalln(fmt.Sprintf("Can not decode configuration: %s", err))
			}

			var wg sync.WaitGroup
			wg.Add(len(Context.NginxConfiguration))
			errs := []error{}

			for i, v := range Context.NginxConfiguration {
				go func(i int, v ConfigurationJson) {
					defer wg.Done()

					err := utils.ValidateAndSetDefaults(t.Metadata, &Context.NginxConfiguration[i])

					if err != nil {
						errs = append(errs, err)
					}
				}(i, v)
			}

			wg.Wait()

			if len(errs) > 0 {
				for _, v := range errs {
					t.Log.Errorln(v)
				}

				return errors.New("Errors encountered while validation.")
			}

			return nil
		},
	}
}

func ReadTemplates() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "template"},
		Task: func(t *utils.Task) error {
			template, err := Templates.ReadFile("templates/server.conf.j2")

			if err != nil {
				return err
			}

			Context.Templates.Server = string(template)

			template, err = Templates.ReadFile("templates/upstream.conf.j2")

			if err != nil {
				return err
			}

			Context.Templates.Upstream = string(template)

			return nil
		}}
}

func GenerateTemplates() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "generate"},
		Task: func(t *utils.Task) error {
			var wg sync.WaitGroup
			wg.Add(len(Context.NginxConfiguration))
			errs := []error{}

			for i, v := range Context.NginxConfiguration {
				go func(i int, conf ConfigurationJson) {
					defer wg.Done()

					id := uuid.New().String()

					t.Log.Debugln(
						fmt.Sprintf(
							"Stream %s will have the id: %s", conf.Server.Listen, id))

					t.Log.Infoln(
						fmt.Sprintf("Creating server template for: %s", conf.Server.Listen),
					)

					tpl, err := pongo2.FromString(Context.Templates.Server)

					if err != nil {
						errs = append(errs, err)

						return
					}

					output, err := tpl.Execute(
						pongo2.Context{
							"listen":   conf.Server.Listen,
							"upstream": id,
							"options":  conf.Server.Options,
						},
					)

					if err != nil {
						errs = append(errs, err)

						return
					}

					t.Log.Debugln(
						fmt.Sprintf("Server template for %s:\n%s", conf.Server.Listen, output),
					)

					p := path.Join(
						NGINX_ROOT_CONFIGURATION_FOLDER,
						TEMPLATE_FOLDER_SERVERS,
						fmt.Sprintf("%s.conf", id),
					)

					t.Log.Debugln(
						fmt.Sprintf(
							"Trying to generate service file for %s: %s",
							conf.Server.Listen,
							p,
						),
					)

					err = os.WriteFile(p, []byte(output), 0644)

					if err != nil {
						errs = append(errs, err)

						return
					}

					t.Log.Infoln(
						fmt.Sprintf("Creating upstream template for: %s", conf.Server.Listen),
					)

					tpl, err = pongo2.FromString(Context.Templates.Upstream)

					if err != nil {
						errs = append(errs, err)

						return
					}

					output, err = tpl.Execute(
						pongo2.Context{
							"upstream": id,
							"servers":  conf.Upstream.Servers,
							"options":  conf.Upstream.Options,
						},
					)

					if err != nil {
						errs = append(errs, err)

						return
					}

					t.Log.Debugln(
						fmt.Sprintf("Upstream template for %s:\n%s", conf.Server.Listen, output),
					)

					p = path.Join(
						NGINX_ROOT_CONFIGURATION_FOLDER,
						TEMPLATE_FOLDER_UPSTREAMS,
						fmt.Sprintf("%s.conf", id),
					)

					t.Log.Debugln(
						fmt.Sprintf(
							"Trying to generate upstream file for %s: %s",
							conf.Server.Listen,
							p,
						),
					)

					err = os.WriteFile(p, []byte(output), 0644)

					if err != nil {
						errs = append(errs, err)

						return
					}

				}(i, v)
			}

			wg.Wait()

			if len(errs) > 0 {
				for _, v := range errs {
					t.Log.Errorln(v)
				}

				return errors.New("Errors encountered while generating templates.")
			}

			return nil
		},
	}
}

func StartNginx() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "nginx"},
		Task: func(t *utils.Task) error {

			cmd := exec.Command("nginx")

			cmd.Args = append(cmd.Args, "-g", "daemon off;")

			t.Command = cmd

			return nil
		},
	}
}
