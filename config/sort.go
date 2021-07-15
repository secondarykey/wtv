package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

type SortType int

const (
	NumericSortAsc SortType = iota
	NumericSortDesc
	AlphamericSortAsc
	AlphamericSortDesc
	ModTimeSortAsc
	ModTimeSortDesc

	NumericSort    = NumericSortAsc
	AlphamericSort = AlphamericSortAsc
	ModTimeSort    = ModTimeSortAsc

	DoNotSort
)

func (t SortType) Order(v bool) bool {
	if t.Asc() {
		return v
	}
	return !v
}

func (t SortType) Asc() bool {
	if t == NumericSortAsc || t == AlphamericSortAsc ||
		t == ModTimeSortAsc {
		return true
	}
	return false
}

func (t SortType) IsNumeric() bool {
	if t == NumericSortAsc || t == NumericSortDesc {
		return true
	}
	return false
}

func (t SortType) IsAlphameric() bool {
	if t == AlphamericSortAsc || t == AlphamericSortDesc {
		return true
	}
	return false
}

func (t SortType) IsModTime() bool {
	if t == ModTimeSortAsc || t == ModTimeSortDesc {
		return true
	}
	return false
}

func (t SortType) Reverse() SortType {
	switch t {
	case NumericSortAsc:
		return NumericSortDesc
	case AlphamericSortAsc:
		return AlphamericSortDesc
	case ModTimeSortAsc:
		return ModTimeSortDesc
	case NumericSortDesc:
		return NumericSortAsc
	case AlphamericSortDesc:
		return AlphamericSortAsc
	case ModTimeSortDesc:
		return ModTimeSortAsc
	}
	return DoNotSort
}

func (t SortType) Less(src []string) func(i, j int) bool {
	if t.IsNumeric() {
		return t.sortNumeric(src)
	} else if t.IsAlphameric() {
		return t.sortAlphameric(src)
	} else if t.IsModTime() {
		return t.sortModTime(src)
	}
	return t.sortNone()
}

func (t SortType) sortNumeric(src []string) func(int, int) bool {
	return func(i, j int) bool {
		p1, err1 := parseNumberByPath(src[i])
		p2, err2 := parseNumberByPath(src[j])

		if err1 != nil && err2 != nil {
			return t.Order(src[i] < src[j])
		} else if err1 != nil {
			return t.Order(false)
		} else if err2 != nil {
			return t.Order(true)
		}

		return t.Order(p1 < p2)
	}
}

func (t SortType) sortAlphameric(src []string) func(int, int) bool {
	return func(i, j int) bool {
		p1 := src[i]
		p2 := src[j]
		return t.Order(p1 < p2)
	}
}

func (t SortType) sortModTime(src []string) func(int, int) bool {
	return func(i, j int) bool {
		p1 := src[i]
		p2 := src[j]
		info1, err := os.Stat(p1)
		if err != nil {
			return t.Order(false)
		}
		info2, err := os.Stat(p2)
		if err != nil {
			return t.Order(false)
		}

		if info1.ModTime().Sub(info2.ModTime()).Seconds() < 0 {
			return t.Order(true)
		}
		return t.Order(false)
	}
}

func (t SortType) sortNone() func(int, int) bool {
	return func(i, j int) bool {
		return i < j
	}
}

func parseNumberByPath(p string) (int, error) {

	name := filepath.Base(p)
	ne := name
	if idx := strings.Index(ne, "."); idx != -1 {
		ne = ne[0:idx]
	}

	n, err := strconv.Atoi(ne)
	if err != nil {
		return -1, xerrors.Errorf("strconv.Atoi() error: %w", err)
	}

	return n, nil
}
