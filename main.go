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

	if !table.Validate() {
		log.Fatal(table.Errors)
		os.Exit(1)
	}

	strategy := Strategy{}

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

	strategy.Table = table

	decision := strategy.Run()

	if !table.FoldButtonIsVisible() && decision == "FOLD" {
		if !table.FastFoldToAnyBetIsChecked() {
			table.FastFoldToAnyBet()
		} else {
			os.Exit(1)
		}
	}

	if !table.HeroMoveInProgress() && decision != "FOLD" {
		os.Exit(1)
	}

	if table.HeroMoveInProgress() && decision != "FOLD" {
		table.BoardRecognize()

		strategy.Table = table

		decision = strategy.Run()
	}

	if table.Window.Id != "" {

		switch decision {

		case "CHECK":
			table.Check()

		case "FOLD":
			if !table.FastFoldToAnyBetIsChecked() {
				table.Fold()
			}

		case "RAISE/FOLD":
			raiseFold(table)
		case "STEAL/FOLD":
			stealFold(table)
		case "3-BET/FOLD if raiser >= EP":
			threeBetFold(table)
		case "3-BET/FOLD if raiser >= MP":
			threeBetFold(table)
		case "3-BET/FOLD if raiser >= LATER":
			threeBetFold(table)
		case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= EP":
			threeBetFold(table)
		case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= MP":
			threeBetFold(table)
		case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= LATER":
			threeBetFold(table)

		case "RAISE/ALL-IN":
			raiseAllIn(table)
		case "STEAL/ALL-IN":
			stealAllIn(table)
		case "3-BET/ALL-IN":
			threeBetAllIn(table)
		case "3-BET/ALL-IN if raiser >= EP":
			threeBetAllIn(table)
		case "3-BET/ALL-IN if raiser >= MP":
			threeBetAllIn(table)
		case "3-BET/ALL-IN if raiser >= LATER":
			threeBetAllIn(table)
		case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= EP":
			threeBetAllIn(table)
		case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= MP":
			threeBetAllIn(table)
		case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= LATER":
			threeBetAllIn(table)

		}
	}

	fmt.Println()

	if args["-v"].(bool) != false {
		fmt.Println(strategy.Messages)
		fmt.Println(table)
	}

	fmt.Println(decision)
}

func raiseFold(table Table) {
	table.Raise()

	if !table.FastFoldToAnyBetIsChecked() {
		table.FastFoldToAnyBet()
	}
}

func raiseAllIn(table Table) {
	flag := fmt.Sprintf("/tmp/croc-allin-%s-%s", table.Hero.Hand, table.Window.Id)

	if !flagFileIsOk(flag) {
		createFlagFile(flag)
		table.Raise()
	} else {
		table.AllIn()
	}
}

func stealFold(table Table) {
	table.Steal()

	if !table.FastFoldToAnyBetIsChecked() {
		table.FastFoldToAnyBet()
	}
}

func stealAllIn(table Table) {
	flag := fmt.Sprintf("/tmp/croc-allin-%s-%s", table.Hero.Hand, table.Window.Id)

	if !flagFileIsOk(flag) {
		createFlagFile(flag)
		table.Steal()
	} else {
		table.AllIn()
	}
}

func threeBetFold(table Table) {
	table.ThreeBet()

	if !table.FastFoldToAnyBetIsChecked() {
		table.FastFoldToAnyBet()
	}
}

func threeBetAllIn(table Table) {
	flag := fmt.Sprintf("/tmp/croc-allin-%s-%s", table.Hero.Hand, table.Window.Id)

	if !flagFileIsOk(flag) {
		createFlagFile(flag)
		table.ThreeBet()
	} else {
		table.AllIn()
	}
}

func createFlagFile(name string) {
	os.Create(name)
}

func flagFileIsOk(flag string) bool {
	file, err := os.Stat(flag)

	if os.IsNotExist(err) {
		return false
	}

	if file.ModTime().Unix() < time.Now().Unix()-60 {
		return false
	}

	return true
}
