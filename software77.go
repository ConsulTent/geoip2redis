package main

import (
	"strconv"
)

type Software77CSV struct {
	GenericCsvFormat
}

func (g Software77CSV) Id() string {
	return g.Ident
}

func (g Software77CSV) DbOutHdr() string {
	var header string

	header = g.Header + strconv.Itoa(g.Formatin)
	return header
}

func (g Software77CSV) IsMaxDB(m int) bool {
	if m < g.Dbrangemax[len(g.Dbrangemax)-1] {
		return false
	}
	return true
}

func (g Software77CSV) MaxDB() int {
	return g.Dbrangemax[len(g.Dbrangemax)-1]
}

// DoSkip skip a line?
func (g Software77CSV) DoSkipLine(s int) bool {
	for i := range g.Skiplines {
		if g.Skiplines[i] == s {
			return true
		}
	}
	return false
}

// DoSkip skip a column?
func (g Software77CSV) DoSkipCol(s int) bool {
	for i := range g.Skipcols {
		if g.Skipcols[i] == s {
			return true
		}
	}
	return false
}

// software77 we define All the Things for software77
func software77(informat int) (g Software77CSV) {

	g = Software77CSV{
		GenericCsvFormat: GenericCsvFormat{
			Ident:      "software77",
			Iprangecol: 0,
			Dbrangemax: []int{0, 1},
			Skiplines:  []int{},
			Skipcols:   []int{1, 3},
			Header:     "BRM",
			Formatin:   informat,
			RedisCMD:   "ZADD",
			Autodetect: false,
			Status:     true}}

	if informat > g.GenericCsvFormat.Dbrangemax[len(g.GenericCsvFormat.Dbrangemax)-1] {
		g.Status = false
	}

	if informat == 1 {
		g.GenericCsvFormat.Skipcols = []int{1, 2, 3, 5}
		g.GenericCsvFormat.Header = "DB"
	}

	return g
}
