package main

import (
	"bytes"
	_ "encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	_ "math/rand"
)

type Field [][]int8

type PatternType int

const (
	_ PatternType = iota
	Glider
	Toad
	Exploder
	SmallExploder
	LightweightSpaceship
	RPentomino
	GosperGliderGun
)

var PatternNames = [...]string{
	"None",
	"Glider",
	"Toad",
	"Exploder",
	"SmallExploder",
	"LightweightSpaceship",
	"RPentomino",
	"GosperGliderGun",
}

var SeedPatterns map[PatternType]Field = make(map[PatternType]Field, 0)

func (s PatternType) String() string {
	return PatternNames[s]
}
func PatternTypeFromString(str string) PatternType {
	var pattern PatternType
	for i, s := range PatternNames {
		if s == str {
			pattern = PatternType(i)
			break
		}
	}

	return pattern

}

func init() {

	glider := Field{
		{0, 1, 0},
		{0, 0, 1},
		{1, 1, 1},
	}
	SeedPatterns[Glider] = glider

	toad := Field{
		{0, 0, 1, 0},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{0, 1, 0, 0},
	}
	SeedPatterns[Toad] = toad

	exploder := Field{
		{1, 0, 1, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 1, 0, 1},
	}
	SeedPatterns[Exploder] = exploder

	lightweight_spaceship := Field{
		{0, 1, 1, 1, 1},
		{1, 0, 0, 0, 1},
		{0, 0, 0, 0, 1},
		{1, 0, 0, 1, 0},
	}
	SeedPatterns[LightweightSpaceship] = lightweight_spaceship

	small_exploder := Field{
		{0, 1, 0},
		{1, 1, 1},
		{1, 0, 1},
		{0, 1, 0},
	}
	SeedPatterns[SmallExploder] = small_exploder

	rpentomino := Field{
		{0, 1, 1},
		{1, 1, 0},
		{0, 1, 0},
	}
	SeedPatterns[RPentomino] = rpentomino

	gosperglidergun := Field{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
		{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	SeedPatterns[GosperGliderGun] = gosperglidergun
}

type Game struct {
	GameId string `json:"id"`

	Generation int   `json:"generation"`
	Field      Field `json:"field"`
	FieldAux   Field `json:"_"`
	Height     int   `json:"height"`
	Width      int   `json:"width"`
}

func NewGame(height int, width int) *Game {

	game := new(Game)
	u1 := uuid.NewV4()
	game.GameId = u1.String()
	game.Generation = 1

	game.Height = height
	game.Width = width

	Logger.Debugf("h = %v, w = %v\n", height, width)
	game.Field = make(Field, width)
	for i := range game.Field {
		game.Field[i] = make([]int8, height)
	}

	game.FieldAux = make(Field, width)
	for i := range game.FieldAux {
		game.FieldAux[i] = make([]int8, height)
	}

	Logger.Debugf("Newgame game = %v\n", game)
	return game
}

func (this Game) String() string {

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "game: id=%v\n", this.GameId)
	return buf.String()
}

func (this *Game) Evolve() {

	Logger.Debugf("Evolve: before: %v\n", this.Field)
	fieldCopy(this.Field, this.FieldAux)

	this.Generation++
	for x := range this.FieldAux {
		for y := range this.FieldAux[x] {
			liveNeighbors := countLiveNeighbors(x, y, this.FieldAux)
			Logger.Debugf("liveNeighbors for (%v,%v) = %v\n", x, y, liveNeighbors)

			if this.FieldAux[x][y] > 0 { // alive
				if liveNeighbors < 2 {
					Logger.Debugf("(%v,%v) dying from lonliness\n", x, y)
					this.Field[x][y] = 0
				} else if liveNeighbors > 3 {
					Logger.Debugf("(%v,%v) dying from overcrowding\n", x, y)
					this.Field[x][y] = 0
				} else { // 2 or 3
					Logger.Debugf("(%v,%v) staying alive\n", x, y)
					this.Field[x][y] = this.Field[x][y]
				}
			} else { // dead
				if liveNeighbors == 3 {
					Logger.Debugf("(%v,%v) coming to life\n", x, y)

					this.Field[x][y] = int8(this.Generation)
				}
			}
		}
	}

	Logger.Debugf("Evolve: after: %v\n", this.Field)

}

func wrapCoordinate(c int, size int) int {
	result := c
	if c < 0 {
		result = size - 1
	} else if c > size-1 {
		result = 0
	}

	return result
}

func countLiveNeighbors(x int, y int, field Field) int {

	height := len(field)
	width := len(field[0])

	count := 0
	for hOffset := -1; hOffset <= 1; hOffset++ {
		for vOffset := -1; vOffset <= 1; vOffset++ {

			if hOffset == 0 && vOffset == 0 {
				continue
			}

			x1 := wrapCoordinate(x+hOffset, width)
			y1 := wrapCoordinate(y+vOffset, height)

			if field[x1][y1] > 0 {
				count++
			}
		}
	}

	return count

}

func fieldCopy(src Field, dst Field) {

	for x := 0; x < len(src); x++ {
		for y := 0; y < len(src[x]); y++ {
			Logger.Debugf("fieldCopy - copying val at %v,%v\n", x, y)
			dst[x][y] = src[x][y]
		}
	}

}

func (field Field) String() string {

	var buf bytes.Buffer
	for y := 0; y < len(field[0]); y++ {
		for x := 0; x < len(field); x++ {
			fmt.Fprintf(&buf, "%d", field[x][y])
		}
		fmt.Fprintf(&buf, "\n")
	}
	return buf.String()

}

func (this *Game) AddPattern(pattern PatternType, xpos int, ypos int) {

	Logger.Debugf("AddPattern - %v, (%v, %v)\n", pattern, xpos, ypos)
	patternField := SeedPatterns[pattern]

	Logger.Debugf("AddPattern height = %v, width = %v\n", len(patternField), len(patternField[0]))

	if len(patternField[0])+xpos >= this.Width {
		xpos = 1
	}
	if len(patternField)+ypos >= this.Height {
		ypos = 1
	}
	for y := 0; y < len(patternField); y++ {
		for x := 0; x < len(patternField[y]); x++ {
			this.Field[x+xpos][y+ypos] = patternField[y][x]
		}
	}
	/*
		for x := 0; x < len(patternField); x++ {
			for y := 0; y < len(patternField[x]); y++ {
				this.Field[x+xpos][y+ypos] = patternField[x][y]
			}
		}
	*/

}

func (this *Game) SeedPattern(pattern PatternType) {

	xpos := this.Width / 2
	ypos := this.Height / 2
	this.AddPattern(pattern, xpos, ypos)

}

func (this *Game) Seed() {

	this.SeedPattern(Glider)

	/*
		size := this.Height * this.Width
		seedSize := size / 10

		for i := 0; i < seedSize; i++ {

			x := rand.Intn(this.Width)
			y := rand.Intn(this.Height)

			this.Field[x][y] = true
		}
	*/
}

func (this *Game) Step() {
	Logger.Debugf("Step\n")
}

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
