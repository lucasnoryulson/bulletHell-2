// Construido como parte da disciplina: Sistemas Distribuidos - PUCRS - Escola Politecnica
// Professor: Fernando Dotti  (https://fldotti.github.io/)
/*
   Este jogo foi adaptado para suportar dois processos distribuídos.
   Cada processo representa um jogador:
     - Processo 0: controla o Player 1
     - Processo 1: controla o Player 2
   Ambos compartilham o acesso à "área crítica" (atualização do jogo) via o módulo DIMEX.
*/

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/brunobaa/bullethell/DIMEX"

	"github.com/nsf/termbox-go"
)

const (
	WorldWidth    = 30
	WorldHeight   = 15
	UpdatesPerSec = 10
	MaxLives      = 5
)

var (
	BulletSpeed     = 1
	BulletsPerSpawn = 1
	SpawnInterval   = 1000
)

type Entity struct {
	X, Y int
	Ch   rune
}

type Bullet struct {
	Entity
	DirectionX int
	DirectionY int
	Active     bool
}

type Player struct {
	Entity
	Lives int
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func render(entities []Entity, bullets []Bullet, player Player, tick int) {
	grid := make([][]rune, WorldHeight)
	for y := range grid {
		grid[y] = make([]rune, WorldWidth)
		for x := range grid[y] {
			grid[y][x] = ' '
		}
	}

	for x := 0; x < WorldWidth; x++ {
		grid[0][x] = '#'
		grid[WorldHeight-1][x] = '#'
	}
	for y := 0; y < WorldHeight; y++ {
		grid[y][0] = '#'
		grid[y][WorldWidth-1] = '#'
	}

	grid[player.Y][player.X] = player.Ch
	for _, e := range entities {
		if e.Y > 0 && e.Y < WorldHeight-1 && e.X > 0 && e.X < WorldWidth-1 {
			grid[e.Y][e.X] = e.Ch
		}
	}
	for _, b := range bullets {
		if b.Active && b.Y > 0 && b.Y < WorldHeight-1 && b.X > 0 && b.X < WorldWidth-1 {
			grid[b.Y][b.X] = b.Ch
		}
	}

	clearScreen()
	for _, row := range grid {
		for _, cell := range row {
			fmt.Print(string(cell))
		}
		fmt.Println()
	}

	fmt.Print("\nVidas: ")
	for i := 0; i < player.Lives; i++ {
		fmt.Print("♥ ")
	}
	for i := player.Lives; i < MaxLives; i++ {
		fmt.Print("♡ ")
	}
	fmt.Printf("\nTick: %d\n", tick)
}

func handleInput(player *Player, done chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			ev := termbox.PollEvent()
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyArrowUp:
					if player.Y > 1 {
						player.Y--
					}
				case termbox.KeyArrowDown:
					if player.Y < WorldHeight-2 {
						player.Y++
					}
				case termbox.KeyArrowLeft:
					if player.X > 1 {
						player.X--
					}
				case termbox.KeyArrowRight:
					if player.X < WorldWidth-2 {
						player.X++
					}
				case termbox.KeyEsc:
					done <- true
					return
				}
			}
		}
	}
}

func checkCollision(bullet Bullet, player Player) bool {
	return bullet.X == player.X && bullet.Y == player.Y
}

func updateBullets(bullets []Bullet, player *Player) []Bullet {
	for i := range bullets {
		if bullets[i].Active {
			bullets[i].X += bullets[i].DirectionX * BulletSpeed
			bullets[i].Y += bullets[i].DirectionY * BulletSpeed

			if checkCollision(bullets[i], *player) {
				bullets[i].Active = false
				if player.Lives > 0 {
					player.Lives--
				}
			}

			if bullets[i].X <= 0 || bullets[i].X >= WorldWidth-1 ||
				bullets[i].Y <= 0 || bullets[i].Y >= WorldHeight-1 {
				bullets[i].Active = false
			}
		}
	}
	return bullets
}

func spawnBullet() Bullet {
	side := rand.Intn(4)
	var bullet Bullet
	bullet.Ch = '*'
	bullet.Active = true

	switch side {
	case 0:
		bullet.X = rand.Intn(WorldWidth-2) + 1
		bullet.Y = 1
		bullet.DirectionX = 0
		bullet.DirectionY = 1
	case 1:
		bullet.X = WorldWidth - 2
		bullet.Y = rand.Intn(WorldHeight-2) + 1
		bullet.DirectionX = -1
		bullet.DirectionY = 0
	case 2:
		bullet.X = rand.Intn(WorldWidth-2) + 1
		bullet.Y = WorldHeight - 2
		bullet.DirectionX = 0
		bullet.DirectionY = -1
	case 3:
		bullet.X = 1
		bullet.Y = rand.Intn(WorldHeight-2) + 1
		bullet.DirectionX = 1
		bullet.DirectionY = 0
	}
	return bullet
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run jogo.go <id> <addr0> <addr1> ...")
		return
	}

	id, _ := strconv.Atoi(os.Args[1])
	addresses := os.Args[2:]
	dmx := DIMEX.NewDIMEX(addresses, id, true)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	var player Player
	if id == 0 {
		player = Player{Entity: Entity{X: 2, Y: 2, Ch: '1'}, Lives: MaxLives}
	} else {
		player = Player{Entity: Entity{X: WorldWidth - 3, Y: WorldHeight - 3, Ch: '2'}, Lives: MaxLives}
	}

	bullets := make([]Bullet, 0)
	done := make(chan bool)
	go handleInput(&player, done)

	ticker := time.NewTicker(time.Second / UpdatesPerSec)
	defer ticker.Stop()

	spawnTicker := time.NewTicker(time.Duration(SpawnInterval) * time.Millisecond)
	defer spawnTicker.Stop()

	tick := 0
	for {
		select {
		case <-done:
			return
		case <-spawnTicker.C:
			for i := 0; i < BulletsPerSpawn; i++ {
				bullets = append(bullets, spawnBullet())
			}
		case <-ticker.C:
			dmx.Req <- DIMEX.ENTER
			<-dmx.Ind

			tick++
			bullets = updateBullets(bullets, &player)

			var other Entity
			if id == 0 {
				other = Entity{X: WorldWidth - 3, Y: WorldHeight - 3, Ch: '2'}
			} else {
				other = Entity{X: 2, Y: 2, Ch: '1'}
			}
			render([]Entity{other}, bullets, player, tick)
			dmx.Req <- DIMEX.EXIT
		}
	}
}
