package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
	"text/template"

	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
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
			template, err := Templates.ReadFile("templates/server.conf.go.tmpl")

			if err != nil {
				return err
			}

			Context.Templates.Server = string(template)

			template, err = Templates.ReadFile("templates/upstream.conf.go.tmpl")

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

					tmpl, err := template.New("server.conf").
						Parse(Context.Templates.Server)

					if err != nil {
						errs = append(errs, err)

						return
					}

					output := new(bytes.Buffer)

					err = tmpl.Execute(output, ServerTemplate{
						Listen:   conf.Server.Listen,
						Upstream: id,
						Options:  conf.Server.Options,
					})

					if err != nil {
						errs = append(errs, err)

						return
					}

					t.Log.Debugln(
						fmt.Sprintf(
							"Server template for %s:\n%s",
							conf.Server.Listen,
							output.String(),
						),
					)

					p := path.Join(
						NGINX_ROOT_CONFIGURATION_FOLDER,
						TEMPLATE_FOLDER_SERVERS,
						fmt.Sprintf("%s.conf", id),
					)

					t.Log.Debugln(
						fmt.Sprintf(
							"Writing service file for %s: %s",
							conf.Server.Listen,
							p,
						),
					)

					err = os.WriteFile(p, output.Bytes(), 0644)

					if err != nil {
						errs = append(errs, err)

						return
					}

					t.Log.Infoln(
						fmt.Sprintf("Creating upstream template for: %s", conf.Server.Listen),
					)

					tmpl, err = template.New("upstream.conf").
						Parse(Context.Templates.Upstream)

					if err != nil {
						errs = append(errs, err)

						return
					}

					output = new(bytes.Buffer)

					err = tmpl.Execute(output, UpstreamTemplate{
						Upstream: id,
						Servers:  conf.Upstream.Servers,
						Options:  conf.Upstream.Options,
					})

					if err != nil {
						errs = append(errs, err)

						return
					}

					t.Log.Debugln(
						fmt.Sprintf(
							"Upstream template for %s:\n%s",
							conf.Server.Listen,
							output.String(),
						),
					)

					p = path.Join(
						NGINX_ROOT_CONFIGURATION_FOLDER,
						TEMPLATE_FOLDER_UPSTREAMS,
						fmt.Sprintf("%s.conf", id),
					)

					t.Log.Debugln(
						fmt.Sprintf(
							"Writing upstream file for %s: %s",
							conf.Server.Listen,
							p,
						),
					)

					err = os.WriteFile(p, output.Bytes(), 0644)

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
