package pipe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/google/uuid"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		ShouldRunBefore(func(t *Task[Pipe]) error {
			if err := json.Unmarshal([]byte(t.Pipe.Nginx.Configuration), &t.Pipe.Ctx.NginxConfiguration); err != nil {
				return fmt.Errorf("Can not decode configuration: %s", err)
			}

			if err := tl.Validate(&t.Pipe.Ctx); err != nil {
				return err
			}

			t.Pipe.Ctx.Directories.ServerConfiguration = path.Join(
				NGINX_ROOT_CONFIGURATION_FOLDER,
				TEMPLATE_FOLDER_SERVERS,
			)

			if err := os.RemoveAll(t.Pipe.Ctx.Directories.ServerConfiguration); err != nil {
				return err
			}

			if err := os.MkdirAll(t.Pipe.Ctx.Directories.ServerConfiguration, os.ModePerm); err != nil {
				return err
			}

			t.Pipe.Ctx.Directories.UpstreamConfiguration = path.Join(
				NGINX_ROOT_CONFIGURATION_FOLDER,
				TEMPLATE_FOLDER_UPSTREAMS,
			)

			if err := os.RemoveAll(t.Pipe.Ctx.Directories.UpstreamConfiguration); err != nil {
				return err
			}

			if err := os.MkdirAll(t.Pipe.Ctx.Directories.UpstreamConfiguration, os.ModePerm); err != nil {
				return err
			}

			return nil
		}).
		Set(func(t *Task[Pipe]) error {

			return nil
		})
}

func ReadTemplates(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("template").
		Set(func(t *Task[Pipe]) error {
			t.CreateSubtask("template:server").Set(func(t *Task[Pipe]) error {
				template, err := Templates.ReadFile("templates/server.conf.go.tmpl")

				if err != nil {
					return err
				}

				t.Pipe.Ctx.Templates.Server = string(template)

				return nil
			}).
				AddSelfToTheParentAsParallel()

			t.CreateSubtask("template:upstream").Set(func(t *Task[Pipe]) error {
				template, err := Templates.ReadFile("templates/upstream.conf.go.tmpl")

				if err != nil {
					return err
				}

				t.Pipe.Ctx.Templates.Upstream = string(template)

				return nil
			}).
				AddSelfToTheParentAsParallel()

			return nil
		}).ShouldRunAfter(func(t *Task[Pipe]) error {
		return t.RunSubtasks()
	})
}

func GenerateTemplates(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("generate").
		Set(func(t *Task[Pipe]) error {
			for i, v := range t.Pipe.Ctx.NginxConfiguration {
				func(i int, conf ConfigurationJson) {
					t.CreateSubtask(fmt.Sprintf("generate:%s", conf.Server.Listen)).
						Set(func(t *Task[Pipe]) error {
							id := uuid.New().String()

							// stream template
							t.CreateSubtask(fmt.Sprintf("generate:%s:server", conf.Server.Listen)).
								Set(func(t *Task[Pipe]) error {
									t.Log.Debugf("Stream %s will have the id: %s", conf.Server.Listen, id)

									t.Log.Infof("Creating server template for: %s", conf.Server.Listen)

									tmpl, err := template.New("server.conf").Parse(t.Pipe.Ctx.Templates.Server)

									if err != nil {
										return err
									}

									output := new(bytes.Buffer)

									if err = tmpl.Execute(output, ServerTemplate{
										Listen:   conf.Server.Listen,
										Upstream: id,
										Options:  conf.Server.Options,
									}); err != nil {
										return err
									}

									t.Log.Debugf(
										"Server template for %s:\n%s",
										conf.Server.Listen,
										output.String(),
									)

									p := path.Join(
										t.Pipe.Ctx.Directories.ServerConfiguration,
										fmt.Sprintf("%s.conf", id),
									)

									t.Log.Debugf(
										"Writing service file for %s: %s",
										conf.Server.Listen,
										p,
									)

									if err := os.WriteFile(p, output.Bytes(), 0644); err != nil {
										return err
									}

									return nil
								}).
								AddSelfToTheParentAsParallel()

								// upstream template
							t.CreateSubtask(fmt.Sprintf("generate:%s:upstream", conf.Server.Listen)).
								Set(func(t *Task[Pipe]) error {
									t.Log.Infof("Creating upstream template for: %s", conf.Server.Listen)

									tmpl, err := template.New("upstream.conf").Parse(t.Pipe.Ctx.Templates.Upstream)

									if err != nil {
										return err
									}

									output := new(bytes.Buffer)

									if err := tmpl.Execute(output, UpstreamTemplate{
										Upstream: id,
										Servers:  conf.Upstream.Servers,
										Options:  conf.Upstream.Options,
									}); err != nil {
										return err
									}

									t.Log.Debugf(
										"Upstream template for %s:\n%s",
										conf.Server.Listen,
										output.String(),
									)

									p := path.Join(
										t.Pipe.Ctx.Directories.UpstreamConfiguration,
										fmt.Sprintf("%s.conf", id),
									)

									t.Log.Debugln(
										fmt.Sprintf(
											"Writing upstream file for %s: %s",
											conf.Server.Listen,
											p,
										),
									)

									if err := os.WriteFile(p, output.Bytes(), 0644); err != nil {
										return err
									}

									return nil
								}).
								AddSelfToTheParentAsParallel()

							return nil
						}).
						AddSelfToTheParentAsParallel().
						ShouldRunAfter(func(t *Task[Pipe]) error {
							return t.RunSubtasks()
						})
				}(i, v)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}
