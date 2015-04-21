package main

var manualHands = []string{
	"AA", "KK", "QQ", "JJ", "TT", "99", "88",
	"77", "66", "55", "44", "33", "22",

	"AK", "AKs",
}

var raiseFoldHandsLatePosition = []string{
	"AQ", "AQs", "AJs",

	"AJ", "KQ",
	"KQs", "KJs", "ATs",

	"AT", "A9", "A8",
	"KJ", "QJ", "KT",

	"A9s", "A8s", "A7s", "A6s", "A5s",
	"KTs", "K9s", "QJs", "QTs", "JTs",
	"T9s",

	"A7", "A6", "A5", "A4", "A3", "A2",
	"K9", "K8",
	"QT", "Q9",
	"JT",
	"T9", "98",

	"A4s", "A3s", "A2s",
	"K8s", "K7s", "K6s",
	"Q9s", "Q8s", "J9s",
	"98s", "87s", "76s", "65s",
}

var raiseFoldHands = map[string][]string{
	"EP": []string{
		"AQ", "AQs", "AJs",
	},
	"MP": []string{
		"AQ", "AQs", "AJs",

		"AJ", "KQ",
		"KQs", "KJs", "ATs",
	},
	"CO": []string{
		"AQ", "AQs", "AJs",

		"AJ", "KQ",
		"KQs", "KJs", "ATs",

		"AT", "A9", "A8",
		"KJ", "QJ", "KT",

		"A9s", "A8s", "A7s", "A6s", "A5s",
		"KTs", "K9s", "QJs", "QTs", "JTs",
		"T9s",
	},
	"BU": raiseFoldHandsLatePosition,
	"SB": raiseFoldHandsLatePosition,
	"BB": raiseFoldHandsLatePosition,
}

var stealFoldHandsBUandSB = []string{
	"99", "88", "77", "66", "55", "44", "33", "22",
	"AQ", "AQs", "AJ", "AJs", "AT", "ATs", "A9", "A9s",
	"A8", "A8s", "A7", "A7s", "A6s", "A5s", "A4s", "A3s", "A2s",
	"KQ", "KQs", "KJ", "KJs", "KT", "KTs",
	"QJ", "QJs", "QT", "QTs",
	"JT", "JTs",
	"T9", "T9s",
	"98s",
	"87s",
	"76s",
}

var stealFoldHands = map[string][]string{
	"CO": []string{
		"99", "88", "77", "66", "55", "44", "33", "22",
		"AQ", "AQs", "AJ", "AJs", "AT", "ATs", "A9s", "A8s", "A7s",
		"KQ", "KQs", "KJ", "KJs", "KT", "KTs",
		"QJ", "QJs", "QTs",
		"JTs",
	},
	"BU": stealFoldHandsBUandSB,
	"SB": stealFoldHandsBUandSB,
}

var restealFoldHands = []string{
	"AQ", "AQs", "AJ", "AJs", "AT", "ATs", "A9", "A9s",
	"99", "88", "77", "66", "55",
}

var threeBetHands = map[string][]string{
	"EP": []string{
		"QQ",
	},
	"MP": []string{
		"QQ", "JJ", "TT",
		"AQ", "AQs", "AK", "AKs",
	},
	"LATER": []string{
		"QQ", "JJ", "TT", "99", "88",
		"AQ", "AQs", "AK", "AKs",
	},
}
