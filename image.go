package wtv

import (
	"image"
	"os"

	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"

	"golang.org/x/xerrors"
)

func Load(name string) (image.Image, error) {

	_, err := os.Stat(name)
	if err != nil {
		return nil, xerrors.Errorf("os.Stat() error: %w", err)
	}

	f, err := os.Open(name)
	if err != nil {
		return nil, xerrors.Errorf("os.Open() error: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, xerrors.Errorf("image.Decode() error: %w", err)
	}

	return img, nil
}

func Scale(img image.Image, scale float64) image.Image {
	src := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, int(float64(src.Dx())*scale), int(float64(src.Dy())*scale)))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, src, draw.Over, nil)
	return dst
}

func WriteImage(name string, img image.Image) error {
	fp, err := os.Create(name)
	if err != nil {
		return xerrors.Errorf("os.Create() error: %w", err)
	}
	defer fp.Close()

	err = jpeg.Encode(fp, img, &jpeg.Options{Quality: 80})
	if err != nil {
		return xerrors.Errorf("jpeg.Encode() error: %w", err)
	}

	return nil
}
