package worker

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-apache-test/lib/backend"
	"github.com/go-apache-test/lib/config"
	"github.com/go-apache-test/lib/cst"
	"github.com/go-apache-test/lib/filesystem"
	"github.com/go-apache-test/lib/logger"
	"github.com/go-apache-test/lib/report"
	"github.com/go-apache-test/lib/utils"
	"github.com/go-apache-test/lib/vhost"
)

type worker struct {
	filer  filesystem.Filer
	config config.Config
	vhosts []vhost.Vhost
	logger logger.Logger
}

func New(filer filesystem.Filer, config config.Config) worker {
	return worker{
		filer:  filer,
		config: config,
		vhosts: []vhost.Vhost{},
		logger: logger.NewLog(config.Debug.LogLevel),
	}
}

func (w worker) Init() error {
	w.logger.Debug("Init")
	w.logger.Debug(w.config.String())

	w.logger.Debug("Load vhosts")
	if err := w.loadVhosts(); err != nil {
		return err
	}

	w.logger.Debug("Update apache2.conf by adding ServerName.\n")
	if err := w.updateHttpdConf(); err != nil {
		return err
	}

	w.logger.Debug("Generate port.conf file.\n")
	if err := w.generatePortConf(); err != nil {
		return err
	}

	w.logger.Debug("Update /etc/hosts with alias.\n")
	if err := w.updateHosts(); err != nil {
		return err
	}

	return nil
}

func (w worker) Run() error {
	w.logger.Debug("Start the backend")
	srv := backend.Start(w.config.Httpd.BackendsPort)

	w.logger.Debug("Run the tests")
	err := w.RunTests()

	w.logger.Debug("Stop the backend")
	backend.Stop(srv)

	return err
}

func (w *worker) loadVhosts() error {
	vhostFiles, err := w.filer.List(cst.DockerVhostsDir, func(filepath string) bool {
		return strings.HasSuffix(filepath, ".conf")
	})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, vhostPath := range vhostFiles {
		v := vhost.New(vhostPath)

		wg.Add(1)
		go func(v vhost.Vhost) {
			defer wg.Done()

			w.logger.Debug(fmt.Sprintf("Get alias: %s", v.OriPath))
			if err := v.GetAlias(w.filer); err != nil {
				w.logger.Error(fmt.Sprintf("%v", err))
				return
			}

			w.logger.Debug(fmt.Sprintf("Get suites: %s", v.OriPath))
			if err := v.GetSuites(w.filer); err != nil {
				w.logger.Error(fmt.Sprintf("%v", err))
				return
			}

			w.logger.Debug(fmt.Sprintf("Copy vhost: %s", v.OriPath))
			if err := w.filer.Copy(v.OriPath, v.NewPath); err != nil {
				w.logger.Error(fmt.Sprintf("%v", err))
				return
			}

			if w.config.Test.InsideDocker {
				filePath := filepath.Base(v.NewPath)
				w.logger.Debug(fmt.Sprintf("a2ensite -m %s", filePath))
				_, stderr, err := utils.Exec("a2ensite", "-m", filePath)
				if err != nil {
					w.logger.Error(stderr)
					return
				}
			}

			w.vhosts = append(w.vhosts, v)
			w.logger.Debug(fmt.Sprintf("Vhost loaded: %s", v.OriPath))
		}(v)
	}
	wg.Wait()

	return nil
}

func (w *worker) RunTests() error {
	vhostFiles, err := w.filer.List(cst.DockerVhostsDir, func(filepath string) bool {
		return strings.HasSuffix(filepath, ".conf")
	})
	if err != nil {
		return err
	}

	report := report.New()

	var wg sync.WaitGroup
	for _, vhostPath := range vhostFiles {
		v := vhost.New(vhostPath)

		wg.Add(1)
		go func(v vhost.Vhost) {
			defer wg.Done()

			if err := v.GetSuites(w.filer); err != nil {
				w.logger.Error(fmt.Sprintf("%v", err))
				return
			}

			for _, s := range v.Suites {
				for _, r := range s.Rules {
					r.Init(v.Expectations, s.Expectations)
					w.logger.Debug(fmt.Sprintf("--> Request %s", r.Name))
					result, err := r.Request()
					w.logger.Debug(fmt.Sprintf("<-- Request %s", r.Name))
					if err != nil {
						failure := report.Add(v.OriPath, s.Name, r.Name, err)
						w.logger.Info(r.StrError(err))
						w.logger.Info(fmt.Sprint(failure))
						return
					}

					w.logger.Debug(fmt.Sprintf("Validate %s", r.Name))
					err = r.ShouldBe(result)
					if err != nil {
						failure := report.Add(v.OriPath, s.Name, r.Name, err)
						w.logger.Info(r.StrError(err))
						w.logger.Info(fmt.Sprint(failure))
						return
					}

					report.Add(v.OriPath, s.Name, r.Name, err)
					w.logger.Info(r.StrSuccess())
				}
			}

		}(v)
	}
	wg.Wait()

	w.logger.Info(fmt.Sprint(report))

	return nil
}

func (w worker) updateHttpdConf() error {
	return w.filer.Append(cst.HttpdConfFile, fmt.Sprintf("\nServerName %s\n", cst.HttpdHost))
}

func (w worker) generatePortConf() error {
	return w.filer.Write(cst.HttpdPortFile, fmt.Sprintf("Listen %s\n", cst.HttpdPort))
}

func (w worker) updateHosts() error {
	domains := []string{}
	domains = append(domains, w.config.Httpd.Backends...)
	for _, v := range w.vhosts {
		domains = append(domains, v.Alias...)
	}

	hosts := fmt.Sprintf("\n127.0.0.1 %s\n", strings.Join(domains, " "))
	w.logger.Debug(hosts)

	return w.filer.Write(cst.DockerHostsFile, hosts)
}

func (w worker) String() string {
	return fmt.Sprintf(`
	config	%v
	vhost	%v
	`, w.config, w.vhosts)
}
