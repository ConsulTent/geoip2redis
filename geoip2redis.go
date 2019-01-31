package main

import (
	"fmt"
	//	"github.com/gocarina/gocsv"
	"github.com/rburmorrison/go-argue"
	//	"os"
	"providers/ip2location_db1"
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

func main() {

	var cmds cmdline

	fmt.Printf("GeoIP2Redis (c) 2019 ConsulTent Ltd. v%s-%s\n", pver, gitver)

	agmt := argue.NewEmptyArgumentFromStruct(&cmds)

	agmt.Dispute(true)

	fmt.Println("cmds.CsvFile:", cmds.CsvFile)

	switch cmds.Format {
	case "maxmind":
		fmt.Println("maxmind is not supported yet")
	case "ip2location":
		fmt.Println("ip2location loading")
		//		import "ip2location"
	default:
		fmt.Println("Unknown format, please see --help")
	}
}
