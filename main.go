package main

import (
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 1040
	screenHeight = 670
	fontHeight   = 12
	fontWidth    = fontHeight - 2
	yLength      = screenHeight / fontHeight
	xLength      = screenWidth / fontWidth
)

var bFont font.Face
var hFont font.Face

var err error

var worldWords string

var colorA = color.RGBA{R: 80, G: 80, B: 80, A: 255}

var colorB = color.RGBA{R: 255, G: 255, B: 255, A: 255}

// 25, 247, 110

// var colorB = color.RGBA{R: 25, G: 247, B: 110, A: 255}

func init() {
	rand.Seed(time.Now().UnixNano())

	textData, err := ioutil.ReadFile("exp.txt")
	if err != nil {
		log.Fatal(err)
	}

	worldWords = strings.ToUpper(string(textData))
	worldWords = strings.ReplaceAll(worldWords, " ", "")
	worldWords = strings.ReplaceAll(worldWords, ".", "")
	worldWords = strings.ReplaceAll(worldWords, ",", "")
	worldWords = strings.ReplaceAll(worldWords, ":", "")
	worldWords = strings.ReplaceAll(worldWords, ";", "")
	worldWords = strings.ReplaceAll(worldWords, "\n", "")

	// fontData, err := ioutil.ReadFile("PTMono-Regular.ttf")
	fontData, err := ioutil.ReadFile("SpaceMono-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 92

	bFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontHeight,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	fontData, err = ioutil.ReadFile("SpaceMono-Bold.ttf")
	if err != nil {
		log.Fatal(err)
	}

	tt, err = opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	hFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontHeight,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	counter int
	board   [][]bool
}

func (g *Game) neighbours(x, y int) int {
	n := 0

	for nx := -1; nx <= 1; nx++ {
		for ny := -1; ny <= 1; ny++ {
			if nx == 0 && ny == 0 {
				continue
			}

			i := x + nx
			j := y + ny

			if i < 0 || i >= xLength || j < 0 || j >= yLength {
				continue
			}

			if g.board[j][i] {
				n++
			}
		}
	}

	return n
}

func (g *Game) nextState() [][]bool {
	nextBoard := make([][]bool, yLength)
	for y := 0; y < yLength; y++ {
		for x := 0; x < xLength; x++ {
			if nextBoard[y] == nil {
				nextBoard[y] = make([]bool, xLength)
			}

			n := g.neighbours(x, y)

			nextBoard[y][x] = n == 3 || (n == 2 && g.board[y][x])
		}
	}

	return nextBoard
}

func (g *Game) Update(screen *ebiten.Image) error {
	xOffset := 0
	yOffset := 12

	letters := strings.Split(worldWords, "")
	index := 0
	for y := 0; y < yLength; y++ {
		for x := 0; x < xLength; x++ {
			i := xOffset + x*fontWidth
			j := yOffset + y*fontHeight
			c := letters[index%len(letters)]

			if g.board[y][x] {
				text.Draw(screen, c, hFont, i, j, colorB)
			} else {
				text.Draw(screen, c, hFont, i, j, colorA)
			}

			index++
		}
	}

	nextBoard := g.nextState()
	g.board = nextBoard

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		board: make([][]bool, yLength),
	}

	for y := 0; y < yLength; y++ {
		for x := 0; x < xLength; x++ {
			if g.board[y] == nil {
				g.board[y] = make([]bool, xLength)
			}
			if rand.Float64() > 0.5 {
				g.board[y][x] = true
			} else {
				g.board[y][x] = false
			}

		}
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("LasHabMun")
	ebiten.SetMaxTPS(16)
	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
