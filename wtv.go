package wtv

import (
	"wtv/config"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/xerrors"
)

const (
	OptimizeDirectory = ".wtv_optimize"
	//height (65536) must be less than or equal to 32768
	//TODO  -2 means -1 is ebiten error
	OpenGLHeight   = 1<<(16-1) - 2
	OptimizeHeight = 1 << 11           //2048
	OptimizeLimit  = OpenGLHeight >> 2 // 9
)

func Show() error {

	err := config.Load()
	if err != nil {
		return xerrors.Errorf("config.Load() error: %w", err)
	}

	conf := config.Get()

	ebiten.SetWindowTitle("Webtoon Viewer")
	ebiten.SetWindowSize(conf.Width, conf.Height)
	ebiten.SetWindowResizable(true)

	p := NewPlayer()

	err = ebiten.RunGame(p)
	if err != nil {
		return xerrors.Errorf("ebiten.RunGame() error: %w", err)
	}
	return nil
}
