package main

import (
	"bytes"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type stateMusic uint8

const (
	musicOn stateMusic = iota
	musicPaused
	musicMuted
)

const (
	sampleRate  = 44100
	volumeMusic = 0.5
)

var (
	//go:embed assets/sounds/RaceToMars.ogg
	bytesMusic   []byte
	audioContext *audio.Context
	playerMusic  *audio.Player
	musicState   stateMusic
)

func init() {

	audioContext = audio.NewContext(sampleRate)
	stream, err := vorbis.DecodeWithSampleRate(sampleRate, bytes.NewReader(bytesMusic))
	panicErr(err)

	playerMusic, err = audioContext.NewPlayer(stream)
	panicErr(err)
	playerMusic.SetVolume(volumeMusic)

	playerMusic.Play()
}
