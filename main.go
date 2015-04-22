package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/docopt/docopt-go"
	"github.com/op/go-logging"
	"github.com/thearkit/runcmd"
)

var (
	log       = logging.MustGetLogger("croc")
	cmdRunner *runcmd.Local
)

const (
	usage = `
	Usage:
	croc [<filepath>] [--wid=<window_id>] [-v] [-a]`
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	logging.SetLevel(logging.NOTICE, "")

	args, err := docopt.Parse(usage, nil, true, "croc", false)
	if err != nil {
		log.Fatal(err)
	}

	cmdRunner, err = runcmd.NewLocalRunner()
	if err != nil {
		log.Fatal(err)
	}

	table := Table{
		BigBlindSize: 2,
	}

	if args["<filepath>"] != nil {
		table.Image.Path = args["<filepath>"].(string)
	} else {
		if args["--wid"] != nil {
			table.Window.Id = args["--wid"].(string)
		} else {
			err = table.Window.ManualSelect()
			if err != nil {
				log.Fatal(err)
			}
		}

		table.Image.Path, err = table.Window.TakeScreenshot()
		if err != nil {
			log.Fatal(err)
		}
	}

	if !table.Validate() {
		log.Fatal(table.Errors)
		os.Exit(1)
	}

	table.Recognize()

	strategy := Strategy{
		Table: table,
	}

	err = strategy.Run()
	if err != nil {
		log.Fatal(err)
	}

	if !table.HeroMoveInProgress() {
		if !table.FoldButtonIsVisible() && strategy.Decision == "FOLD" {
			os.Exit(1)
		}

		if strategy.Decision != "FOLD" {
			os.Exit(1)
		}
	}

	if args["-a"].(bool) != false && table.Window.Id != "" {
		table.PerformAutomatedActions(strategy.Decision)
	}

	if args["-v"].(bool) != false {
		fmt.Println(table.Image.Path)
		fmt.Println(table)
		fmt.Println(strategy.Messages)
	}

	fmt.Println(strategy.Decision)
}
