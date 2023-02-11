package pipe

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"syscall"
	"text/template"

	"github.com/google/uuid"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func Tasks(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tasks", "parent").
		SetJobWrapper(func(job Job) Job {
			return tl.JobSequence(
				Setup(tl).Job(),

				tl.JobSequence(
					ReadTemplates(tl).Job(),
					tl.JobParallel(
						GenerateNginxConfigurationTemplate(tl).Job(),
						GenerateTemplates(tl).Job(),
					),
				),
			)
		})
}

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		ShouldRunBefore(func(t *Task[Pipe]) error {
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
			t.CreateSubtask("nginx").Set(func(t *Task[Pipe]) error {
				template, err := Templates.ReadFile("templates/nginx.conf.go.tmpl")

				if err != nil {
					return err
				}

				t.Pipe.Ctx.Templates.Nginx = string(template)

				return nil
			}).
				AddSelfToTheParentAsParallel()

			t.CreateSubtask("server").Set(func(t *Task[Pipe]) error {
				template, err := Templates.ReadFile("templates/server.conf.go.tmpl")

				if err != nil {
					return err
				}

				t.Pipe.Ctx.Templates.Server = string(template)

				return nil
			}).
				AddSelfToTheParentAsParallel()

			t.CreateSubtask("upstream").Set(func(t *Task[Pipe]) error {
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

func GenerateNginxConfigurationTemplate(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("generate", "nginx").
		Set(func(t *Task[Pipe]) error {
			t.Log.Infof("Creating then Nginx configuration template.")

			tmpl, err := template.New("nginx.conf").Parse(t.Pipe.Ctx.Templates.Nginx)

			if err != nil {
				return err
			}

			output := new(bytes.Buffer)

			var rLimit syscall.Rlimit

			if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
				return err
			}

			cores := runtime.NumCPU()

			if err = tmpl.Execute(output, NginxTemplate{
				CpuCores:          cores,
				RLimit:            rLimit.Cur,
				WorkerConnections: rLimit.Cur / uint64(cores),
			}); err != nil {
				return err
			}

			t.Log.Debugf(
				"Nginx configuration template:\n%s",
				output.String(),
			)

			p := path.Join(
				NGINX_ROOT_CONFIGURATION_FOLDER,
				NGINX_CONFIGURATION,
			)

			t.Log.Debugln(
				"Writing Nginx configuration file.",
			)

			if err := os.WriteFile(p, output.Bytes(), 0600); err != nil {
				return err
			}

			return nil
		})
}

func GenerateTemplates(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("generate").
		Set(func(t *Task[Pipe]) error {
			for i, v := range t.Pipe.Nginx.Configuration {
				func(_ int, conf ConfigurationJson) {
					t.CreateSubtask(conf.Server.Listen).
						Set(func(t *Task[Pipe]) error {
							id := uuid.New().String()

							// stream template
							t.CreateSubtask("server").
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

									if err := os.WriteFile(p, output.Bytes(), 0600); err != nil {
										return err
									}

									return nil
								}).
								AddSelfToTheParentAsParallel()

								// upstream template
							t.CreateSubtask("upstream").
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

									if err := os.WriteFile(p, output.Bytes(), 0600); err != nil {
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
