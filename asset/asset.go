package asset

import "embed"

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

	ImageHeart = "heart.png"
	ImageArrow = "arrow.png"

	FontMinecraft = "fonts/Minecraft.ttf"

	Music          = "sounds/RaceToMars.ogg"
	SoundExplosion = "sounds/explosion.wav"

	Map = "gameMap.tmx"
)

func Bytes(path string) []byte {
	bytes, err := fs.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return bytes
}
