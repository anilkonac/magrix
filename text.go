package main

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	dpi                   = 72
	fontSizeIntro         = 20
	textIntro             = "Loading the locations of terminals that lift plasma gates..."
	textTerminalBlue      = "Lifting the blue plasma gate"
	textTerminalOrange    = "Lifting the orange plasma gate"
	durationTextIntroSec  = 3.0
	durationTextTerminals = 2.0
	introTextShiftY       = screenHeight / 4.0
)

var (
	//go:embed assets/fonts/Minecraft.ttf
	bytesFontMinecraft            []byte
	fontFaceIntro                 font.Face
	showTextIntro                 bool
	showTextTerminalBlue          bool
	showTextTerminalOrange        bool
	imageTextIntro                *ebiten.Image
	imageTextTerminalBlue         *ebiten.Image
	imageTextTerminalOrange       *ebiten.Image
	imageTextButton               *ebiten.Image
	drawOptionsTextIntro          ebiten.DrawImageOptions
	drawOptionsTextTerminalBlue   ebiten.DrawImageOptions
	drawOptionsTextTerminalOrange ebiten.DrawImageOptions
)

func init() {
	tt, err := opentype.Parse(bytesFontMinecraft)
	panicErr(err)

	fontFaceIntro, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeIntro,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	panicErr(err)

	// Prepare intro text
	boundText := text.BoundString(fontFaceIntro, textIntro)
	boundTextSize := boundText.Size()
	imageTextIntro = ebiten.NewImage(boundTextSize.X, boundTextSize.Y)
	text.Draw(imageTextIntro, textIntro, fontFaceIntro, -boundText.Min.X, -boundText.Min.Y, colorOrange)
	drawOptionsTextIntro.GeoM.Reset()
	drawOptionsTextIntro.GeoM.Translate(
		float64((screenWidth-boundTextSize.X)/2.0-boundText.Min.X),
		float64((screenHeight-boundTextSize.Y)/2.0-boundText.Min.Y)+introTextShiftY)

	// Prepare blue terminal text
	boundText = text.BoundString(fontFaceIntro, textTerminalBlue)
	boundTextSize = boundText.Size()
	imageTextTerminalBlue = ebiten.NewImage(boundTextSize.X, boundTextSize.Y)
	text.Draw(imageTextTerminalBlue, textTerminalBlue, fontFaceIntro, -boundText.Min.X, -boundText.Min.Y, colorBlue)
	drawOptionsTextTerminalBlue.GeoM.Reset()
	drawOptionsTextTerminalBlue.GeoM.Translate(
		float64((screenWidth-boundTextSize.X)/2.0-boundText.Min.X),
		float64((screenHeight-boundTextSize.Y)/2.0-boundText.Min.Y)+introTextShiftY)

	// Prepare orange terminal text
	boundText = text.BoundString(fontFaceIntro, textTerminalOrange)
	boundTextSize = boundText.Size()
	imageTextTerminalOrange = ebiten.NewImage(boundTextSize.X, boundTextSize.Y)
	text.Draw(imageTextTerminalOrange, textTerminalOrange, fontFaceIntro, -boundText.Min.X, -boundText.Min.Y, colorOrange)
	drawOptionsTextTerminalOrange.GeoM.Reset()
	drawOptionsTextTerminalOrange.GeoM.Translate(
		float64((screenWidth-boundTextSize.X)/2.0-boundText.Min.X),
		float64((screenHeight-boundTextSize.Y)/2.0-boundText.Min.Y)+introTextShiftY)
}
