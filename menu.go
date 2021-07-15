package wtv

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/xerrors"
)

type Direction int

const (
	N Direction = iota
	S
	E
	W
)

func (d Direction) Odds() int {
	if d == E || d == S {
		return -1
	}
	return 1
}

type Menu struct {
	direction Direction
	state     MenuState
	area      int
	limit     int

	relativeX int
	relativeY int

	move           int
	areaMovement   int
	activeMovement int

	x int
	y int

	img *ebiten.Image

	*Components
}

const (
	defaultAreaMovement   = 2
	defaultActiveMovement = 5
)

type MenuState int

const (
	MenuOFFState MenuState = iota
	MenuAreaState
	MenuActiveState
	MenuHideState
)

func (m MenuState) String() string {
	switch m {
	case MenuOFFState:
		return "OFF"
	case MenuAreaState:
		return "Area"
	case MenuActiveState:
		return "Active"
	case MenuHideState:
		return "Hide"
	}
	return "None"
}

//m.area zero is MenuAreaState
func NewMenu(d Direction, area, limit int) *Menu {
	var m Menu

	m.direction = d
	m.state = MenuOFFState
	m.area = area
	m.limit = limit
	m.areaMovement = defaultAreaMovement
	m.activeMovement = defaultActiveMovement
	m.img = nil
	//TODO Menu implement Component
	m.Components = NewComponents()

	return &m
}

func (m *Menu) updatePosition(w, h int) bool {

	m.relativeX = -1
	m.relativeY = -1
	nowX, nowY := ebiten.CursorPosition()

	rtn := false
	area := m.area
	if m.area == 0 || m.state == MenuActiveState || m.state == MenuHideState {
		area = m.limit
	}

	dw := w
	dh := h

	switch m.direction {
	case N:
		m.x = 0
		m.y = m.limit * -1
		dh = m.limit
		if nowY <= area {
			rtn = true
			if m.state == MenuActiveState {
				m.relativeX = nowX
				m.relativeY = nowY
			}
		}
	case S:
		m.x = 0
		m.y = h
		dh = m.limit
		if nowY >= h-area {
			rtn = true
			if m.state == MenuActiveState {
				m.relativeX = nowX
				m.relativeY = m.limit - (h - nowY)
			}
		}
	case W:
		m.x = m.limit * -1
		m.y = 0
		dw = m.limit
		if nowX <= area {
			rtn = true
			if m.state == MenuActiveState {
				m.relativeX = nowX
				m.relativeY = nowY
			}
		}
	case E:
		m.x = w
		m.y = 0
		dw = m.limit
		if nowX >= w-area {
			rtn = true
			if m.state == MenuActiveState {
				m.relativeX = w - nowX
				m.relativeY = nowY
			}
		}
	}

	m.img = ebiten.NewImage(dw, dh)
	m.img.Fill(color.RGBA{0, 0, 0, 255})

	if m.area != 0 &&
		(m.state == MenuAreaState || (m.state == MenuActiveState && m.move != m.limit)) {
		center := w / 2
		leng := 80
		uy := m.limit - m.area + 5
		dy := m.limit - 5
		vertecies := []ebiten.Vertex{
			newVertex(center-(leng/2), uy), newVertex(center+(leng/2), uy), newVertex(center, dy),
		}
		indices := []uint16{0, 1, 2}
		white := ebiten.NewImage(10, 10)
		white.Fill(color.RGBA{200, 200, 200, 255})
		m.img.DrawTriangles(vertecies, indices, white, nil)
	}

	if !rtn {
		if m.state == MenuActiveState {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				m.state = MenuHideState
			} else if m.area == 0 {
				m.state = MenuHideState
			}
		}
	}

	if nowX <= 0 || nowX >= w {
		return false
	}
	if nowY <= 0 || nowY >= h {
		return false
	}
	return rtn
}

func newVertex(x, y int) ebiten.Vertex {
	return ebiten.Vertex{
		DstX: float32(x), DstY: float32(y),
		ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1,
	}
}

func (m *Menu) Update(w, h int) error {

	ebiten.SetCursorShape(ebiten.CursorShapeDefault)

	area := m.updatePosition(w, h)
	if !area {
		if m.state == MenuAreaState {
			m.state = MenuHideState
		}
	}

	//TODO event作り直しが必要
	err := m.Components.Update(m.relativeX, m.relativeY)
	if err != nil {
		return xerrors.Errorf("Components.Update() error: %w", err)
	}

	switch m.state {
	case MenuOFFState:
		m.move = 0
		if area {
			m.state = MenuAreaState
			if m.area == 0 {
				m.state = MenuActiveState
			}
		}
	case MenuAreaState:
		m.move += m.areaMovement
		if m.move > m.area {
			m.move = m.area
		}
		//TODO 中央？
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			m.state = MenuActiveState
		}

	case MenuHideState:
		m.move -= m.activeMovement
		if m.move <= 0 {
			m.move = 0
			m.state = MenuOFFState
		}
	case MenuActiveState:
		m.move += m.activeMovement
		if m.move > m.limit {
			m.move = m.limit
		}
	}

	if m.x != 0 {
		m.x += m.move * m.direction.Odds()
	}
	if m.y != 0 {
		m.y += m.move * m.direction.Odds()
	}

	return nil
}

func (m *Menu) Active() bool {
	return m.state != MenuOFFState
}

func (m *Menu) Draw(img *ebiten.Image) error {
	var err error
	switch m.state {
	case MenuAreaState:
		err = m.drawArea(img)
	case MenuActiveState, MenuHideState:
		err = m.drawActive(img)
	}
	if err != nil {
		return xerrors.Errorf("draw error: %w", err)
	}
	return nil
}

func (m *Menu) drawArea(img *ebiten.Image) error {

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(m.x), float64(m.y))
	op.ColorM.Scale(1, 1, 1, 0.5)

	img.DrawImage(m.img, op)

	return nil
}

func (m *Menu) drawActive(img *ebiten.Image) error {

	err := m.Components.Draw(m.img)
	if err != nil {
		return xerrors.Errorf("Components.Draw() error: %w", err)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(m.x), float64(m.y))
	op.ColorM.Scale(1, 1, 1, 0.9)
	img.DrawImage(m.img, op)
	return nil
}
