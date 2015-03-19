package main

type Number struct {
	Digits []ImageSnippet
}

func (image Image) NumberTypeRecognize(
	snippet ImageSnippet,
	typeSamples string,
) (string, error) {
	typeString, err := recognize(
		image.Crop(snippet), typeSamples,
	)

	if err != nil {
		return "", err
	}

	return typeString, err
}

func getDigitsImageSnippets(
	width int,
	height int,
	offsetY int,
	offsets []int,
) []ImageSnippet {
	imageSnippets := make([]ImageSnippet, len(offsets))

	for index, offset := range offsets {
		imageSnippets[index] = ImageSnippet{
			Width:   width,
			Height:  height,
			OffsetX: offset,
			OffsetY: offsetY,
		}
	}

	return imageSnippets
}
