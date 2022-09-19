package main

import (
	"bytes"
	"time"

	"github.com/anilkonac/magrix/asset"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type stateMusic uint8

const (
	musicOn stateMusic = iota
	musicPaused
	musicMuted
)

const (
	sampleRate      = 44100
	volumeMusic     = 0.3
	volumeExplosion = 0.5
	musicCheckSec   = 5.0
)

var (
	playerMusic     *audio.Player
	playerExplosion *audio.Player
	musicState      stateMusic
)

func init() {

	audioContext := audio.NewContext(sampleRate)
	streamMusic, err := vorbis.DecodeWithSampleRate(sampleRate, bytes.NewReader(asset.Bytes(asset.Music)))
	panicErr(err)

	playerMusic, err = audioContext.NewPlayer(streamMusic)
	panicErr(err)
	playerMusic.SetVolume(volumeMusic)

	playerMusic.Play()
	go repeatMusic()

	streamSound, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(asset.Bytes(asset.SoundExplosion)))
	panicErr(err)

	playerExplosion, err = audioContext.NewPlayer(streamSound)
	panicErr(err)

}

// Goroutine
func repeatMusic() {
	ticker := time.NewTicker(time.Second * musicCheckSec)
	for range ticker.C {
		if (musicState == musicOn) && !playerMusic.IsPlaying() {
			playerMusic.Rewind()
			playerMusic.Play()
		}
	}
}
