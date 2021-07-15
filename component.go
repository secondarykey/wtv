package wtv

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/xerrors"
)

type Component interface {
	Shape
	Layout
	Update(int, int) error
	Draw(*ebiten.Image) error
}

type ComponentObserver struct {
	Component
}

func NewComponent(c Component) *ComponentObserver {
	var o ComponentObserver
	o.Component = c
	return &o
}

type Components struct {
	parent   Component
	children []Component
}

func NewComponents() *Components {
	var c Components
	//TODO parent
	return &c
}

func (c *Components) Add(comp Component) {
	c.children = append(c.children, comp)
}

func (c *Components) Update(x, y int) error {
	for idx, comp := range c.children {
		err := comp.Update(x, y)
		if err != nil {
			return xerrors.Errorf("Components[%d Update() error: %w", idx, err)
		}
	}
	return nil
}

func (c *Components) Draw(img *ebiten.Image) error {
	for idx, comp := range c.children {
		err := comp.Draw(img)
		if err != nil {
			return xerrors.Errorf("Components[%d Draw() error: %w", idx, err)
		}
	}
	return nil
}
