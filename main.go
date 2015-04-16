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

	potSaneLimitForThreeBet = 18
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

		switch decision {

		case "CHECK":
			table.Check()
		case "FOLD":
			table.Fold()

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

		case "RESTEAL/ALL-IN\n3-BET/FOLD if raiser >= EP":
			threeBetFold(table)
		case "RESTEAL/ALL-IN\n3-BET/FOLD if raiser >= MP":
			threeBetFold(table)
		case "RESTEAL/ALL-IN\n3-BET/FOLD if raiser >= LATER":
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

		case "RESTEAL/ALL-IN\n3-BET/ALL-IN":
			threeBetAllIn(table)
		case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= EP":
			threeBetAllIn(table)
		case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= MP":
			threeBetAllIn(table)
		case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= LATER":
			threeBetAllIn(table)

		case "FLOP BET/ALL-IN":
			flopBetAllIn(table)

		case "TURN BET/ALL-IN":
			turnBetAllIn(table)

		}
	}

	fmt.Println()

	if args["-v"].(bool) != false {
		fmt.Println(strategy.Messages)
		fmt.Println(table)
		fmt.Println(table.Image.Path)
	}

	fmt.Println(decision)
}

func performTwoActions(firstAction func(), secondAction func(), flag string) {
	if !flagFileIsOk(flag) {
		createFlagFile(flag)
		firstAction()
	} else {
		secondAction()
	}
}

func raiseFold(table Table) {
	performTwoActions(
		table.Raise, table.Fold,
		fmt.Sprintf("/tmp/croc-fold-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func raiseAllIn(table Table) {
	performTwoActions(
		table.Raise, table.AllIn,
		fmt.Sprintf("/tmp/croc-allin-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func stealFold(table Table) {
	performTwoActions(
		table.Steal, table.Fold,
		fmt.Sprintf("/tmp/croc-fold-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func stealAllIn(table Table) {
	performTwoActions(
		table.Steal, table.AllIn,
		fmt.Sprintf("/tmp/croc-allin-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func threeBetFold(table Table) {
	flag := fmt.Sprintf("/tmp/croc-fold-%s-%s", table.Hero.Hand, table.Window.Id)

	if !flagFileIsOk(flag) && table.Pot <= potSaneLimitForThreeBet {
		createFlagFile(flag)
		table.ThreeBet()
	} else {
		table.Fold()
	}
}

func threeBetAllIn(table Table) {
	performTwoActions(
		table.ThreeBet, table.AllIn,
		fmt.Sprintf("/tmp/croc-allin-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func flopBetAllIn(table Table) {
	performTwoActions(
		table.Bet, table.AllIn,
		fmt.Sprintf("/tmp/croc-flop-allin-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func turnBetAllIn(table Table) {
	performTwoActions(
		table.Bet, table.AllIn,
		fmt.Sprintf("/tmp/croc-turn-allin-%s-%s", table.Hero.Hand, table.Window.Id),
	)
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
