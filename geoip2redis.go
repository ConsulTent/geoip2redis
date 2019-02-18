package main

import (
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/gocarina/gocsv"
	"github.com/keimoon/gore"
	"github.com/rburmorrison/go-argue"
	"os"
	"strconv"
)

type cmdline struct {
	CsvFile      string `options:"required,positional" help:"CSV file with GeoIP data"`
	Format       string `init:"f" options:"required" help:"ip2location|maxmind"`
	RedisHost    string `init:"r" help:"Redis Host, default 127.0.0.1"`
	RedisPort    int    `init:"p" help:"Redis Port, default 6379"`
	RedisPass    string `init:"a" help:"Redis DB password, default none (unused)"`
	InPrecision  int    `init:"i" options:"required" help:"Input precision.  This would be db file number. 1=DB1 for ip2location.  See README.TXT"`
	OutPrecision int    `init:"o" help:"Output Precision.  Default: 0 (match input), see README.TXT (unused)"`
	SkipHeader   bool   `init:"s" help:"Skip the first CSV line. Default: don't skip, see README.TXT"`
}

const pver = "0.0.2"

var gitver = "undefined"

type csvstruct []struct{}

var DEBUG = false

func main() {

	var i uint64
	var x uint8
	var skipcol uint8
	var rediscmd string
	var cmds cmdline
	var samples [][]string
	var DBHDR string
	var index int
	var iprange int
	var bcounter int
	//	var zcmd int64
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
		if DBHDR == "DB0" {
			fmt.Println("WARNING: Format out of range for ip2locartion.\nProceeding with set DB0")
		}
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

	if cmds.SkipHeader == true {
		fmt.Println("WARNING: Skipping header!")
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

	bcounter = len(samples)

	/*
		if DEBUG == true {
			fmt.Printf("%#v", samples)
		}
	*/
	// Put the redis connection stuff here

	if DEBUG == true {
		fmt.Printf("RedisHost: %s, RedisPort: %d\n", cmds.RedisHost, cmds.RedisPort)
	}

	redisdb, err := gore.Dial(cmds.RedisHost + ":" + strconv.Itoa(cmds.RedisPort))
	if err != nil {
		panic(err)
	}

	defer redisdb.Close()

	bar := pb.StartNew(bcounter)

	i = 0 // i is the outer line range
	for _, sample := range samples {
		x = 0 // x is the inner column range
		for _, cell := range sample {
			if !(i == 0 && cmds.SkipHeader == true) {
				if !(skipcol == x) {
					//push unto string array
					index = len(rediscmd)
					if x == 0 {
						iprange, _ = strconv.Atoi(cell)
						//		rediscmd = rediscmd[:index] + cell + " \"" + cell
						rediscmd = rediscmd[:index] + cell
					} else {
						rediscmd = rediscmd[:index] + "|" + cell
					}
				} //skipcol
			} //SkipHeader
			x++
		} // cell range

		// Use rediscmd and purge it
		if len(rediscmd) == 0 {
			fmt.Printf("Skipped Header\n")
		} else {
			//index = len(rediscmd)
			//		rediscmd = rediscmd[:index] + "\""
			// if DEBUG == true { fmt.Printf("REDIS<: %s\nx: %d, i: %d\n", rediscmd, x, i) }
			_, err = gore.NewCommand("ZADD", DBHDR, iprange, rediscmd).Run(redisdb)
			if err != nil {
				panic(err)
			}
			rediscmd = ""
			bar.Increment()
		}
		i++
	}

	bar.Finish()
	fmt.Printf("Loaded %d entries into Redis\n", i)
}
