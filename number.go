package main

type Number struct {
	Digits []ImageSnippet
}

func (image Image) NumberTypeRecognize(
	snippet ImageSnippet,
	typeSamples string,
	compareThreshold float64,
) (string, error) {
	typeString, err := recognize(
		image.Crop(snippet), typeSamples, compareThreshold,
	)

	if err != nil {
		return "", err
	}

	return typeString, err
}
