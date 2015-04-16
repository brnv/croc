package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

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
	croc [<filepath>] [--wid=<window_id>] [-v]`
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

	table := Table{}

	if args["<filepath>"] != nil {
		table.Image.Path = args["<filepath>"].(string)
	} else {
		if args["--wid"] != nil {
			table.Window.Id = args["--wid"].(string)
		} else {
			table.Window.InitId()
		}

		table.Image.Path = table.Window.Screenshot()
	}

	if !table.Validate() {
		log.Fatal(table.Errors)
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}

	wg.Add(5)

	go func() {
		table.HandRecognize()
		wg.Done()
	}()

	go func() {
		table.ButtonRecognize()
		table.HeroPositionRecognize()
		wg.Done()
	}()

	go func() {
		table.LimpersRecognize()
		wg.Done()
	}()

	go func() {
		table.PotRecognize()
		wg.Done()
	}()

	go func() {
		//table.HeroChipsRecognize()
		wg.Done()
	}()

	wg.Wait()

	strategy := Strategy{}

	strategy.Table = table
	decision := strategy.Run()

	if !table.FoldButtonIsVisible() && decision == "FOLD" {
		os.Exit(1)
	}

	if !table.HeroMoveInProgress() && decision != "FOLD" {
		os.Exit(1)
	}

	if table.HeroMoveInProgress() {
		table.BoardRecognize()

		strategy = Strategy{}

		strategy.Table = table
		decision = strategy.Run()
	}

	if table.Window.Id != "" {
		table.PerformAutomatedActions(decision)
	}

	if args["-v"].(bool) != false {
		fmt.Println(strategy.Messages)
		fmt.Println(table)
		fmt.Println(table.Image.Path)
	}

	fmt.Println(decision)
}
