package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/docopt/docopt-go"
	"github.com/op/go-logging"
	"github.com/thearkit/runcmd"
)

var (
	log                 = logging.MustGetLogger("croc")
	reCompareErrorLevel = regexp.MustCompile("\\((.*)\\)$")
	compareThreshold    = 0.05
	cmdRunner           *runcmd.Local
	compareCmd          = "/bin/compare -quiet -metric RMSE %s %s NULL:"
	cardSamples         = "cards/*"
)

const usage = `
	Usage:
	croc
	croc <filepath>
`

type Image struct {
	Path string
}

func (image Image) Crop(snippet ImageSnippet) string {
	croppedPath, _ := getTmpFilename()

	command, _ := cmdRunner.Command(fmt.Sprintf(
		"/bin/convert -crop %dx%d+%d+%d %s %s",
		snippet.Width, snippet.Height, snippet.OffsetX, snippet.OffsetY,
		image.Path, croppedPath),
	)
	_, _ = command.Run()

	return croppedPath
}

func main() {
	var err error

	cmdRunner, err = runcmd.NewLocalRunner()
	if err != nil {
		log.Fatal(err)
	}

	image := Image{}

	args, _ := docopt.Parse(usage, nil, true, "croc", false)
	if args["<filepath>"] != nil {
		image.Path = args["<filepath>"].(string)
	} else {
		image.Path, err = makeScreenshot()
		if err != nil {
			log.Fatal(err)
		}
	}

	hand := Hand{
		LeftCard: Card{
			ImageSnippet: ImageSnippet{
				Width:   46,
				Height:  30,
				OffsetX: 346,
				OffsetY: 340,
			},
		},
		RightCard: Card{
			ImageSnippet: ImageSnippet{
				Width:   46,
				Height:  30,
				OffsetX: 396,
				OffsetY: 340,
			},
		},
	}

	leftCard, _ := recognize(image.Crop(hand.LeftCard.ImageSnippet), cardSamples)
	rightCard, _ := recognize(image.Crop(hand.RightCard.ImageSnippet), cardSamples)
	log.Notice("%v", leftCard)
	log.Notice("%v", rightCard)
}

func recognize(
	input string,
	samplesFilepathPattern string,
) (string, error) {
	samples, _ := filepath.Glob(samplesFilepathPattern)

	for _, sample := range samples {
		command, _ := cmdRunner.Command(fmt.Sprintf(compareCmd, sample, input))

		_, err := command.Run()
		if err == nil {
			return path.Base(sample), nil
		}

		errorLevel, _ := strconv.ParseFloat(
			reCompareErrorLevel.FindStringSubmatch(err.Error())[1],
			32,
		)

		if errorLevel < compareThreshold {
			return path.Base(sample), nil
		}
	}

	return "", errors.New("Recognition failed")
}

func makeScreenshot() (string, error) {
	screenshot, err := getTmpFilename()
	if err != nil {
		return "", err
	}

	command, err := cmdRunner.Command(fmt.Sprintf(
		"/bin/import png:%s", screenshot),
	)

	if err != nil {
		return "", err
	}

	_, err = command.Run()
	if err != nil {
		return "", err
	}

	return screenshot, nil
}

func getTmpFilename() (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "croc")
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}
