package main

import "strconv"

//@TODO: move to config
var (
	opponentsSamples          = "opponents/*"
	opponentsCompareThreshold = 0.03
	opponentsWidth            = 6
	opponentsHeight           = 10
)

var opponents map[int]ImageSnippet

func (image Image) OpponentsRecognize() string {
	opponents := map[int]ImageSnippet{
		2: ImageSnippet{
			opponentsWidth, opponentsHeight, 92, 357,
		},

		3: ImageSnippet{
			opponentsWidth, opponentsHeight, 0, 241,
		},

		4: ImageSnippet{
			opponentsWidth, opponentsHeight, 43, 128,
		},

		5: ImageSnippet{
			opponentsWidth, opponentsHeight, 199, 72,
		},

		6: ImageSnippet{
			opponentsWidth, opponentsHeight, 587, 72,
		},

		7: ImageSnippet{
			opponentsWidth, opponentsHeight, 743, 128,
		},

		8: ImageSnippet{
			opponentsWidth, opponentsHeight, 786, 241,
		},

		9: ImageSnippet{
			opponentsWidth, opponentsHeight, 695, 357,
		},
	}

	opponentsSeats := ""

	for seat, opponent := range opponents {
		_, err := recognize(
			image.Crop(opponent),
			opponentsSamples,
			opponentsCompareThreshold,
		)

		if err != nil {
			continue
		}

		opponentsSeats += strconv.Itoa(seat)
	}

	return opponentsSeats
}
