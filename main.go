package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

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

	if table.Window.Id != "" {
		switch decision {

		case "CHECK":
			table.ClickCheck()

		case "FOLD":
			table.ClickFold()

		case "RAISE/FOLD":
			raiseFold(table)

		case "RAISE/ALL-IN":
			table.ClickRaise()

		case "STEAL/FOLD":
			stealFold(table)

		case "STEAL/ALL-IN":
			table.ClickSteal()

		case "3-BET/FOLD if raiser >= EP":
			threeBetFold(table)

		case "3-BET/ALL-IN if raiser >= EP":
			table.ClickThreeBet()

		case "3-BET/FOLD if raiser >= MP":
			threeBetFold(table)

		case "3-BET/ALL-IN if raiser >= MP":
			table.ClickThreeBet()

		case "3-BET/FOLD if raiser >= LATER":
			threeBetFold(table)

		case "3-BET/ALL-IN if raiser >= LATER":
			table.ClickThreeBet()

		case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= EP":
			threeBetFold(table)

		case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= MP":
			threeBetFold(table)

		case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= LATER":
			threeBetFold(table)

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

func raiseFold(table Table) {
	flag := getTimeFlagname(
		fmt.Sprintf("/tmp/croc-raise-%s-%s", table.Hero.Hand, table.Window.Id),
	)

	if _, err := os.Stat(flag); os.IsNotExist(err) {
		createFlagFile(flag)
		table.ClickRaise()
	} else {
		table.ClickFold()
	}
}

func stealFold(table Table) {
	flag := getTimeFlagname(
		fmt.Sprintf("/tmp/croc-raise-%s-%s", table.Hero.Hand, table.Window.Id),
	)

	if _, err := os.Stat(flag); os.IsNotExist(err) {
		createFlagFile(flag)
		table.ClickSteal()
	} else {
		table.ClickFold()
	}
}

func threeBetFold(table Table) {
	flag := getTimeFlagname(
		fmt.Sprintf("/tmp/croc-raise-%s-%s", table.Hero.Hand, table.Window.Id),
	)

	if _, err := os.Stat(flag); os.IsNotExist(err) {
		createFlagFile(flag)
		table.ClickThreeBet()
	} else {
		table.ClickFold()
	}
}

func createFlagFile(name string) {
	os.Create(name)
}

func getTimeFlagname(basename string) string {
	timeNow := time.Now().Unix()
	return fmt.Sprintf("%s-%d", basename, timeNow-timeNow%60)
}
