package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cheggaaa/pb"
	m2i "github.com/consultent/geoip2redis/pkg/maxmind_ip2location"
	spinner "github.com/consultent/go-spinner"
	"github.com/go-redis/redis"
	argue "github.com/rburmorrison/go-argue"
)

type cmdline struct {
	CsvFile         string `options:"required,positional" help:"CSV file with GeoIP data"`
	Format          string `init:"f" options:"required" help:"ip2location|maxmind"`
	RedisHost       string `init:"r" help:"Redis Host, default 127.0.0.1"`
	RedisPort       int    `init:"p" help:"Redis Port, default 6379"`
	RedisPass       string `init:"a" help:"Redis DB password, default none"`
	InPrecision     int    `init:"i" options:"required" help:"Input precision. Optional. This would be db file number. 1=DB1 for ip2location. Default is autodetect.  See README.TXT"`
	ForceDbhdr      string `init:"d" help:"Force a custom subkey where the GeoIP data will be stored, instead of using defaults."`
	ForceAutodetect bool   `init:"t" help:"Force autodetect of database type, NOT format.  Optional.  This will ignore input precision, and set a default header"`
	SkipHeader      bool   `init:"s" help:"Force skip the first CSV line. Default: follows format, see README.TXT"`
	TempDir         string `init:"c" help:"Set temporary work directory for conversions. Default: ./"`
	MaxmindLocation string `init:"m" help:"Maxmind 'location' CSV file. Only specify when '--format maxmind'"`
	UseTimezone     bool   `init:"z" help:"Fallback to Timezone city, when there's no data. (For MaxMind)"`
}

// MaxMindCSV is now CsvFile, Ip2Location string is now a temp file (which will be the real CsvFile)
const pver = "1.0.0"

var gitver = "undefined"

var DEBUG = false

var TempFile string

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

	//	var CSVFile string

	//	var zcmd int64
	//      var fakedata struct{}

	fmt.Printf("GeoIP2Redis (c) 2021 ConsulTent Pte. Ltd. v%s build:%s\n", pver, gitver)

	agmt := argue.NewEmptyArgumentFromStruct(&cmds)

	agmt.Dispute(true)

	if len(cmds.TempDir) != 0 {
		_, err := os.Stat(cmds.TempDir)
		if os.IsNotExist(err) {
			cleanupTemp()
			log.Fatal(fmt.Sprintf("Temp directory %s does not exist.", cmds.TempDir))
		}
	}

	switch cmds.Format {
	case "maxmind":
		if len(cmds.MaxmindLocation) == 0 {
			fmt.Println("Please specify MaxMind location csv data file.")
			os.Exit(1)
		}
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
		go func() {
			<-sigs
			cleanupTemp()
			os.Exit(0)
		}()

		TempFile = m2i.MaxMind_Merge(cmds.CsvFile, cmds.MaxmindLocation, cmds.UseTimezone, cmds.TempDir)

		fallthrough
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
		fmt.Println("ASN not supported yet.")
		os.Exit(2)
		if cmds.ForceAutodetect == false {
			fmt.Println("ASN format only available with autodetect")
			os.Exit(1)
		}
		Ip2info = ip2location(cmds.InPrecision, cmds.ForceAutodetect)
		CSVinfo = Ip2info.GenericCsvFormat
		CSVinfo.Header = "ASN"
	case "software77":
		fmt.Println("WARNING: Software77 support is legacy and deprecated.")
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

	s := spinner.StartNew("Reading in CSV data.")
	s.SetCharset([]string{"\U0001F311", "\U0001F312", "\U0001F313", "\U0001F314", "\U0001F315", "\U0001F316", "\U0001F317", "\U0001F318"})
	csvFile, err := os.Open(TempFile)
	if err != nil {
		cleanupTemp()
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
		cleanupTemp()
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
			fmt.Printf("\nWarning: Detected too many fields (%d), max is (%d), proceed with caution!\n", csvreader.FieldsPerRecord, CSVinfo.MaxDB())
		}
		fmt.Printf("\nUsing set %s for autodetection\n", DBHDR)
	}

	if DEBUG == true {
		fmt.Printf("DBHDR: %s, %d, %s\n", DBHDR, CSVinfo.Formatin, CSVinfo.DbOutHdr())
	} else {
		fmt.Printf("\nLoading into set %s\n", DBHDR)
	}

	bcounter = len(samples)

	s.Stop()

	// Put the redis connection stuff here
	s = spinner.StartNew(fmt.Sprintf("Connecting to Redis server: %s:%d", cmds.RedisHost, cmds.RedisPort))

	redisClient := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%v:%v", cmds.RedisHost, cmds.RedisPort),
		Password: cmds.RedisPass, DB: 0, ReadTimeout: time.Minute, WriteTimeout: time.Minute})
	_, err = redisClient.Ping().Result()
	if err != nil {
		cleanupTemp()
		log.Printf("[ERROR] unable to ping redis client, error : %v", err)
		os.Exit(1)
	}

	defer redisClient.Close()

	s.Stop()

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
				cleanupTemp()
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
		cleanupTemp()
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

	cleanupTemp()

	fmt.Printf("Loaded %d entries into Redis\n", i)
}

func cleanupTemp() {
	if len(TempFile) != 0 {
		os.Remove(TempFile)
	}
}
