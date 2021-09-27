package main

type Ip2LocationCSV struct {
	GenericCsvFormat
}

// ip2location we define All the Things for ip2location
func ip2location(informat int, autodetect bool) (g Ip2LocationCSV) {

	g = Ip2LocationCSV{
		GenericCsvFormat: GenericCsvFormat{
			Ident:      "ip2location",
			Iprangecol: 0,
			Dbrangemax: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24},
			Skiplines:  []int{},
			Skipcols:   []int{1},
			Header:     "DB",
			Formatin:   informat,
			RedisCMD:   "ZADD",
			Autodetect: autodetect,
			Status:     true}}

	if informat > g.GenericCsvFormat.Dbrangemax[len(g.GenericCsvFormat.Dbrangemax)-1] {
		g.GenericCsvFormat.Status = false
	}

	if autodetect {
		g.GenericCsvFormat.Autodetect = true
	}

	return g
}
