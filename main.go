package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/go-apache-test/lib/config"
	"github.com/go-apache-test/lib/cst"
	"github.com/go-apache-test/lib/filesystem"
	"github.com/go-apache-test/lib/worker"
)

func main() {
	fmt.Println("- Go Load Test -")

	action := flag.String("x", "test", "Action")
	flag.Parse()

	err := execute(*action)
	handlerErr(err)
}

func execute(action string) error {
	filer := filesystem.New()
	config, err := config.New(cst.DockerConfigFile, filer)
	config.Test.InsideDocker = action != "test"
	handlerErr(err)
	worker := worker.New(filer, config)

	switch action {
	case "init":
		err = worker.Init()
		handlerErr(err)
	case "run":
		worker.Run()
	case "test":
		err = worker.Init()
		handlerErr(err)
		worker.Run()
	default:
		fmt.Println(action)
	}

	return nil
}

func handlerErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
