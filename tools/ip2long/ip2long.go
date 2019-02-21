package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/rburmorrison/go-argue"
	"net"
)

func ip2long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

type cmdline struct {
	Rawip  string `options:"required,positional" help:"IPv4 address"`
	Silent bool   `init:"s" help:"Silent operation (for piping)"`
}

const pver = "1.0"

var gitver = "undefined"

func main() {
	var cmds cmdline
	var long uint32

	agmt := argue.NewEmptyArgumentFromStruct(&cmds)

	agmt.Dispute(true)

	long = ip2long(cmds.Rawip)

	if cmds.Silent == false {
		fmt.Printf("ip2long (c) 2019 ConsulTent Ltd. v%s-%s\n", pver, gitver)
		fmt.Printf("IP4v: %s -> %d\n", cmds.Rawip, long)
	} else {
		fmt.Println(long)
	}

}
