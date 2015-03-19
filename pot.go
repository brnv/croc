package main

//@TODO: move to config
var (
	potNumberSamples = "pot_numbers/*"
	potTypeWidth     = 14
	potTypeHeight    = 13
	potTypeOffsetX   = 360
	potOffsetY       = 154
	potNumberWidth   = 9
	potNumberHeight  = 13
)

type Pot struct {
	Type       ImageSnippet
	PotNumbers []ImageSnippet
}

func (image Image) PotRecognize() {
	pot := Pot{
		Type: ImageSnippet{
			Width:   potTypeWidth,
			Height:  potTypeHeight,
			OffsetX: potTypeOffsetX,
			OffsetY: potOffsetY,
		},
		PotNumbers: []ImageSnippet{},
	}

	potType, err := recognize(
		image.Crop(pot.Type), "pot_types/*",
	)

	if err != nil {
		log.Notice("%v", err.Error())
	}

	switch potType {
	case "3":
		pot.PotNumbers = getImageSnippetsForThreePot()
	case "4":
		pot.PotNumbers = getImageSnippetsForFourPot()
	}

	potSize := ""

	for _, potNumber := range pot.PotNumbers {
		number, err := recognize(image.Crop(potNumber), potNumberSamples)
		if err != nil {
			log.Notice("%v", err.Error())
		}
		potSize += number
	}

	log.Notice("%v", potSize)
}

func getImageSnippetsForTwoPot() []ImageSnippet {
	return []ImageSnippet{
		getImageSnippetForOffset(407),
		getImageSnippetForOffset(416),
	}
}

func getImageSnippetsForThreePot() []ImageSnippet {
	return []ImageSnippet{
		getImageSnippetForOffset(403),
		getImageSnippetForOffset(412),
		getImageSnippetForOffset(421),
	}
}

func getImageSnippetsForFourPot() []ImageSnippet {
	return []ImageSnippet{
		getImageSnippetForOffset(397),
		getImageSnippetForOffset(409),
		getImageSnippetForOffset(418),
		getImageSnippetForOffset(427),
	}
}

func getImageSnippetForOffset(offsetX int) ImageSnippet {
	return ImageSnippet{
		Width:   potNumberWidth,
		Height:  potNumberHeight,
		OffsetX: offsetX,
		OffsetY: potOffsetY,
	}
}
