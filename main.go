package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/docopt/docopt-go"
	"github.com/op/go-logging"
	"github.com/thearkit/runcmd"
)

var (
	log        = logging.MustGetLogger("croc")
	reWindowId = regexp.MustCompile("window id # of group leader: (0x.*)")
	cmdRunner  *runcmd.Local
)

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

	args, _ := docopt.Parse(usage, nil, true, "croc", false)
	if args["<filepath>"] != nil {
		log.Notice("%v", "File mode")
	} else {
		log.Notice("%v", "Window mode")

		tableScreenshot, err := makeTableScreenshot()
		if err != nil {
			log.Fatal(err)
		}

		log.Notice("%v", tableScreenshot)
	}
}

func makeTableScreenshot() (string, error) {
	tableScreenshot, err := getTmpFilename("table")
	if err != nil {
		return "", err
	}

	c, err := cmdRunner.Command(fmt.Sprintf(
		"/bin/import png:%s", tableScreenshot),
	)

	if err != nil {
		return "", err
	}

	_, err = c.Run()
	if err != nil {
		return "", err
	}

	return tableScreenshot, nil
}

func getWindowId() (string, error) {
	c, err := cmdRunner.Command("/bin/xprop")

	output, err := c.Run()
	if err != nil {
		return "", err
	}

	matches := reWindowId.FindStringSubmatch(output[11])

	if len(matches) != 0 {
		return matches[1], nil
	}

	return "", errors.New("No window id found")
}

func getTmpFilename(postfix string) (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "croc-"+postfix)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

//@TODO:
// 1) parse button's position, hero position
// 1.1) positions: dealer, sb, bb, mp2, mp3, cutoff
// 2) parse current blinds
