package main

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/rburmorrison/go-argue"
	"os"
)

type cmdline struct {
	CsvFile        string `options:"required,positional" help:"CSV file with GeoIP data"`
	Format         string `init:"f" options:"required" help:"ip2location|maxmind"`
	RedisHost      string `init:"r" help:"Redis Host, default 127.0.0.1"`
	RedisPort      int    `init:"p" help:"Redis Port, default 6379"`
	InPrecision    int    `init:"i" options:"required" help:"Input precision.  This would be db file number.  See README.TXT"`
	OutPrecision   int    `init:"o" help:"Output Precision.  Default: 0 (match input), see README.TXT"`
	DontSkipHeader bool   `init:"s" help:"DON'T Skip the first CSV line. Default: skip, see README.TXT"`
}

const pver = "0.0.1"

var gitver = "undefined"

type csvstruct []struct{}

var skiphdr bool

var DEBUG = true

func main() {

	var i uint64
	var x uint8
	var skipcol uint8
	var rediscmd string
	var cmds cmdline
	var samples [][]string
	var DBHDR string
	var index int
	//      var fakedata struct{}

	fmt.Printf("GeoIP2Redis (c) 2019 ConsulTent Ltd. v%s-%s\n", pver, gitver)

	agmt := argue.NewEmptyArgumentFromStruct(&cmds)

	agmt.Dispute(true)

	switch cmds.Format {
	case "maxmind":
		fmt.Println("maxmind is not supported yet")
		os.Exit(1)
	case "ip2location":
		DBHDR = ip2location(cmds.InPrecision)
		skipcol = 1
		if DEBUG == true {
			fmt.Printf("ip2location skipcol: %d\n", skipcol)
		}

		if DEBUG == true {
			fmt.Printf("%#v\n", DBHDR)
		} else {
			fmt.Printf("Using %s with %s format.\n", cmds.Format, DBHDR)
		}

	default:
		fmt.Println("Unknown format, please see --help")
		os.Exit(1)
	}

	if len(cmds.RedisHost) == 0 || cmds.RedisPort == 0 {
		fmt.Println("Please specify a redis host and port.  Even if it's the default:\nExample: -r 127.0.0.1 -p 6379")
		os.Exit(1)
	}

	if cmds.DontSkipHeader == true {
		fmt.Println("WARNING: Not skipping header!")
	} else {
		skiphdr = true
	}

	csvFile, err := os.OpenFile(cmds.CsvFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	csvreader := gocsv.LazyCSVReader(csvFile)
	samples, err = csvreader.ReadAll()
	if err != nil {
		panic(err)
	}

	/*
		if DEBUG == true {
			fmt.Printf("%#v", samples)
		}
	*/
	// Put the redis connection shit here

	i = 0 // i is the outer line range
	for _, sample := range samples {
		x = 0 // x is the inner column range
		for _, cell := range sample {
			if !(i == 0 && skiphdr == true) {
				if !(skipcol == x) {
					//push unto string array
					index = len(rediscmd)
					rediscmd = rediscmd[:index] + " " + cell
				} //skipcol
			} //skiphdr
			x++
		} // cell range

		// Use rediscmd and purge it
		if len(rediscmd) == 0 {
			fmt.Printf("Skipped Header\n")
		} else {
			fmt.Printf("REDIS<: %s, x: %d, i: %d\n", rediscmd, x, i)
			rediscmd = ""
		}
		i++
	}

}
