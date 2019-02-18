package main

import (
	"strconv"
)

func ip2location(format int) string {

	iprangemax := 24

	if format > iprangemax {
		return "DB0"
	}
	return "DB" + strconv.Itoa(format)
	// ip2location testing: https://play.golang.org/p/OWgH6N_Kmil
}
