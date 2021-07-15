package wtv

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Slider struct {
	max        int
	current    int
	changeFunc func(v int) error

	rect image.Rectangle
}

func NewSlider() *Slider {
	var s Slider
	s.current = 1
	s.max = 1

	s.rect = image.Rect(int(SliderStartX), int(SliderCurrentY),
		int(SliderWidth)+int(SliderStartX), int(SliderCurrentHeight)+int(SliderCurrentY))
	return &s
}

func (s *Slider) Changed(fn func(v int) error) {
	s.changeFunc = fn
}

func (s *Slider) SetMax(max int) {
	s.max = max
}

func (s *Slider) Set(w, h int) {
	//TODO Not Implemented
}

func (s *Slider) SetValue(current int) {
	s.current = current
}

func (s *Slider) GetValue() int {
	return s.current
}

func (s *Slider) Move(x, y int) error {
	return fmt.Errorf("Not Implemented")
}

func (s *Slider) Point() (float64, float64) {
	log.Println("Slider Point Not Implemented")
	return -1, -1
}

func (s *Slider) In(x, y int) bool {
	return image.Pt(x, y).In(s.rect)
}

func (s *Slider) Update(x, y int) error {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if s.In(x, y) {
			v := int((float64(x)-SliderStartX)/(SliderWidth-SliderCurrentWidth)*(float64(s.max)-1.0)) + 1
			s.changeFunc(v)
			s.current = v
		}
	}
	return nil
}

const (
	SliderStartX        = 20.0
	SliderWidth         = 500.0
	SliderCurrentY      = 30.0
	SliderCurrentWidth  = 10.0
	SliderCurrentHeight = 15.0
)

func (s *Slider) Draw(img *ebiten.Image) error {
	x, y, w, h := SliderStartX, SliderCurrentY+5.0, SliderWidth, 5.0
	ebitenutil.DrawRect(img, x, y, w, h, color.White)

	p := 0.0
	if s.max != 1 {
		p = float64(s.current-1) / float64(s.max-1)
	}
	currentP := (SliderWidth-SliderCurrentWidth)*p + SliderStartX

	cx, cy, cw, ch := currentP, SliderCurrentY, SliderCurrentWidth, SliderCurrentHeight
	ebitenutil.DrawRect(img, cx, cy, cw, ch, color.White)

	tx := SliderStartX + SliderWidth + 10

	text.Draw(img, fmt.Sprintf("%d/%d", s.current, s.max), defaultFont,
		int(tx), SliderCurrentY+SliderCurrentHeight, color.White)

	return nil
}
