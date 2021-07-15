package wtv

import (
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/xerrors"
)

var logger *log.Logger
var dw *DisplayWriter

func init() {
	dw = NewDisplayWriter()
	logger = log.New(dw, "", 0)
}

func setDebugDisplay(img *ebiten.Image) {
	dw.Draw(img)
}

type DisplayWriter struct {
	Display bool
	lines   []string
	img     *ebiten.Image
	builder *strings.Builder

	fontHeight int
}

func NewDisplayWriter() *DisplayWriter {
	var dw DisplayWriter
	dw.Display = true

	dw.fontHeight = defaultFont.Metrics().XHeight.Ceil()

	var builder strings.Builder
	dw.builder = &builder

	return &dw
}

func (w *DisplayWriter) Write(buf []byte) (int, error) {

	if !w.Display {
		return len(buf), nil
	}

	if w == nil {
		return 0, xerrors.Errorf("Writer or Writer image is nil")
	}

	//現在表示できる部分を表示
	return w.builder.Write(buf)
}

const (
	DisplayWriterMargin = 12
)

func (w *DisplayWriter) Draw(img *ebiten.Image) error {

	b := img.Bounds()
	if dw.img == nil || b.Dy() != dw.img.Bounds().Dy() {
		dw.img = ebiten.NewImage(b.Dx(), b.Dy())
	}

	dw.img.Clear()
	dw.img.Fill(color.RGBA{0, 0, 0, 30})

	buf := w.builder.String()
	lines := strings.Split(buf, "\n")

	h := b.Dy()
	dm := w.fontHeight + DisplayWriterMargin

	startY := 0
	startIdx := 0
	for leng := len(lines); (leng * dm) > h; leng-- {
		startIdx++
	}

	var writeLine []string
	for idx := startIdx; idx < len(lines); idx++ {
		line := lines[idx]
		writeLine = append(writeLine, line)
	}

	var builder strings.Builder
	for idx, txt := range writeLine {
		dy := (dm)*(idx+1) + startY
		text.Draw(w.img, txt, defaultFont, 10, dy, color.RGBA{23, 200, 0, 255})
		builder.WriteString(txt)
		if idx+1 != len(writeLine) {
			builder.WriteString("\n")
		}
	}

	dw.builder = &builder

	img.DrawImage(w.img, nil)
	return nil
}
