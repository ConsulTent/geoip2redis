package main

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/rburmorrison/go-argue"
	"os"
)

type cmdline struct {
	CsvFile      string `options:"required,positional" help:"CSV file with GeoIP data"`
	Format       string `init:"f" options:"required" help:"ip2location|maxmind"`
	RedisHost    string `init:"r" help:"Redis Host, default 127.0.0.1"`
	RedisPort    int    `init:"p" help:"Redis Port, default 6379"`
	InPrecision  int    `init:"i" options:"required" help:"Input precision.  This would be db file number.  See README.TXT"`
	OutPrecision int    `init:"o" help:"Output Precision.  Default: 0 (match input), see README.TXT"`
}

const pver = "0.0.1"

var gitver = "undefined"

type csvstruct []struct{}

var DEBUG bool = true

func main() {

	var cmds cmdline
	var samples [][]string
	var DBHDR string
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

		if DEBUG == true {
			fmt.Printf("%#v", DBHDR)
		}

	default:
		fmt.Println("Unknown format, please see --help")
		os.Exit(1)
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

	if DEBUG == true {
		fmt.Printf("%#v", samples)
	}

}
