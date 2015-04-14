package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
)

var (
	compareCmd = "/bin/compare -dissimilarity-threshold 1 " +
		"-quiet -metric RMSE %s %s NULL:"

	convertCmd = "/bin/convert -crop %dx%d+%d+%d %s %s"

	reCompareErrorLevel = regexp.MustCompile("\\((.*)\\).*$")

	subimageSearchCmd = "visgrep"

	reSubimageOffsets = regexp.MustCompile("^(\\d+),(\\d+)")
)

func (image Image) Crop(snippet ImageSnippet) string {
	croppedPath, _ := getTmpFilename()

	command, _ := cmdRunner.Command(fmt.Sprintf(
		convertCmd,
		snippet.Width, snippet.Height, snippet.OffsetX, snippet.OffsetY,
		image.Path, croppedPath),
	)
	_, _ = command.Run()

	return croppedPath
}

func recognize(
	input string,
	samplesFilepathPattern string,
	compareThreshold float64,
) (string, error) {
	samples, _ := filepath.Glob(samplesFilepathPattern)

	for _, sample := range samples {
		command, _ := cmdRunner.Command(fmt.Sprintf(compareCmd, sample, input))

		_, err := command.Run()
		if err == nil {
			return path.Base(sample), nil
		}

		compareErrorLevel := reCompareErrorLevel.FindStringSubmatch(err.Error())

		if len(compareErrorLevel) > 0 {
			errorLevel, _ := strconv.ParseFloat(
				compareErrorLevel[1],
				32,
			)

			if errorLevel < compareThreshold {
				return path.Base(sample), nil
			}

		}
	}

	return "", errors.New(fmt.Sprintf("%s failed!", input))
}

func getSampleOffsets(
	image string,
	subimage string,
) (int, int) {
	cmd := exec.Command(subimageSearchCmd, image, subimage)

	var out bytes.Buffer
	cmd.Stdout = &out

	cmd.Run()

	offsets := reSubimageOffsets.FindStringSubmatch(out.String())

	x, y := -1, -1

	if len(offsets) == 3 {
		x, _ = strconv.Atoi(offsets[1])
		y, _ = strconv.Atoi(offsets[2])
	}

	return x, y
}
