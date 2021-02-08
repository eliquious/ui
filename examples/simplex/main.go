package main

import (
	"bytes"
	"context"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/eliquious/ui"
	"github.com/hajimehoshi/ebiten/v2"
	simplex "github.com/ojrac/opensimplex-go"
)

func main() {
	ui.EnableHighDPI()

	ctx := context.Background()
	display := ui.New(ctx, &ui.DisplaySettings{
		Title:      "MQDB UI v0.1.0",
		Width:      1024,
		Height:     768,
		HideCursor: true,
		// BackgroundColor: color.RGBA{26, 27, 31, 255},
	})

	// screen size is dpubled due to high dpi
	display.SetBackground(NewTiledImageSimplexBackground(readTile(), 1024*2, 768*2))

	display.Add(ui.Rectangle(ui.Rect(0, 0, 2048, 768*2), &ui.RectangleOptions{
		FillColor: color.RGBA{0, 0, 0, 125},
	}))
	display.Add(ui.Rectangle(ui.Rect(32, 32, 512, 768*2-64), &ui.RectangleOptions{
		FillColor: color.RGBA{0, 0, 0, 64},
		Border:    ui.StrokeBorder(color.RGBA{100, 100, 100, 50}, 2),
	}))

	// display.Add(ui.FPSDisplay())
	display.SetCursor(ui.ImageCursor(readCursorTile(), true))
	if err := display.Show(); err != nil {
		log.Fatal(err)
	}
}

func readCursorTile() *ebiten.Image {

	// read file
	b, err := ioutil.ReadFile("cursor.png")
	if err != nil {
		panic(err)
	}

	// decode PNG file
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(img)
}

func readTile() *ebiten.Image {

	// read file
	b, err := ioutil.ReadFile("/Users/eliquious/Data/Sprite-0001.png")
	if err != nil {
		panic(err)
	}

	// decode PNG file
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(img)
}

// NewTiledImageSimplexBackground creates an animated background.
func NewTiledImageSimplexBackground(tile *ebiten.Image, w, h int) ui.Component {
	iw, ih := tile.Size()
	noiseTile := ebiten.NewImage(iw, ih)
	noiseTile.Fill(color.White)

	data := make([]float64, (h/ih)*(w/iw))
	s := &SimplexBackground{
		Tile:      tile,
		Width:     w,
		Height:    h,
		data:      data,
		noiseTile: noiseTile,
		iWidth:    iw,
		iHeight:   ih,
		dataCh:    make(chan []float64),
	}
	go s.calculateNoise()
	return s
}

// SimplexBackground is a simplex noise background
type SimplexBackground struct {
	Tile          *ebiten.Image
	Width, Height int

	noiseTile       *ebiten.Image
	iWidth, iHeight int
	data            []float64
	dataCh          chan []float64
}

func (s *SimplexBackground) calculateNoise() {
	ticker := time.NewTicker(time.Millisecond * 24)
	defer ticker.Stop()

	var tick float64
	noise := simplex.New(rand.Int63())

	data := make([]float64, (s.Height/s.iHeight)*(s.Width/s.iWidth))
	for {
		select {
		case <-ticker.C:
			tick += 0.006125

			var offset int
			for j := 0; j < s.Height; j += s.iHeight {
				for i := 0; i < s.Width; i += s.iWidth {
					data[offset] = noise.Eval3(float64(i/s.iWidth)/8, float64(j/s.iHeight)/8, tick) + 1/2
					offset++
				}
			}
			s.dataCh <- data
		}
	}
}

// Update updates the component
func (s *SimplexBackground) Update(ctx *ui.UpdateContext) error {
	select {
	case noise := <-s.dataCh:
		s.data = noise
	default:
	}
	return nil
}

// Display draws the component
func (s *SimplexBackground) Display(ctx *ui.DisplayContext) {
	op := &ebiten.DrawImageOptions{}
	sw, sh := ctx.Image().Size()

	// op := &ebiten.DrawTrianglesOptions{Filter: ebiten.FilterNearest}
	var val float64
	var offset int
	for j := 0; j < sh; j += s.iHeight {
		for i := 0; i < sw; i += s.iWidth {
			// val := noise.Eval3(float64(i/iw)/8, float64(j/ih)/8, tick) + 1/2
			val = s.data[offset]
			offset++

			// c := color.RGBA{0xff, 0xff, 0xff, uint8((0.02 + 0.08*val) * 128)}
			// vs, is := RectVertices(i, j, i+iw, j+ih, c)
			// ctx.DrawTriangles(vs, is, op)

			op.ColorM.Scale(1, 1, 1, 0.12+0.08*val)
			op.GeoM.Translate(float64(i), float64(j))
			ctx.DrawImage(s.noiseTile, op)
			op.ColorM.Reset()
			op.GeoM.Reset()

			op.ColorM.Scale(1, 1, 1, 0.25+val/2)
			op.GeoM.Translate(float64(i), float64(j))
			// op.GeoM.Scale(2, 2)
			ctx.DrawImage(s.Tile, op)
			op.GeoM.Reset()
			op.ColorM.Reset()

		}
	}
}

// RectVertices returns the vertices for a rectangle.
func RectVertices(px0, py0, px1, py1 int, c color.Color) ([]ebiten.Vertex, []uint16) {
	x0, y0 := float32(px0), float32(py0)
	x1, y1 := float32(px1), float32(py1)

	r0, g0, b0, a0 := c.RGBA()
	clr := color.RGBA{uint8(r0), uint8(g0), uint8(b0), uint8(a0)}

	r := float32(clr.R) / 0xff
	g := float32(clr.G) / 0xff
	b := float32(clr.B) / 0xff
	a := float32(clr.A) / 0xff

	return []ebiten.Vertex{
		{
			// bottom left
			DstX:   x0,
			DstY:   y1,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			// top left
			DstX:   x0,
			DstY:   y0,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			// top right
			DstX:   x1,
			DstY:   y0,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			// bottom right
			DstX:   x1,
			DstY:   y1,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	}, []uint16{0, 1, 2, 1, 2, 3}
}
