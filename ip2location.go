package main

import (
	"fmt"
	"os"
)

func ip2location(format int) string {

	// ip2location testing: https://play.golang.org/p/OWgH6N_Kmil
	switch format {
	case 1:
		if DEBUG == true {
			fmt.Println("ip2location DB1")
		}
		return "DB1"
	case 11:
		if DEBUG == true {
			fmt.Println("ip2location DB11")
		}
		return "DB11"
	default:
		fmt.Println("ip2location format unknown")
		os.Exit(1)
	}

	return "DB0"
}
