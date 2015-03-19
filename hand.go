package main

type ImageSnippet struct {
	Width   int
	Height  int
	OffsetX int
	OffsetY int
}

type Card struct {
	ImageSnippet
}

type Hand struct {
	LeftCard  Card
	RightCard Card
}
