package wtv

import (
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/xerrors"
)

type Button interface {
	Component
}

var buttonColor = color.RGBA{128, 128, 128, 255}

type ButtonObserver struct {
	img   *ebiten.Image
	focus bool
	click func() error

	Button
}

func NewButton(b Button) *ButtonObserver {
	var bo ButtonObserver
	bo.Button = b
	return &bo
}

func (bo *ButtonObserver) Click(fn func() error) {
	bo.click = fn
}

func (bo *ButtonObserver) Update(x, y int) error {

	if bo == nil {
		return nil
	}
	bo.focus = false

	if bo.Button.In(x, y) {
		bo.focus = true
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			err := bo.click()
			if err != nil {
				return xerrors.Errorf("click() error: %w", err)
			}
		}
	}
	return nil
}

func (bo *ButtonObserver) Draw(img *ebiten.Image) error {

	if bo == nil || bo.img == nil {
		return nil
	}

	cop := &ebiten.DrawImageOptions{}

	cop.GeoM.Translate(bo.Point())

	if bo.focus {
		cop.ColorM.Scale(1, 1, 1, 0.8)
	}
	img.DrawImage(bo.img, cop)

	return nil
}

type RectButton struct {
	Shape
	*ButtonObserver
}

func NewRectButton(x, y int) *RectButton {
	var r RectButton
	r.ButtonObserver = NewButton(&r)
	return &r
}

type TextButton struct {
	*RectButton
}

func NewTextButton(txt string, x, y, w, h int) *TextButton {
	var t TextButton

	t.RectButton = NewRectButton(x, y)
	btn := ebiten.NewImage(w, h)
	btn.Fill(buttonColor)

	tw := font.MeasureString(defaultFont, txt).Ceil()
	th := defaultFont.Metrics().XHeight.Ceil()

	dx := (w / 2) - (tw / 2)
	dy := (h / 2) + (th / 2) + 2

	t.Shape = NewRectangle(x, y, w, h)
	text.Draw(btn, txt, defaultFont, dx, dy, color.Black)

	t.img = btn

	return &t
}

type CircleButton struct {
	*ButtonObserver
	Shape
}

func NewCircleButton(x, y, r int) *CircleButton {

	var c CircleButton
	c.ButtonObserver = NewButton(&c)
	c.Shape = NewCircle(x, y, r)

	c.setBaseImage()

	return &c
}

func (c *CircleButton) setBaseImage() {

	cs, ok := c.Shape.(*Circle)
	if !ok {
		log.Println("Shape Cast(Circle) error")
	}

	r := cs.r
	rc, gc, bc, _ := buttonColor.RGBA()

	dc := gg.NewContext(r*2, r*2)
	dc.DrawCircle(float64(r), float64(r), float64(r))
	dc.SetRGB(float64(rc), float64(gc), float64(bc))
	dc.Fill()

	img := ebiten.NewImageFromImage(dc.Image())
	c.img = img
}

func (c *CircleButton) PasteImage(name ResourceName) {

	c.setBaseImage()

	res := ebiten.NewImageFromImage(GetImage(name))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(8), float64(8))
	c.img.DrawImage(res, op)
}
