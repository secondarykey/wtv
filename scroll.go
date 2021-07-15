package wtv

import (
	"image"
	"image/draw"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/xerrors"
)

type ScrollMenu struct {
	loading sync.Once

	index     int
	start     int
	scrollMap map[int]int

	selectedIndex int
	selectedPos   int

	img *ebiten.Image

	*Menu
}

func NewScrollMenu(m *Menu) *ScrollMenu {
	var sm ScrollMenu
	sm.Menu = m
	sm.selectedIndex = -1
	sm.selectedPos = -1
	return &sm
}

func (sm *ScrollMenu) Load(b *Book, idx int, width, pos int) error {
	sm.loading.Do(func() {
		sm.load(b, idx, width, pos)
	})
	return nil
}

func (sm *ScrollMenu) load(b *Book, idx int, width, pos int) error {

	mib := sm.Menu.img.Bounds()
	w := mib.Dx()
	h := mib.Dy()

	sm.scrollMap = make(map[int]int)

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	img1, err := b.Load(idx)
	if err != nil {
		return xerrors.Errorf("Book Load() error: %w", err)
	}

	half := h / 2
	b1 := img1.Bounds()

	//現在表示中の倍率を取得
	orgS := float64(width) / float64(b1.Dx())
	//画面サイズの高さを実際の中央位置を算出
	orgP := float64(pos+half) / orgS

	//元の画像での座標を算出
	scale := float64(w) / float64(b1.Dx())
	p1 := Scale(img1, scale)
	pb1 := p1.Bounds()

	p := int(orgP * scale)

	//現状のpを半分の位置に表示
	startY := half - p
	sm.index = idx

	sm.start = startY

	draw.Draw(
		img,
		image.Rect(0, startY, pb1.Dx(), pb1.Dy()+startY),
		p1, image.Point{0, 0}, draw.Over)
	sm.scrollMap[idx] = pb1.Dy()

	y := startY
	for i := idx - 1; i >= 0; i-- {
		prev, err := b.Load(i)
		if err != nil {
			return xerrors.Errorf("book preb load() error: %w", err)
		}
		prevB := prev.Bounds()
		prevS := float64(w) / float64(prevB.Dx())

		ps := Scale(prev, prevS)

		y = y - ps.Bounds().Dy()
		py := 0

		dy := ps.Bounds().Dy()
		sm.scrollMap[i] = dy

		if y < 0 {
			dy = dy + y
			py = y * -1
			y = 0
		}

		draw.Draw(
			img,
			image.Rect(0, y, prevB.Dx(), y+dy),
			ps, image.Point{0, py}, draw.Over)

		if y <= 0 {
			break
		}
	}

	y = startY + pb1.Dy()
	for i := idx + 1; i < b.Page(); i++ {
		next, err := b.Load(i)
		if err != nil {
			return xerrors.Errorf("book preb load() error: %w", err)
		}
		nextB := next.Bounds()
		nextS := float64(w) / float64(nextB.Dx())

		ns := Scale(next, nextS)

		dy := ns.Bounds().Dy()
		sm.scrollMap[i] = dy

		if y+dy > h {
			dy = h - y
		}

		draw.Draw(
			img,
			image.Rect(0, y, nextB.Dx(), y+dy),
			ns, image.Point{0, 0}, draw.Over)

		if y >= h {
			break
		}
		y = y + dy
	}

	sm.img = ebiten.NewImageFromImage(img)
	return nil
}

func (sm *ScrollMenu) Reset() {
	if sm == nil {
		return
	}

	if sm.img != nil {
		sm.loading = sync.Once{}
		sm.img = nil
		sm.selectedIndex = -1
		sm.selectedPos = -1
		sm.scrollMap = make(map[int]int)
	}
}

func (sm *ScrollMenu) Update(w, h int) error {

	err := sm.Menu.Update(w, h)
	if err != nil {
		return xerrors.Errorf("Update() error: %w", err)
	}

	if sm.Active() {

		ebiten.SetCursorShape(ebiten.CursorShapePointer)

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

			half := h / 2
			updown := sm.Menu.relativeY - half
			extent := 1
			if updown <= 0 {
				extent = -1
			}

			idx := sm.index
			pos := 0

			for {

				l, ok := sm.scrollMap[idx]
				if !ok {
					break
				}

				s := h / l

				updown = updown + l*extent*-1
				if extent <= 0 {
					if updown > 0 {
						logger.Println(s, updown)
						pos = updown*s - h
						break
					}
				} else {
					if updown < 0 {
						pos = updown*s - h
						break
					}
				}

				idx = idx + extent
			}

			sm.selectedIndex = idx
			sm.selectedPos = pos

			logger.Println(idx, pos)
		}
	}

	return nil
}

func (sm *ScrollMenu) Active() bool {
	return sm.Menu.Active()
}

func (sm *ScrollMenu) Draw(img *ebiten.Image) {
	if sm.Active() {
		if sm.img != nil {
			op := &ebiten.DrawImageOptions{}
			sm.Menu.img.DrawImage(sm.img, op)
		}
		sm.Menu.Draw(img)
	}
}
