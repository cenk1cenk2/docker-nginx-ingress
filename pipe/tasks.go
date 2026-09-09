package pipe

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"syscall"
	"text/template"
	"uuid"

	. "github.com/cenk1cenk2/plumber/v7"
)

func Tasks(tl *TaskList) *Task {
	return tl.CreateTask("tasks", "parent").
		SetJobWrapper(func(_ Job, _ *Task) Job {
			return JobSequence(
				Setup(tl).Job(),

				JobSequence(
					ReadTemplates(tl).Job(),
					JobParallel(
						GenerateNginxConfigurationTemplate(tl).Job(),
						GenerateTemplates(tl).Job(),
					),
				),
			)
		})
}

func Setup(tl *TaskList) *Task {
	return tl.CreateTask("init").
		Set(func(_ context.Context, t *Task) error {
			C.Directories.ServerConfiguration = path.Join(
				NGINX_ROOT_CONFIGURATION_FOLDER,
				TEMPLATE_FOLDER_SERVERS,
			)

			if err := os.RemoveAll(C.Directories.ServerConfiguration); err != nil {
				return err
			}

			if err := os.MkdirAll(C.Directories.ServerConfiguration, os.ModePerm); err != nil {
				return err
			}

			C.Directories.UpstreamConfiguration = path.Join(
				NGINX_ROOT_CONFIGURATION_FOLDER,
				TEMPLATE_FOLDER_UPSTREAMS,
			)

			if err := os.RemoveAll(C.Directories.UpstreamConfiguration); err != nil {
				return err
			}

			return os.MkdirAll(C.Directories.UpstreamConfiguration, os.ModePerm)
		})
}

func ReadTemplates(tl *TaskList) *Task {
	return tl.CreateTask("template").
		Set(func(_ context.Context, t *Task) error {
			t.CreateSubtask("nginx").Set(func(_ context.Context, t *Task) error {
				template, err := Templates.ReadFile("templates/nginx.conf.go.tmpl")

				if err != nil {
					return err
				}

				C.Templates.Nginx = string(template)

				return nil
			}).
				AddSelfToTheParentAsParallel()

			t.CreateSubtask("server").Set(func(_ context.Context, t *Task) error {
				template, err := Templates.ReadFile("templates/server.conf.go.tmpl")

				if err != nil {
					return err
				}

				C.Templates.Server = string(template)

				return nil
			}).
				AddSelfToTheParentAsParallel()

			t.CreateSubtask("upstream").Set(func(_ context.Context, t *Task) error {
				template, err := Templates.ReadFile("templates/upstream.conf.go.tmpl")

				if err != nil {
					return err
				}

				C.Templates.Upstream = string(template)

				return nil
			}).
				AddSelfToTheParentAsParallel()

			return nil
		}).ShouldRunAfter(func(ctx context.Context, t *Task) error {
		return t.RunSubtasks(ctx)
	})
}

func GenerateNginxConfigurationTemplate(tl *TaskList) *Task {
	return tl.CreateTask("generate", "nginx").
		Set(func(_ context.Context, t *Task) error {
			t.Log.Info("Creating then Nginx configuration template.")

			tmpl, err := template.New("nginx.conf").Parse(C.Templates.Nginx)

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

			t.Log.Debug(
				fmt.Sprintf(
					"Nginx configuration template:\n%s",
					output.String(),
				),
			)

			p := path.Join(
				NGINX_ROOT_CONFIGURATION_FOLDER,
				NGINX_CONFIGURATION,
			)

			t.Log.Debug(
				"Writing Nginx configuration file.",
			)

			return os.WriteFile(p, output.Bytes(), 0600)
		})
}

func GenerateTemplates(tl *TaskList) *Task {
	return tl.CreateTask("generate").
		Set(func(_ context.Context, t *Task) error {
			for _, conf := range P.Nginx.Configuration {
				t.CreateSubtask(conf.Server.Listen).
					Set(func(_ context.Context, t *Task) error {
						id := uuid.New().String()

						// stream template
						t.CreateSubtask("server").
							Set(func(_ context.Context, t *Task) error {
								t.Log.Debug(fmt.Sprintf("Stream %s will have the id: %s", conf.Server.Listen, id))

								t.Log.Info(fmt.Sprintf("Creating server template for: %s", conf.Server.Listen))

								tmpl, err := template.New("server.conf").Parse(C.Templates.Server)

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

								t.Log.Debug(
									fmt.Sprintf(
										"Server template for %s:\n%s",
										conf.Server.Listen,
										output.String(),
									),
								)

								p := path.Join(
									C.Directories.ServerConfiguration,
									fmt.Sprintf("%s.conf", id),
								)

								t.Log.Debug(
									fmt.Sprintf(
										"Writing service file for %s: %s",
										conf.Server.Listen,
										p,
									),
								)

								return os.WriteFile(p, output.Bytes(), 0600)
							}).
							AddSelfToTheParentAsParallel()

							// upstream template
						t.CreateSubtask("upstream").
							Set(func(_ context.Context, t *Task) error {
								t.Log.Info(fmt.Sprintf("Creating upstream template for: %s", conf.Server.Listen))

								tmpl, err := template.New("upstream.conf").Parse(C.Templates.Upstream)

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

								t.Log.Debug(
									fmt.Sprintf(
										"Upstream template for %s:\n%s",
										conf.Server.Listen,
										output.String(),
									),
								)

								p := path.Join(
									C.Directories.UpstreamConfiguration,
									fmt.Sprintf("%s.conf", id),
								)

								t.Log.Debug(
									fmt.Sprintf(
										"Writing upstream file for %s: %s",
										conf.Server.Listen,
										p,
									),
								)

								return os.WriteFile(p, output.Bytes(), 0600)
							}).
							AddSelfToTheParentAsParallel()

						return nil
					}).
					AddSelfToTheParentAsParallel().
					ShouldRunAfter(func(ctx context.Context, t *Task) error {
						return t.RunSubtasks(ctx)
					})
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunSubtasks(ctx)
		})
}
