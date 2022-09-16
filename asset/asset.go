package asset

import "embed"

//go:embed
var fs embed.FS

var (
	AnimEnemy1Idle = "enemy1_idle.png"
	AnimPlayerIdle = "player_idle.png"
	AnimPlayerWalk = "player_walk.png"
	AnimRocket     = "rocket_anim.png"
	AnimExplosion  = "Explosion_duplicateframes.png"

	SpriteElectricBlue   = "anim_electric_blue.png"
	SpriteElectricOrange = "anim_electric_orange.png"
	SpriteTerminalBlue   = "anim_electric_orange.png"
	SpriteTerminalOrange = "anim_electric_orange.png"
	SpriteButton         = "theButton.png"

	ImageHeart = "heart.png"
	ImageArrow = "arrow.png"
)

func GetBytes(path string) []byte {
	bytes, err := fs.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}
