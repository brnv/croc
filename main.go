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

	if !table.Check() {
		log.Fatal(table.Errors)
		os.Exit(1)
	}

	strategy := Strategy{
		Table: &table,
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

	if table.FastFoldButtonIsVisible() {
		decision := strategy.Run()

		if decision == "FOLD" {
			fmt.Print("\n")

			if args["-v"].(bool) != false {
				fmt.Println(strategy.Messages)
				fmt.Println(table)
			}

			table.ClickFold()

			fmt.Println("FAST FOLD")

			os.Exit(0)
		}

		os.Exit(1)
	} else if !table.HeroMoveIsPending() {
		os.Exit(1)
	}

	table.BoardRecognize()

	decision := strategy.Run()

	//@TODO: remember decision for current window id
	//@TODO: implement tasks logic for current situation
	// to make program perform 2 steps decisions by itself

	if table.Window.Id != "" {
		switch decision {
		case "CHECK":
			table.ClickCheck()
		case "FOLD":
			table.ClickFold()

		case "RAISE/FOLD":
			table.ClickRaise()
		case "RAISE/ALL-IN":
			table.ClickRaise()

		case "STEAL/FOLD":
			table.ClickSteal()
		case "STEAL/ALL-IN":
			table.ClickSteal()

		case "3-BET/FOLD if raiser >= EP":
			table.ClickThreeBet()
		case "3-BET/ALL-IN if raiser >= EP":
			table.ClickThreeBet()

		case "3-BET/FOLD if raiser >= MP":
			table.ClickThreeBet()
		case "3-BET/ALL-IN if raiser >= MP":
			table.ClickThreeBet()

		case "3-BET/FOLD if raiser >= LATER":
			table.ClickThreeBet()
		case "3-BET/ALL-IN if raiser >= LATER":
			table.ClickThreeBet()

		case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= EP":
			table.ClickThreeBet()
		case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= MP":
			table.ClickThreeBet()
		case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= LATER":
			table.ClickThreeBet()

		case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= EP":
			table.ClickThreeBet()
		case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= MP":
			table.ClickThreeBet()
		case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= LATER":
			table.ClickThreeBet()
		}

	}

	fmt.Print("\n")

	if args["-v"].(bool) != false {
		fmt.Println(strategy.Messages)
		fmt.Println(table)
	}

	fmt.Println(decision)
}
