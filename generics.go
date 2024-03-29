package main

import (
	"strconv"
	"strings"
)

type CsvFormat interface {
	Id() string
	DbOutHdr() string
	IsMaxDB(int) bool
	MaxDB() int
	DoSkipLine(int) bool
	DoSkipCol(int) bool
}

type GenericCsvFormat struct {
	Ident      string
	Iprangecol int
	Dbrangemax []int
	Skiplines  []int
	Skipcols   []int
	Header     string
	Formatin   int
	RedisCMD   string
	Autodetect bool
	Status     bool
}

func (g GenericCsvFormat) Id() string {
	return g.Ident
}

func (g GenericCsvFormat) DbOutHdr() string {
	var header string
	if g.Autodetect == false {
		header = g.Header + strconv.Itoa(g.Formatin)
	} else {
		header = g.Header + "0"
	}
	return header
}

func (g GenericCsvFormat) IsMaxDB(m int) bool {
	if m < g.Dbrangemax[len(g.Dbrangemax)-1] {
		return false
	}
	return true
}

func (g GenericCsvFormat) MaxDB() int {
	return g.Dbrangemax[len(g.Dbrangemax)-1]
}

// DoSkip skip a line?
func (g GenericCsvFormat) DoSkipLine(s int) bool {
	for i := range g.Skiplines {
		if g.Skiplines[i] == s {
			return true
		}
	}
	return false
}

// DoSkip skip a column?
func (g GenericCsvFormat) DoSkipCol(s int) bool {
	for i := range g.Skipcols {
		if g.Skipcols[i] == s {
			return true
		}
	}
	return false
}

func ip2long(ip string) int {
	var data string
	pos := strings.Index(ip, "/")

	if !(pos == -1) {
		data = ip[0:pos]
	} else {
		data = ip
	}

	a := uint32(data[12])
	b := uint32(data[13])
	c := uint32(data[14])
	d := uint32(data[15])

	converted := int(a<<24 | b<<16 | c<<8 | d)

	return converted
}

// genericCsv is just a template
/*
func genericCsv(informat int, autodetect bool) (g GenericCsvFormat) {

	g = GenericCsvFormat{
		Ident:      "ip2location",
		Iprangecol: 0,
		Dbrangemax: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24},
		Skiplines:  []int{},
		Skipcols:   []int{1},
		Header:     "DB",
		Formatin:   informat,
		RedisCMD:   "ZADD",
		Autodetect: autodetect,
		Status:     true}

	if informat > g.Dbrangemax[len(g.Dbrangemax)-1] {
		g.Status = false
	}

	if autodetect == true {
		g.Autodetect = true
	}

	return g
}
*/
