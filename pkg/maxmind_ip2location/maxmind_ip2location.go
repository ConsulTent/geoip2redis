package main

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/rburmorrison/go-argue"
)

// ip2location header format:
//  start_ip, end_ip, iso_country, country_name, State, City
// "16777216","16777471","US","United States of America","California","Los Angeles"
// "34877440","34877695","US","United States of America","Texas","Dallas"

// Maxmind GeoIP2-City-Blocks-IPv4.csv format:
// network,geoname_id,registered_country_geoname_id,represented_country_geoname_id,is_anonymous_proxy,is_satellite_provider,postal_code,latitude,longitude,accuracy_radius
//1.0.67.0/25,1862415,1861060,,0,0,730-0000,34.4000,132.4500,20
// 1.0.67.128/25,1863018,1861060,,0,0,738-0031,34.2833,132.2667,20
// We only need: 1,2
// network(converted),

// MAxmind GeoIP2-City-Locations-en.csv format:
// geoname_id,locale_code,continent_code,continent_name,country_iso_code,country_name,subdivision_1_iso_code,subdivision_1_name,subdivision_2_iso_code,subdivision_2_name,city_name,metro_code,time_zone,is_in_european_union
// 4178972,en,NA,"North America",US,"United States",FL,Florida,,,"Zolfo Springs",539,America/New_York,0
// 4178992,en,NA,"North America",US,"United States",GA,Georgia,,,Abbeville,503,America/New_York,0
// We only need:
// country_iso_code,country_name,subdivision_1_name,city_name
// Skip: 1,2,3,4,7,9,10,12,13,14
// Keep: 5,6,8,11

func ip2long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

type cmdline struct {
	MaxMindCSV  string `options:"required,positional" help:"Maxmind 'blocks' CSV file."`
	LocationCSV string `options:"required,positional" help:"Maxmind 'location' CSV file."`
	Ip2Location string `options:"required,positional" help:"Output CSV file in ip2location format."`
	UseTimezone bool   `init:"t" help:"Fallback to Timezone city, when there's no data."`
}

const pver = "0.1"

var gitver = "undefined"

var DEBUG = false

