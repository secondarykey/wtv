package wtv

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DragState int

const (
	DragNoneState DragState = iota
	DragStartState
	DraggingState
	DragFinishState
)

func (d DragState) Get() DragState {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return DragStartState
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return DragFinishState
	}

	if d.started() {
		return DraggingState
	}
	return DragNoneState
}

func (d DragState) started() bool {
	if d == DragStartState || d == DraggingState {
		return true
	}
	return false
}
