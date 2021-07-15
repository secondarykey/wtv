package config

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/xerrors"
)

var gConf *Config

type Config struct {
	Directory string
	Direction Direction
	Effect    Effect
	FitMode   bool
	Width     int
	Height    int
	Sort      SortType
}

const (
	defaultConfigFileName = ".wtv_config_gob"
)

func init() {
	gob.Register(&Config{})
	gConf = defaultConfig()
}

func defaultConfig() *Config {
	var cnf Config
	cnf.Directory = ""
	cnf.Direction = Down
	cnf.Effect = Scroll
	cnf.Sort = NumericSort
	cnf.FitMode = true
	cnf.Width = 500
	cnf.Height = 800
	return &cnf
}

func Get() *Config {
	return gConf
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Effect int

const (
	Fadein Effect = iota
	Scroll
)

func Load() error {
	p := getPath()

	if _, err := os.Stat(p); err != nil {
		//TODO Create?
		return nil
	}

	err := gConf.load(p)
	if err != nil {
		return xerrors.Errorf("load() error: %w", err)
	}
	return nil
}

func (c *Config) load(name string) error {

	fp, err := os.Open(name)
	if err != nil {
		return xerrors.Errorf("os.Create() error: %w")
	}
	defer fp.Close()

	dec := gob.NewDecoder(fp)

	var cnf Config
	err = dec.Decode(&cnf)
	if err != nil {
		return xerrors.Errorf("Decode() error: %w")
	}

	gConf = &cnf

	return nil
}

func Save() error {

	p := getPath()
	err := gConf.save(p)
	if err != nil {
		return xerrors.Errorf("save() error: %w", err)
	}
	return nil
}

func (c *Config) save(name string) error {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(c)
	if err != nil {
		return xerrors.Errorf("Encode() error: %w")
	}

	fp, err := os.Create(name)
	if err != nil {
		return xerrors.Errorf("os.Create() error: %w")
	}
	defer fp.Close()

	_, err = fp.Write(buf.Bytes())
	if err != nil {
		return xerrors.Errorf("Write() error: %w")
	}

	return nil
}

func getPath() string {
	path := getHome()
	return filepath.Join(path, defaultConfigFileName)
}

func getHome() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	}
	return os.Getenv(env)
}
