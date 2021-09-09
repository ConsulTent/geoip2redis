package main

import (
	"fmt"
	"os"
	"time"

	m2i "github.com/consultent/geoip2redis/pkg/maxmind_ip2location"
	spinner "github.com/consultent/go-spinner"
	argue "github.com/rburmorrison/go-argue"
)

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
	var cmds cmdline
	var TempFile string

	fmt.Printf("MaxMind->Ip2Location CSV converter. (c) 2021 ConsulTent Pte. Ltd. v%s-%s\n", pver, gitver)

	agmt := argue.NewEmptyArgumentFromStruct(&cmds)

	agmt.Dispute(true)

	_, err := os.Stat(cmds.MaxMindCSV)
	if err != nil {
		fmt.Printf("%s doesn't exit.", cmds.MaxMindCSV)
		panic(err)
	}

	_, err = os.Stat(cmds.LocationCSV)
	if err != nil {
		fmt.Printf("%s doesn't exit.", cmds.LocationCSV)
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("WARNING: Couldn't get current working directory (%s).  Using workaround for error %s.\n.", cwd, err)
		cwd = "."
	}

	TempFile = m2i.MaxMind_Merge(cmds.MaxMindCSV, cmds.LocationCSV, cmds.UseTimezone, cwd)

	s := spinner.StartNew("Moving CSV from Temporary Directory.")
	s.SetSpeed(100 * time.Millisecond)
	s.SetCharset([]string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"})
	err = os.Rename(TempFile, cmds.Ip2Location)
	if err != nil {
		fmt.Printf("Error moving %s -> %s\n", TempFile, cmds.Ip2Location)
		os.Remove(TempFile)
		panic(err)
	}
	s.Stop()

	fmt.Printf("Created: %s\n", cmds.Ip2Location)

}
