package wtv

import (
	"errors"
	"log"
	"wtv/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sqweek/dialog"
	"golang.org/x/xerrors"
)

type Player struct {
	width  int
	height int

	viewRedraw bool

	viewer *Viewer

	topMenu      *Menu
	scrollMenu   *ScrollMenu
	controllMenu *Menu
}

func NewPlayer() *Player {

	var p Player
	p.width = 0
	p.height = 0
	p.viewRedraw = true
	p.viewer = NewViewer()

	p.topMenu = NewMenu(N, 30, 110)

	sortBtn1 := NewTextButton("Numeric", 100, 10, 90, 30)
	sortBtn2 := NewTextButton("Alphanumeric", 200, 10, 90, 30)
	sortBtn3 := NewTextButton("Modtime", 300, 10, 90, 30)
	sortBtn1.Click(func() error {
		err := changeSortConfig(config.NumericSort)
		if err != nil {
			return xerrors.Errorf("changeSortConfig() numeric error: %w", err)
		}
		p.viewRedraw = true
		p.viewer.reset()
		return nil
	})
	sortBtn2.Click(func() error {
		err := changeSortConfig(config.AlphamericSort)
		if err != nil {
			return xerrors.Errorf("changeSortConfig() numeric error: %w", err)
		}
		p.viewRedraw = true
		p.viewer.reset()
		return nil
	})
	sortBtn3.Click(func() error {
		err := changeSortConfig(config.ModTimeSort)
		if err != nil {
			return xerrors.Errorf("changeSortConfig() numeric error: %w", err)
		}
		p.viewRedraw = true
		p.viewer.reset()
		return nil
	})

	autoBtn := NewCircleButton(500, 45, 32)
	autoBtn.Click(func() error {
		if p.viewer.playMode == AutoPlayMode {
			p.viewer.playMode = NormalPlayMode
			autoBtn.PasteImage(ResPlay)
		} else {
			p.viewer.playMode = AutoPlayMode
			autoBtn.PasteImage(ResPause)
		}
		p.topMenu.state = MenuHideState
		return nil
	})
	autoBtn.PasteImage(ResPlay)

	btn := NewCircleButton(50, 45, 32)
	btn.PasteImage(ResFolder)

	slider := NewSlider()

	btn.Click(func() error {

		conf := config.Get()

		t := "Load Webtoon Directory"
		dir := conf.Directory

		builder := dialog.Directory().Title(t)
		builder.StartDir = dir
		dir, err := builder.Browse()
		if err != nil {
			if errors.Is(err, dialog.ErrCancelled) {
				return nil
			}
			return xerrors.Errorf("dialog.Directory() error: %w", err)
		}

		p.viewer.SetBook(dir)
		p.topMenu.state = MenuHideState
		p.viewRedraw = true

		slider.SetMax(len(p.viewer.book.files))
		slider.SetValue(1)

		conf.Directory = dir
		err = config.Save()
		if err != nil {
			return xerrors.Errorf("config.Save() error: %w", err)
		}
		return nil
	})

	slider.Changed(func(v int) error {

		p.viewer.reset()
		p.viewer.index = v - 1
		p.viewer.pos = 0
		p.viewRedraw = true

		return nil
	})

	p.topMenu.Add(btn)
	p.topMenu.Add(sortBtn1)
	p.topMenu.Add(sortBtn2)
	p.topMenu.Add(sortBtn3)
	p.topMenu.Add(autoBtn)

	p.topMenu.state = MenuActiveState

	m := NewMenu(E, 0, 100)
	p.scrollMenu = NewScrollMenu(m)

	cm := NewMenu(S, 0, 80)

	cm.Add(slider)
	p.controllMenu = cm

	return &p
}

func changeSortConfig(t config.SortType) error {
	conf := config.Get()
	conf.Sort = t
	//Save is redraw
	return nil
}

func (p *Player) Layout(w, h int) (int, int) {

	//TODO 変更になるものに通知する
	// Resizable ?
	// Replaceable ?

	if p.width != w || p.height != h {
		p.width = w
		p.height = h
		p.viewRedraw = true
	} else if p.viewRedraw {
		if p.isView() {
			p.viewer.Redraw(w, h)

			conf := config.Get()
			conf.Width = w
			conf.Height = h
			err := config.Save()
			if err != nil {
				log.Println(err)
			}
		} else {
			p.viewer.width = w
			p.viewer.height = h
		}
		p.viewRedraw = false
	}

	return w, h
}

func (p *Player) Update() error {

	// TODO Updateが必要かどうか？

	if !p.viewer.Dragging() {
		if !p.topMenu.Active() && !p.controllMenu.Active() {

			p.scrollMenu.Update(p.width, p.height)
			idx := p.scrollMenu.selectedIndex
			if idx != -1 {
				p.viewer.reset()
				p.viewer.index = idx
				p.viewer.pos = p.scrollMenu.selectedPos
				p.viewer.Redraw(p.width, p.height)
				p.scrollMenu.selectedIndex = -1
				p.scrollMenu.selectedPos = -1
				p.scrollMenu.state = MenuHideState

				//TODO Slider
			}
		}

		if !p.scrollMenu.Active() && !p.controllMenu.Active() {
			p.topMenu.Update(p.width, p.height)
		}

		if !p.scrollMenu.Active() && !p.topMenu.Active() {
			for _, comp := range p.controllMenu.Components.children {
				if v, ok := comp.(*Slider); ok {
					v.SetValue(p.viewer.index + 1)
				}
			}
			p.controllMenu.Update(p.width, p.height)
		}
	}

	if p.scrollMenu.Active() {

		book, idx := p.viewer.GetBook()
		if book != nil {
			err := p.scrollMenu.Load(book, idx, p.viewer.width, p.viewer.pos)
			if err != nil {
				return xerrors.Errorf("ScrollMenu Load() error: %w")
			}
		}
		return nil
	} else {
		p.scrollMenu.Reset()
	}

	if p.topMenu.Active() {
		return nil
	}
	if p.scrollMenu.Active() {
		return nil
	}

	if p.isView() {
		err := p.viewer.Update()
		if err != nil {
			return xerrors.Errorf("viewer Update() error: %w", err)
		}
	}
	return nil
}

func (p *Player) Draw(screen *ebiten.Image) {

	if p.isView() {
		p.viewer.Draw(screen)
	}

	p.topMenu.Draw(screen)
	if p.isView() {
		p.scrollMenu.Draw(screen)
	}

	p.controllMenu.Draw(screen)

	setDebugDisplay(screen)
}

func (p *Player) isView() bool {
	if p.viewer.book == nil {
		return false
	}
	return true
}
