package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/docopt/docopt-go"
	"github.com/op/go-logging"
	"github.com/thearkit/runcmd"
)

var (
	log                 = logging.MustGetLogger("croc")
	reWindowId          = regexp.MustCompile("window id # of group leader: (0x.*)")
	reMatchedSample     = regexp.MustCompile("[[:word:]]+/(.*)\\..*$")
	reCompareErrorLevel = regexp.MustCompile("\\((.*)\\)$")
	compareThreshold    = 0.05
	cmdRunner           *runcmd.Local
	samplesFiles        map[string][]string
)

type Table struct {
	Screenshot string

	HandFirstCard  string
	HandSecondCard string

	Pot string
}

const usage = `
	Usage:
	croc
	croc <filepath>
`

func main() {
	var err error

	cmdRunner, err = runcmd.NewLocalRunner()
	if err != nil {
		log.Fatal(err)
	}

	// find out table type: six or nine players in
	table := Table{}

	args, _ := docopt.Parse(usage, nil, true, "croc", false)
	if args["<filepath>"] != nil {
		table.Screenshot = args["<filepath>"].(string)
	} else {
		table.Screenshot, err = makeTableScreenshot()
		if err != nil {
			log.Fatal(err)
		}
	}

	wg := &sync.WaitGroup{}

	wg.Add(2)
	go func() {
		table.HandRecognition()
		wg.Done()
	}()

	go func() {
		table.PotRecognition()
		wg.Done()
	}()

	wg.Wait()

	log.Notice("%v", table)
}

//@TODO: detect pot digit size
// refactor duplicated code
func (table *Table) PotRecognition() {
	// for 4 digit pot
	digitOne := ""
	digitTwo := ""
	digitThree := ""
	digitFour := ""

	wg := &sync.WaitGroup{}

	wg.Add(4)

	go func() {
		digitOne = recognize(
			table.Crop(9, 13, 397, 154),
			"pot_samples/*",
		)
		wg.Done()
	}()

	go func() {
		digitTwo = recognize(
			table.Crop(9, 13, 409, 154),
			"pot_samples/*",
		)
		wg.Done()
	}()

	go func() {
		digitThree = recognize(
			table.Crop(9, 13, 418, 154),
			"pot_samples/*",
		)
		wg.Done()
	}()

	go func() {
		digitFour = recognize(
			table.Crop(9, 13, 427, 154),
			"pot_samples/*",
		)
		wg.Done()
	}()

	wg.Wait()

	table.Pot = fmt.Sprintf("%s%s%s%s",
		digitOne, digitTwo, digitThree, digitFour,
	)
}

func (table *Table) HandRecognition() {
	wg := &sync.WaitGroup{}

	wg.Add(2)

	go func() {
		table.HandFirstCard = recognize(
			table.Crop(46, 30, 346, 340),
			"card_samples/*/*",
		)
		wg.Done()
	}()

	go func() {
		table.HandSecondCard = recognize(
			table.Crop(46, 30, 396, 340),
			"card_samples/*/*",
		)
		wg.Done()
	}()

	wg.Wait()
}

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil
}

//@TODO: refactor duplicated code
func recognize(search string, samplesFilenamePattern string) string {
	files, _ := filepath.Glob(samplesFilenamePattern)

	for _, file := range files {
		// try -metric AE
		command, _ := cmdRunner.Command(fmt.Sprintf(
			"/bin/compare -quiet -metric RMSE %s %s NULL:",
			file, search))

		_, err := command.Run()
		if err == nil {
			matches := reMatchedSample.FindStringSubmatch(file)

			if len(matches) != 0 {
				return strings.Replace(matches[1], "/", "", 1)
			}
		} else {
			matches := reCompareErrorLevel.FindStringSubmatch(err.Error())

			if len(matches) != 0 {
				errorLevel, _ := strconv.ParseFloat(matches[1], 32)
				if errorLevel < compareThreshold {
					matches := reMatchedSample.FindStringSubmatch(file)

					if len(matches) != 0 {
						return strings.Replace(matches[1], "/", "", 1)
					}
				}
			}
		}
	}

	return ""
}

func (table *Table) Crop(
	cropWidth int,
	cropHeight int,
	offsetX int,
	offsetY int,
) string {
	cropped, _ := getTmpFilename()

	command, _ := cmdRunner.Command(fmt.Sprintf(
		"/bin/convert -crop %dx%d+%d+%d %s %s",
		cropWidth, cropHeight, offsetX, offsetY,
		table.Screenshot, cropped),
	)

	_, _ = command.Run()

	return cropped
}

func makeTableScreenshot() (string, error) {
	tableScreenshot, err := getTmpFilename()
	if err != nil {
		return "", err
	}

	command, err := cmdRunner.Command(fmt.Sprintf(
		"/bin/import png:%s", tableScreenshot),
	)

	if err != nil {
		return "", err
	}

	_, err = command.Run()
	if err != nil {
		return "", err
	}

	return tableScreenshot, nil
}

func getWindowId() (string, error) {
	command, err := cmdRunner.Command("/bin/xprop")

	output, err := command.Run()
	if err != nil {
		return "", err
	}

	matches := reWindowId.FindStringSubmatch(output[11])

	if len(matches) != 0 {
		return matches[1], nil
	}

	return "", errors.New("No window id found")
}

func getTmpFilename() (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "croc")
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

//@TODO:
// 1) parse button's position, hero position
// 1.1) positions: dealer, sb, bb, mp2, mp3, cutoff
// 2) parse current blinds
