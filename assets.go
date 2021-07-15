package wtv

import (
	"embed"
	"image"
	"image/png"
	"io/fs"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/xerrors"
)

//go:embed _assets
var embAssets embed.FS
var assets fs.FS

type ResourceName string

const (
	ResFolder ResourceName = "baseline_folder_black_24dp.png"
	ResPlay   ResourceName = "outline_play_arrow_black_24dp.png"
	ResPause  ResourceName = "outline_pause_black_24dp.png"
	ResSwap   ResourceName = "outline_swap_vert_black_24dp.png"
	ResTop    ResourceName = "outline_vertical_align_top_black_24dp.png"
)

var imageNames = []ResourceName{ResFolder, ResPlay, ResPause, ResSwap, ResTop}

func (n ResourceName) Name() string {
	switch n {
	case ResFolder:
		return "Folder"
	case ResPlay:
		return "Play"
	case ResPause:
		return "Pause"
	case ResSwap:
		return "Swap"
	case ResTop:
		return "Top"
	default:
		panic("not found resource name")
	}
	return ""
}

var assetImages map[string]image.Image

func init() {

	var err error
	assets, err = fs.Sub(embAssets, "_assets")
	if err != nil {
		panic(err)
	}
	err = loadImages(assets)
	if err != nil {
		panic(err)
	}
	//clear
	embAssets = embed.FS{}

	err = initFont()
	if err != nil {
		panic(err)
	}

}

func loadImages(fsys fs.FS) error {
	assetImages = make(map[string]image.Image)
	for _, elm := range imageNames {
		f, err := fsys.Open(string(elm))
		if err != nil {
			return xerrors.Errorf("open() error: %w", err)
		}
		defer f.Close()

		img, err := png.Decode(f)
		if err != nil {
			return xerrors.Errorf("png.Decode() error: %w", err)
		}
		assetImages[elm.Name()] = img
	}
	return nil
}

func GetImage(res ResourceName) image.Image {
	return assetImages[res.Name()]
}

var defaultFont font.Face

const defaultFontName = "OdibeeSans-Regular.ttf"
const defaultFontDPI = 72
const defaultFontSize = 20

func initFont() error {

	f, err := assets.Open(defaultFontName)
	if err != nil {
		return xerrors.Errorf("assets.Open() error: %w", err)
	}

	info, err := f.Stat()
	if err != nil {
		return xerrors.Errorf("file.Stat() error: %w", err)
	}

	sz := info.Size()
	data := make([]byte, sz)
	_, err = f.Read(data)
	if err != nil {
		return xerrors.Errorf("file.Read() error: %w", err)
	}

	font, err := opentype.Parse(data)
	if err != nil {
		return xerrors.Errorf("file.Read() error: %w", err)
	}

	defaultFont, err = opentype.NewFace(font, &opentype.FaceOptions{
		Size: defaultFontSize,
		DPI:  defaultFontDPI,
	})
	return nil
}
