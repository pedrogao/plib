package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/songgao/water"
	"golang.org/x/net/ipv4"
)

const (
	bufferSize = 1500
	mtu        = "1300"
)

var (
	localIP = flag.String("local", "", "Local tun interface IP/MASK like 192.0.2.1/24")
)

func runIP(args ...string) {
	cmd := exec.Command("/sbin/ip", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if nil != err {
		log.Fatalln("Error running /sbin/ip:", err)
	}
}

func main() {
	flag.Parse()
	// check if we have anything
	if "" == *localIP {
		flag.Usage()
		log.Fatalln("\nlocal ip is not specified")
	}

	// create TUN interface
	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if nil != err {
		log.Fatalln("Unable to allocate TUN interface :", err)
	}
	log.Println("Interface allocated:", iface.Name())
	// set interface parameters
	runIP("link", "set", "dev", iface.Name(), "mtu", mtu)
	runIP("addr", "add", *localIP, "dev", iface.Name())
	runIP("link", "set", "dev", iface.Name(), "up")

	// and one more loop
	packet := make([]byte, bufferSize)
	for {
		plen, err := iface.Read(packet)
		if err != nil {
			break
		}
		// debug :)
		header, _ := ipv4.ParseHeader(packet[:plen])
		fmt.Printf("Rece: %+v (%+v)\n", header, err)
	}

}