func main() {
	var i, x int
	var cmds cmdline
	//	var long uint32
	var locs map[int]string
	var outline string
	var isocityreserve bool
	var locations [][]string
	var blocks [][]string
	var index int
	var geonameid int
	var network_cidr string
	//var ip2record []string
	var ipStart string
	var ipEnd string
	var Error error
	//  var iprange int

	agmt := argue.NewEmptyArgumentFromStruct(&cmds)

	agmt.Dispute(true)

	//	long = ip2long(cmds.Rawip)

	fmt.Printf("maxmind->ip2location (c) 2021 ConsulTent Pte. Ltd. v%s-%s\n", pver, gitver)

	DataFile, err := os.Open(cmds.LocationCSV)
	if err != nil {
		panic(err)
	}
	defer DataFile.Close()

	csvreader := csv.NewReader(DataFile)
	//Configure reader options Ref http://golang.org/src/pkg/encoding/csv/reader.go?s=#L81
	csvreader.Comma = ','         //field delimiter
	csvreader.Comment = '#'       //Comment character
	csvreader.FieldsPerRecord = 0 //Number of records per record. Set to Negative value for variable
	csvreader.TrimLeadingSpace = true
	csvreader.ReuseRecord = true
	csvreader.LazyQuotes = false

	locations, err = csvreader.ReadAll()
	if err != nil {
		panic(err)
	}

	if DEBUG == true {
		fmt.Printf("Columns: %d\n", csvreader.FieldsPerRecord)
	}

	fmt.Printf("Loading locations into hash table.\n")

	locs = make(map[int]string)
	bar := pb.StartNew(len(locations))

	isocityreserve = false

	i = 1 // i is the outer line range
	for _, location := range locations {
		x = 1 // x is the inner column range
		for _, cell := range location {
			if i != 1 {
				// Skip: 1,2,3,4,7,9,10,12,13,14  Keep: 5,6,8,11
				if x == 1 {
					geonameid, _ = strconv.Atoi(cell)
				}
				if x == 5 || x == 6 || x == 8 || x == 11 {
					//push unto string array
					index = len(outline)
					if index <= 1 {
						outline = "\"" + cell + "\""
						//	outline = cell
					} else {
						if x == 11 && len(cell) == 0 && cmds.UseTimezone == true {
							isocityreserve = true
						} else {
							outline = outline[:index] + ",\"" + cell + "\""
							//	 outline = outline[:index] + "," + cell
						}
					}
				} else {
					if x == 13 && isocityreserve == true && cmds.UseTimezone == true {
						isocity := strings.Split(cell, "/")
						outline = outline[:index] + ",\"" + isocity[1] + "\""
					}
				}
			} //SkipHeader
			x++
		} // cell range
		locs[geonameid] = outline
		bar.Increment()
		outline = ""
		isocityreserve = false
		i++
	}

	bar.Finish()
	/*			 if DEBUG == true {
				   fmt.Printf("Totals(locs): %d\n",len(locs))
					   for key, value := range locs {
	             fmt.Println("Key:", key, "Value:", value)
	          }
			   }
	*/

	// End of locations, now we do the data.

	DataFile, err = os.Open(cmds.MaxMindCSV)
	if err != nil {
		panic(err)
	}
	defer DataFile.Close()

	csvreader = csv.NewReader(DataFile)
	//Configure reader options Ref http://golang.org/src/pkg/encoding/csv/reader.go?s=#L81
	csvreader.Comma = ','         //field delimiter
	csvreader.Comment = '#'       //Comment character
	csvreader.FieldsPerRecord = 0 //Number of records per record. Set to Negative value for variable
	csvreader.TrimLeadingSpace = true
	csvreader.ReuseRecord = true
	csvreader.LazyQuotes = false

	blocks, err = csvreader.ReadAll()
	if err != nil {
		panic(err)
	}

	if DEBUG == true {
		fmt.Printf("Columns: %d\n", csvreader.FieldsPerRecord)
	}

	fmt.Printf("Cross-referencing and creating data.\n")

	csvOut, err := os.OpenFile(cmds.Ip2Location, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}

	bar = pb.StartNew(len(blocks))

	i = 1 // i is the outer line range
	for _, block := range blocks {
		x = 1 // x is the inner column range
		for _, cell := range block {
			/*							 if DEBUG == true {
										 fmt.Printf("i: %d, x: %d, cell: %s\n",i,x, cell)
									 } */
			if x > 10 {
				fmt.Printf("Warning: there are %d extra columns.", x-10)
			}
			if i != 1 {
				// Keep: 1,2
				if x == 1 {
					network_cidr = cell
				}
				if x == 2 {
					geonameid, _ = strconv.Atoi(cell)
				}
				if x == 3 {
					if geonameid == 0 {
						geonameid, _ = strconv.Atoi(cell)
						if DEBUG == true {
							fmt.Println("geonameid was 0")
						}
					}
				}
			} //SkipHeader
			x++
		} // cell range

		if i != 1 {
			if DEBUG == true {
				fmt.Printf("CIDR: %s\n", network_cidr)
			}
			if len(network_cidr) == 0 {
				fmt.Println("Received blank value for CIDR in Maxmind source, exiting.")
				os.Exit(1)
			}
			ipStart, ipEnd, Error = CIDRRangeToIPv4Range(network_cidr)
			if Error != nil {
				panic(Error)
			}
			outline = locs[geonameid]
			/*									if len(outline) <= 1 { fmt.Println("Location data was blank, exiting.")
												os.Exit(1)
												} */
			//								  ip2record = fmt.Sprintf("\"%s\",\"%s\",\"%s\"",iPv4ToUint32(ipStart), iPv4ToUint32(ipEnd), outline)
			//										ip2record = append(ip2record[:0],string(iPv4ToUint32(ipStart)),string(iPv4ToUint32(ipEnd)))
			//									ip2record = "\" + iPv4ToUint32(ipStart) + \" + \" + iPv4ToUint32(ipEnd) + \""
			if DEBUG == true {
				fmt.Printf("line: %d, outline: %s\n", i, outline)
			}
			//ip2record = strings.Split(fmt.Sprintf("\"%s\",\"%s\",\"%s\"",string(iPv4ToUint32(ipStart)), string(iPv4ToUint32(ipEnd)), outline), ",")
			//ip2record = strings.Split(fmt.Sprintf("\"%d\",\"%d\",%s",iPv4ToUint32(ipStart), iPv4ToUint32(ipEnd), outline), ",")
			if iPv4ToUint32(ipStart) == 16777216 && i == 2 {
				_, err := io.WriteString(csvOut, fmt.Sprintf("\"0\",\"16777215\",\"-\",\"-\",\"-\",\"-\"\n"))
				if err != nil {
					panic(err)
				}
				if DEBUG == true {
					fmt.Println("iPv4ToUint32(ipStart) == 16777216 && i == 2")
				}
			}

			if len(outline) > 1 {
				_, err := io.WriteString(csvOut, fmt.Sprintf("\"%d\",\"%d\",%s\n", iPv4ToUint32(ipStart), iPv4ToUint32(ipEnd), outline))
				if err != nil {
					panic(err)
				}
				if DEBUG == true {
					fmt.Printf("outlinesize: %d\n", len(strings.Split(outline, ",")))
				}
			} else {
				if DEBUG == true {
					fmt.Printf("outlinesize: %d - SKIPPED\n", len(strings.Split(outline, ",")))
				}
			}

		}
		bar.Increment()
		outline = ""
		i++
	}

	bar.Finish()

	csvOut.Sync()
	csvOut.Close()

}

// Thank you Valentyn Ponomarenko! https://gist.github.com/P-A-R-U-S/3c54bacef489499e2b44a075fdab6af0

// Convert CIDR to IPv4 range
func CIDRRangeToIPv4Range(cidr string) (ipStart string, ipEnd string, err error) {

	var ip uint32 // ip address

	var ipS uint32 // Start IP address range
	var ipE uint32 // End IP address range

	cidrParts := strings.Split(cidr, "/")

	ip = iPv4ToUint32(cidrParts[0])
	bits, _ := strconv.ParseUint(cidrParts[1], 10, 32)

	if ipS == 0 || ipS > ip {
		ipS = ip
	}

	ip = ip | (0xFFFFFFFF >> bits)

	if ipE < ip {
		ipE = ip
	}

	ipStart = uInt32ToIPv4(ipS)
	ipEnd = uInt32ToIPv4(ipE)

	return ipStart, ipEnd, err
}

//Convert IPv4 to uint32
func iPv4ToUint32(iPv4 string) uint32 {

	ipOctets := [4]uint64{}

	for i, v := range strings.SplitN(iPv4, ".", 4) {
		ipOctets[i], _ = strconv.ParseUint(v, 10, 32)
	}

	result := (ipOctets[0] << 24) | (ipOctets[1] << 16) | (ipOctets[2] << 8) | ipOctets[3]

	return uint32(result)
}

//Convert uint32 to IP
func uInt32ToIPv4(iPuInt32 uint32) (iP string) {
	iP = fmt.Sprintf("%d.%d.%d.%d",
		iPuInt32>>24,
		(iPuInt32&0x00FFFFFF)>>16,
		(iPuInt32&0x0000FFFF)>>8,
		iPuInt32&0x000000FF)
	return iP
}
