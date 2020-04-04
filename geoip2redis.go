package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/go-redis/redis"
	argue "github.com/rburmorrison/go-argue"
)

type cmdline struct {
	CsvFile         string `options:"required,positional" help:"CSV file with GeoIP data"`
	Format          string `init:"f" options:"required" help:"ip2location|software77"`
	RedisHost       string `init:"r" help:"Redis Host, default 127.0.0.1"`
	RedisPort       int    `init:"p" help:"Redis Port, default 6379"`
	RedisPass       string `init:"a" help:"Redis DB password, default none"`
	InPrecision     int    `init:"i" options:"required" help:"Input precision. Optional. This would be db file number. 1=DB1 for ip2location. Default is autodetect.  See README.TXT"`
	ForceDbhdr      string `init:"d" help:"Force a custom subkey where the GeoIP data will be stored, instead of using defaults."`
	ForceAutodetect bool   `init:"t" help:"Force autodetect.  Optional.  This will ignore input precision, and set a default header"`
	SkipHeader      bool   `init:"s" help:"Foce skip the first CSV line. Default: follows format, see README.TXT"`
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

const pver = "0.9.3"

var gitver = "undefined"

var DEBUG = false

func main() {

	var i int
	var x int
	var rediscmd string
	var cmds cmdline
	var samples [][]string
	var DBHDR string
	var DBHDRtemp string
	var index int
	var iprange int
	var bcounter int
	var finalcount int

	var CSVinfo GenericCsvFormat
	var Ip2info Ip2LocationCSV
	var S77info Software77CSV

	//	var zcmd int64
	//      var fakedata struct{}

	fmt.Printf("GeoIP2Redis (c) 2020 ConsulTent Ltd. v%s-%s\n", pver, gitver)

	agmt := argue.NewEmptyArgumentFromStruct(&cmds)

	agmt.Dispute(true)

	switch cmds.Format {
	case "maxmind":
		fmt.Println("maxmind is not supported yet")
		os.Exit(1)
	case "ip2location":
		Ip2info = ip2location(cmds.InPrecision, cmds.ForceAutodetect)
		CSVinfo = Ip2info.GenericCsvFormat
		if CSVinfo.Status == false {
			fmt.Printf("Database format out of range for %s\n", cmds.Format)
			os.Exit(1)
		}
		if CSVinfo.Autodetect == true {
			fmt.Println("Autodetecting ip2location db format.")
		}
	case "ip2location/asn":
		// ASN is a hack and broken.  We need to create a new definition for it
		// and use ip2long
		fmt.Println("ASN not supported.")
		os.Exit(2)
		if cmds.ForceAutodetect == false {
			fmt.Println("ASN format only available with autodetect")
			os.Exit(1)
		}
		Ip2info = ip2location(cmds.InPrecision, cmds.ForceAutodetect)
		CSVinfo = Ip2info.GenericCsvFormat
		CSVinfo.Header = "ASN"
	case "software77":
		S77info = software77(cmds.InPrecision)
		CSVinfo = S77info.GenericCsvFormat
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

	csvFile, err := os.Open(cmds.CsvFile)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	csvreader := csv.NewReader(csvFile)
	//Configure reader options Ref http://golang.org/src/pkg/encoding/csv/reader.go?s=#L81
	csvreader.Comma = ','         //field delimiter
	csvreader.Comment = '#'       //Comment character
	csvreader.FieldsPerRecord = 0 //Number of records per record. Set to Negative value for variable
	csvreader.TrimLeadingSpace = true
	csvreader.ReuseRecord = true

	samples, err = csvreader.ReadAll()
	if err != nil {
		panic(err)
	}

	if DEBUG == true {
		fmt.Printf("Columns: %d\n", csvreader.FieldsPerRecord)
	}

	if len(cmds.ForceDbhdr) == 0 {
		DBHDR = CSVinfo.DbOutHdr()
	} else {
		DBHDR = cmds.ForceDbhdr
	}

	if CSVinfo.Autodetect == true {
		if CSVinfo.IsMaxDB(csvreader.FieldsPerRecord) == true {
			fmt.Printf("Warning: Detected too many fields (%d), max is (%d), proceed with caution!\n", csvreader.FieldsPerRecord, CSVinfo.MaxDB())
		}
		fmt.Printf("Using set %s for autodetection\n", DBHDR)
	}

	if DEBUG == true {
		fmt.Printf("DBHDR: %s, %d, %s\n", DBHDR, CSVinfo.Formatin, CSVinfo.DbOutHdr())
	} else {
		fmt.Printf("Loading into set %s\n", DBHDR)
	}

	bcounter = len(samples)

	// Put the redis connection stuff here

	if DEBUG == true {
		fmt.Printf("RedisHost: %s, RedisPort: %d\n", cmds.RedisHost, cmds.RedisPort)
	}

	redisClient := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%v:%v", cmds.RedisHost, cmds.RedisPort),
		Password: cmds.RedisPass, DB: 0, ReadTimeout: time.Minute, WriteTimeout: time.Minute})
	_, err = redisClient.Ping().Result()
	if err != nil {
		log.Printf("[ERROR] unable to ping redis client, error : %v", err)
		os.Exit(1)
	}

	defer redisClient.Close()

	rep, err := redisClient.Keys(DBHDR).Result()
	if err != nil {
		fmt.Println("No subkeys found.")
	}

	if len(rep) > 0 {
		DBHDRtemp = DBHDR
		DBHDR = DBHDR + "X"
		fmt.Println("Live migration detected.  Using temporary set: ", DBHDR)
	}

	bar := pb.StartNew(bcounter)

	i = 0 // i is the outer line range
	for _, sample := range samples {
		x = 0 // x is the inner column range
		for _, cell := range sample {
			if CSVinfo.DoSkipLine(i) == false {
				if CSVinfo.DoSkipCol(x) == false {
					//push unto string array
					index = len(rediscmd)
					if x == CSVinfo.Iprangecol {
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
			_, err = redisClient.Do(CSVinfo.RedisCMD, DBHDR, iprange, rediscmd).Result()
			if err != nil {
				panic(err)
			}

			rediscmd = ""
			bar.Increment()
		}
		i++
	}

	bar.Finish()

	// Verify loaded count with bcounter
	fcount := redisClient.ZCount(DBHDR, "0", "+inf")

	finalcount = int(fcount.Val())

	if finalcount-bcounter != 0 {
		fmt.Println("Loaded count mismatch: ", finalcount)
		fmt.Println("Please do manual clean up of ", DBHDR)
		os.Exit(2)
	}

	if len(DBHDRtemp) > 0 {
		_, err = redisClient.Do("DEL", DBHDRtemp).Result()
		if err != nil {
			fmt.Println("Error deleting ", DBHDRtemp)
			fmt.Println(err)
		}
		_, err = redisClient.Do("RENAME", DBHDR, DBHDRtemp).Result()
		if err != nil {
			fmt.Println("Error renaming ", DBHDR)
			fmt.Println(err)
		}
	}

	fmt.Printf("Loaded %d entries into Redis\n", i)
}
