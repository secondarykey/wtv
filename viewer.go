package wtv

import (
	"errors"
	"fmt"
	"image"
	"log"
	"path/filepath"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/xerrors"
)

type PlayMode int

const (
	NormalPlayMode PlayMode = iota
	AutoPlayMode
)

type Viewer struct {
	book *Book

	prev    *ebiten.Image
	current *ebiten.Image
	next    *ebiten.Image
	index   int

	loadingPrev sync.Once
	loadingNext sync.Once

	playMode  PlayMode
	dragState DragState
	pos       int
	startPos  int

	width  int
	height int
}

func NewViewer() *Viewer {
	var v Viewer
	v.playMode = NormalPlayMode
	return &v
}

func (v *Viewer) SetBook(dir string) error {

	opti := false
	if existsOptimizeDirectory(dir) {
		dir = filepath.Join(dir, OptimizeDirectory)
		opti = true
	}

	b, err := NewBook(dir)
	if err != nil {
		return xerrors.Errorf("NewBook() error: %w", err)
	}
	b.optimize = opti

	if b.canOptimize(v.width, v.height) {
		nb, err := b.Optimize(v.width, v.height)
		if err != nil {
			return xerrors.Errorf("Book optimize error: %w", err)
		}
		b = nb
	}

	v.book = b

	err = v.reset()
	if err != nil {
		return xerrors.Errorf("reset() error: %w", err)
	}

	//v.playMode = AutoPlayMode
	return nil
}

func (v *Viewer) GetBook() (*Book, int) {
	return v.book, v.index
}

func (v *Viewer) reset() error {
	v.prev = nil
	v.current = nil
	v.next = nil
	v.index = 0
	v.pos = 0
	v.loadingNext = sync.Once{}
	v.loadingPrev = sync.Once{}
	return nil
}

func (v *Viewer) Dragging() bool {
	return v.dragState == DraggingState
}

func (v *Viewer) enable() bool {
	if v.book == nil {
		return false
	}
	if v.current == nil {
		return false
	}
	return true
}

func (v *Viewer) Redraw(w, h int) {

	img, err := v.book.Load(v.index)
	if err != nil {
		log.Println(err)
		return
	}

	v.width, v.height = w, h
	v.current = v.resize(img)

	//TODO  next prev
	if v.prev != nil {
		v.prev = v.resize(v.prev)
	}
	if v.next != nil {
		v.next = v.resize(v.next)
	}

	return
}

func (v *Viewer) resize(src image.Image) *ebiten.Image {
	s := float64(v.width) / float64(src.Bounds().Dx())
	h := float64(src.Bounds().Dy())

	if (h * s) > OpenGLHeight {
		orgS := s
		s = float64(OpenGLHeight) / float64(h)
		fmt.Printf("Due to height restrictions,the magnification will be changed\n%0.2f -> %0.2f\n", orgS, s)
	}

	dst := Scale(src, s)

	return ebiten.NewImageFromImage(dst)
}

func (v *Viewer) Update() error {

	if !v.enable() {
		return nil
	}

	v.loadPrev()
	v.loadNext()

	if v.playMode == AutoPlayMode {
		v.pos += 5
		fmt.Printf("\r%10d", v.pos)
		return nil
	}

	_, nowY := ebiten.CursorPosition()

	b := v.current.Bounds()
	_, dy := ebiten.Wheel()

	v.dragState = v.dragState.Get()

	if v.dragState == DragStartState {
		v.startPos = nowY
	} else if v.dragState == DraggingState || dy != 0 {

		my := v.startPos - nowY
		if dy != 0 {
			my = int(dy * 80 * -1)
		}

		v.pos, v.startPos = v.pos+my, nowY
		m := 0

		if v.pos < m {
			if v.prev == nil {
				v.pos = 0
			}
		} else if v.pos > b.Dy()-v.height-m {
			if v.next == nil {
				v.pos = b.Dy() - v.height
			}
		}
	}

	return nil
}

func (v *Viewer) loadPrev() error {

	if v.prev != nil {
		return nil
	}

	go v.loadingPrev.Do(func() {

		img, err := v.book.Load(v.index - 1)
		if err != nil {
			if !errors.Is(err, BookIndexError) {
				log.Println(err)
			}
			return
		}
		v.prev = v.resize(img)
	})
	return nil
}

func (v *Viewer) loadNext() error {

	if v.next != nil {
		return nil
	}

	go v.loadingNext.Do(func() {

		img, err := v.book.Load(v.index + 1)
		if err != nil {
			if !errors.Is(err, BookIndexError) {
				log.Println(err)
			}
			return
		}
		v.next = v.resize(img)
	})

	return nil
}

func (v *Viewer) Draw(screen *ebiten.Image) {

	if !v.enable() {
		return
	}

	op := &ebiten.DrawImageOptions{}
	py := float64(v.pos * -1)
	op.GeoM.Translate(0, py)
	screen.DrawImage(v.current, op)

	v.drawPrev(screen, py)
	v.drawNext(screen, py)
}

func (v *Viewer) drawPrev(screen *ebiten.Image, by float64) {

	if v.prev == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	iy := v.prev.Bounds().Dy() * -1
	ty := iy + int(by)
	op.GeoM.Translate(0, float64(ty))
	screen.DrawImage(v.prev, op)

	if ty > iy+v.height {
		go func() {
			v.current, v.next, v.prev = v.prev, v.current, nil
			v.pos = ty * -1
			v.index--
			v.loadingPrev = sync.Once{}
		}()
	}
}

func (v *Viewer) drawNext(screen *ebiten.Image, by float64) {

	if v.next == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	ty := v.current.Bounds().Dy() + int(by)
	op.GeoM.Translate(0, float64(ty))
	screen.DrawImage(v.next, op)

	if ty < 0 {
		go func() {
			v.current, v.prev, v.next = v.next, v.current, nil
			v.pos = ty * -1
			v.index++
			v.loadingNext = sync.Once{}
		}()
	}
}
