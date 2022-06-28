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
	fontSizeIntro         = 24
	fontSizeButton        = 48
	textIntro             = "Locating terminals controlling plasma walls..."
	textTerminalBlue      = "Disabling the blue plasma wall..."
	textTerminalOrange    = "Disabling the orange plasma wall..."
	textButton            = "Mission Accomplished!"
	textFail              = "Mission Failed!"
	durationTextIntroSec  = 2.5
	durationTextTerminals = 2.0
	introTextShiftY       = screenHeight / 4.0
)

var (
	//go:embed assets/fonts/Minecraft.ttf
	bytesFontMinecraft            []byte
	fontFaceIntro                 font.Face
	fontFaceButton                font.Face
	showTextIntro                 bool
	showTextTerminalBlue          bool
	showTextTerminalOrange        bool
	showTextButton                bool
	imageTextIntro                *ebiten.Image
	imageTextTerminalBlue         *ebiten.Image
	imageTextTerminalOrange       *ebiten.Image
	imageTextButton               *ebiten.Image
	imageTextFail                 *ebiten.Image
	drawOptionsTextIntro          ebiten.DrawImageOptions
	drawOptionsTextTerminalBlue   ebiten.DrawImageOptions
	drawOptionsTextTerminalOrange ebiten.DrawImageOptions
	drawOptionsTextButton         ebiten.DrawImageOptions
	drawOptionsTextFail           ebiten.DrawImageOptions
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

	fontFaceButton, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeButton,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	panicErr(err)

	// Prepare intro text
	boundText := text.BoundString(fontFaceIntro, textIntro)
	boundTextSize := boundText.Size()
	imageTextIntro = ebiten.NewImage(boundTextSize.X, boundTextSize.Y)
	text.Draw(imageTextIntro, textIntro, fontFaceIntro, -boundText.Min.X, -boundText.Min.Y, colorGreen)
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

	// Prepare final(button) text
	boundText = text.BoundString(fontFaceButton, textButton)
	boundTextSize = boundText.Size()
	imageTextButton = ebiten.NewImage(boundTextSize.X, boundTextSize.Y)
	text.Draw(imageTextButton, textButton, fontFaceButton, -boundText.Min.X, -boundText.Min.Y, colorGreen)
	drawOptionsTextButton.GeoM.Reset()
	drawOptionsTextButton.GeoM.Translate(
		float64((screenWidth-boundTextSize.X)/2.0-boundText.Min.X),
		float64((screenHeight-boundTextSize.Y)/2.0-boundText.Min.Y)+introTextShiftY)

	// Prepare Fail text
	boundText = text.BoundString(fontFaceButton, textFail)
	boundTextSize = boundText.Size()
	imageTextFail = ebiten.NewImage(boundTextSize.X, boundTextSize.Y)
	text.Draw(imageTextFail, textFail, fontFaceButton, -boundText.Min.X, -boundText.Min.Y, colorGunAttract)
	drawOptionsTextFail.GeoM.Reset()
	drawOptionsTextFail.GeoM.Translate(
		float64((screenWidth-boundTextSize.X)/2.0-boundText.Min.X),
		float64((screenHeight-boundTextSize.Y)/2.0-boundText.Min.Y)+introTextShiftY)
}
