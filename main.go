package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 1040
	screenHeight = screenWidth * 9 / 16
	fontHeight   = 12
	totalStates  = 16
	fontWidth    = fontHeight - 2
	yLength      = screenHeight / fontHeight
	xLength      = screenWidth / fontWidth

	siteDir      = "site"
	assetsDir    = "assets"
	textFilepath = "exp.txt"
	fontFilepath = "SpaceMono-Bold.ttf"
)

var (
	hFont       font.Face
	err         error
	worldWords  string
	lowerColor  colorful.Color
	higherColor colorful.Color
	baseColor   colorful.Color
)

func init() {
	rand.Seed(time.Now().UnixNano())

	baseColor, _ = colorful.Hex("#060606")
	higherColor, _ = colorful.Hex("#f52552")
	lowerColor, _ = colorful.Hex("#126cfc")

	var textData, fontData []byte

	if runtime.GOOS == "js" {
		r, err := http.Get(path.Join(assetsDir, textFilepath))
		if err != nil {
			log.Fatal(err)
		}

		textData, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		r, err = http.Get(path.Join(assetsDir, fontFilepath))
		if err != nil {
			log.Fatal(err)
		}

		fontData, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		textData, err = ioutil.ReadFile(path.Join(siteDir, assetsDir, textFilepath))
		if err != nil {
			log.Fatal(err)
		}

		// fontData, err := ioutil.ReadFile("PTMono-Regular.ttf")
		fontData, err = ioutil.ReadFile(path.Join(siteDir, assetsDir, fontFilepath))
		if err != nil {
			log.Fatal(err)
		}
	}

	worldWords = strings.ToUpper(string(textData))
	// worldWords = strings.ReplaceAll(worldWords, " ", "")
	// worldWords = strings.ReplaceAll(worldWords, ".", "")
	// worldWords = strings.ReplaceAll(worldWords, ",", "")
	// worldWords = strings.ReplaceAll(worldWords, ":", "")
	// worldWords = strings.ReplaceAll(worldWords, ";", "")
	worldWords = strings.ReplaceAll(worldWords, "\n", " ")

	tt, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 92

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

// Game ...
type Game struct {
	counter int
	board   [][]uint8
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

			if g.board[j][i] == totalStates {
				n++
			}
		}
	}

	return n
}

func (g *Game) nextState() [][]uint8 {
	nextBoard := make([][]uint8, yLength)
	for y := 0; y < yLength; y++ {
		for x := 0; x < xLength; x++ {
			if nextBoard[y] == nil {
				nextBoard[y] = make([]uint8, xLength)
			}

			n := g.neighbours(x, y)

			if n == 3 || (n == 2 && g.board[y][x] == totalStates) {
				nextBoard[y][x] = totalStates
			} else {
				if g.board[y][x] > 0 {
					nextBoard[y][x] = g.board[y][x] - 1
				}
			}

		}
	}

	return nextBoard
}

// Update ...
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

			if g.board[y][x] == 0 { // dead
				text.Draw(screen, c, hFont, i, j, baseColor)
			} else { // live
				r := float64(g.board[y][x]) / float64(totalStates)
				col := lowerColor.BlendLuv(higherColor, r)
				text.Draw(screen, c, hFont, i, j, col)
			}

			index++
		}
	}

	nextBoard := g.nextState()
	g.board = nextBoard

	return nil
}

// Layout ...
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		board: make([][]uint8, yLength),
	}

	for y := 0; y < yLength; y++ {
		for x := 0; x < xLength; x++ {
			if g.board[y] == nil {
				g.board[y] = make([]uint8, xLength)
			}
			if rand.Float64() > 0.5 {
				g.board[y][x] = totalStates
			} else {
				g.board[y][x] = 0
			}

		}
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ruins Automaton")
	ebiten.SetRunnableInBackground(true)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetMaxTPS(20)
	// ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
