package main


import (
	"strconv"
)

type Ip2LocationCSV struct {
	GenericCsvFormat
}


func (g Ip2LocationCSV) Id() string {
	return g.Ident
}

func (g Ip2LocationCSV) DbOutHdr() string {
	var header string
	if g.Autodetect == false {
	header = g.Header + strconv.Itoa(g.Formatin)
} else {
	header = g.Header + "0"
}
	return header
}

func (g Ip2LocationCSV) IsMaxDB(m int) bool {
	if m < g.Dbrangemax[len(g.Dbrangemax)-1] {
		return false
	}
	return true
}

func (g Ip2LocationCSV) MaxDB() int {
	return g.Dbrangemax[len(g.Dbrangemax)-1]
}

// DoSkip skip a line?
func (g Ip2LocationCSV) DoSkipLine(s int) bool {
	for i := range g.Skiplines {
		if g.Skiplines[i] == s {
			return true
		}
	}
	return false
}

// DoSkip skip a column?
func (g Ip2LocationCSV) DoSkipCol(s int) bool {
	for i := range g.Skipcols {
		if g.Skipcols[i] == s {
			return true
		}
	}
	return false
}


// ip2location we define All the Things for ip2location
func ip2location(informat int,autodetect bool) (g Ip2LocationCSV) {

	g = Ip2LocationCSV {
		GenericCsvFormat: GenericCsvFormat {
		Ident: "ip2location",
		Iprangecol: 0,
		Dbrangemax: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14 ,15 ,16, 17, 18, 19, 20, 21, 22, 23, 24},
		Skiplines: []int{ },
		Skipcols: []int{ 1 },
		Header: "DB",
		Formatin: informat,
	 	RedisCMD: "ZADD",
		Autodetect: autodetect,
	  Status: true }}


			if informat > g.GenericCsvFormat.Dbrangemax[len(g.GenericCsvFormat.Dbrangemax)-1] {
				g.GenericCsvFormat.Status = false
			}

			if autodetect == true {
				g.GenericCsvFormat.Autodetect = true
			}

	return g
}
