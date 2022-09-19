package asset

import (
	"bytes"
	"embed"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *.png sounds fonts gameMap.tmx
var fs embed.FS

var (
	AnimEnemy1Idle     = "enemy1_idle.png"
	AnimPlayerIdle     = "player_idle.png"
	AnimPlayerWalk     = "player_walk.png"
	AnimRocket         = "rocket_anim.png"
	AnimExplosion      = "Explosion_duplicateframes.png"
	AnimElectricBlue   = "anim_electric_blue.png"
	AnimElectricOrange = "anim_electric_orange.png"

	SpriteTerminalBlue   = "terminal_blue.png"
	SpriteTerminalOrange = "terminal_orange.png"
	SpriteTerminalGreen  = "terminal_green.png"
	SpriteButton         = "theButton.png"
	SpriteGun            = "sprite_gun.png"

	ImageHeart               = "heart.png"
	ImageArrow               = "arrow.png"
	ImageMapLayerPlatforms   = "map_layer_platforms.png"
	ImageMapLayerDecorations = "map_layer_decorations.png"

	FontMinecraft = "fonts/Minecraft.ttf"

	Music          = "sounds/RaceToMars.ogg"
	SoundExplosion = "sounds/explosion.wav"

	Map = "gameMap.tmx"
)

func Bytes(path string) []byte {
	bytes, err := fs.ReadFile(path)
	panik(err)

	return bytes
}

func Image(path string) *ebiten.Image {
	img, err := png.Decode(bytes.NewReader(Bytes(path)))
	panik(err)

	return ebiten.NewImageFromImage(img)
}

func panik(err error) {
	if err != nil {
		panic(err)
	}
}
