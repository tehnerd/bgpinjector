package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tehnerd/bgp2go"
)

func ReadFile(fname string, to chan bgp2go.BGPProcessMsg) {
	fd, err := os.Open(fname)
	defer fd.Close()
	if err != nil {
		fmt.Printf("file doesnt exists")
		os.Exit(-1)
	}
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "N", "N6":
			to <- bgp2go.BGPProcessMsg{Cmnd: "AddNeighbour", Data: strings.Join(fields[1:], " ")}
		case "R":
			to <- bgp2go.BGPProcessMsg{Cmnd: "AddV4Route", Data: fields[1]}
		case "R6":
			to <- bgp2go.BGPProcessMsg{Cmnd: "AddV6Route", Data: fields[1]}
		}
	}
}

func main() {
	duration := flag.Uint64("duration", 3600, "Time to sleep before exit")
	asn := flag.Uint64("ASN", 65101, "Our AS number")
	rid := flag.Uint64("RID", 1, "Our router's ID")
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Println("exiting. not enough arguments")
		os.Exit(-1)
	}
	fmt.Println("starting to inject bgp routes")
	to := make(chan bgp2go.BGPProcessMsg)
	from := make(chan bgp2go.BGPProcessMsg)
	bgpContext := bgp2go.BGPContext{ASN: uint32(*asn), RouterID: uint32(*rid)}
	go bgp2go.StartBGPProcess(to, from, bgpContext)
	ReadFile(os.Args[1], to)
	time.Sleep(time.Duration(*duration) * time.Second)
	fmt.Println("i've slept enough. waking up...")
}
