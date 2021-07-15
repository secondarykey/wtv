package wtv

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"wtv/config"

	"golang.org/x/xerrors"
)

type Book struct {
	dir      string
	files    []string
	optimize bool
}

func NewBook(dir string) (*Book, error) {

	var b Book
	b.optimize = false
	b.dir = dir
	files, err := getFiles(dir)
	if err != nil {
		return nil, xerrors.Errorf("getFiles() error: %w", err)
	}

	b.files = files
	return &b, nil
}

func (b *Book) Page() int {
	return len(b.files)
}

var BookIndexError = fmt.Errorf("Book Index Error")

func (b *Book) Load(idx int) (image.Image, error) {

	if idx < 0 || idx >= len(b.files) {
		return nil, BookIndexError
	}

	img, err := Load(b.files[idx])
	if err != nil {
		return nil, xerrors.Errorf("Load() error: %w", err)
	}
	return img, nil
}

func (b *Book) String() string {
	return fmt.Sprintf("%v", b.files)
}

func (b *Book) canOptimize(w, h int) bool {

	if b.optimize {
		return false
	}

	name := b.files[0]
	img, err := Load(name)
	if err != nil {
		return false
	}

	bo := img.Bounds()
	s := float64(w) / float64(bo.Dx())
	nowH := float64(bo.Dy()) * s

	if nowH > (OpenGLHeight / 4.0) {
		return true
	}
	return false
}

func (b *Book) Optimize(w, h int) (*Book, error) {

	path := filepath.Join(b.dir, OptimizeDirectory)
	err := os.Mkdir(path, 0777)
	if err != nil {
		return nil, xerrors.Errorf("already exitst?[%s]", path)
	}

	var newB Book
	newB.optimize = true
	newB.dir = path

	for _, name := range b.files {

		img, err := Load(name)
		if err != nil {
			return nil, xerrors.Errorf("Load() error: %w", err)
		}

		bou := img.Bounds()
		s := float64(w) / float64(bou.Dx())

		nowH := float64(bou.Dy()) * s

		div := int(nowH / OptimizeHeight)

		nn := filepath.Base(name)
		if idx := strings.LastIndex(nn, "."); idx != -1 {
			nn = nn[0:idx]
		}
		nameFmt := filepath.Join(path, nn+"_%02d.jpg")

		divH := bou.Dy() / div
		modH := bou.Dy() % div

		cutH := 0

		for idx := 0; idx < div; idx++ {

			flg := false
			if idx == div-1 {
				divH += modH
				flg = true
			}

			r := image.Rect(0, 0, bou.Dx(), divH)
			newImg := image.NewRGBA(r)
			draw.Draw(newImg, r, img, image.Point{0, cutH}, draw.Src)

			fn := fmt.Sprintf(nameFmt, idx)
			err := WriteImage(fn, newImg)
			if err != nil {
				return nil, xerrors.Errorf("WriteImage() error: %w", err)
			}

			newB.files = append(newB.files, fn)

			if flg {
				break
			}
			cutH += divH
		}
	}

	return &newB, nil
}

func getFiles(dir string) ([]string, error) {

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, xerrors.Errorf("os.ReadDir() error: %w", err)
	}

	var rtn []string
	for _, entry := range entries {

		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			files, err := getFiles(path)
			if err != nil {
				return nil, xerrors.Errorf("getFiles(r) error: %w", err)
			}
			rtn = append(rtn, files...)
		} else {
			rtn = append(rtn, path)
		}
	}

	conf := config.Get()
	sort.Slice(rtn, conf.Sort.Less(rtn))

	return rtn, nil
}

func existsOptimizeDirectory(dir string) bool {
	path := filepath.Join(dir, OptimizeDirectory)
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
