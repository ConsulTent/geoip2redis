package ip2location

import (
	"fmt"
)

var DBKEY string

type csvdata struct{}

type CSVdataDB1 struct {
	ipFROM       uint   `csv:"ipfrom"`
	ipTO         uint   `csv:"-"`
	countrySHORT string `csv:"countrycode"`
	countryLONG  string `csv:"countryname"`
}

func ip2location(format int) *csvdata {
	switch format {
	case 1:
		DBKEY = "DB1"
		csvdata := new(CSVdataDB1)
	case 11:
		DBKEY = "DB11"
	default:
		fmt.Println("ip2location format unknown")
	}

	return csvdata
}
