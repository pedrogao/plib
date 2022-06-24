package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	"github.com/songgao/water"

	"github.com/pedrogao/plib/pkg/net/packet"
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

	iface, err := water.New(water.Config{
		DeviceType: water.TAP,
	})

	if nil != err {
		log.Fatalln("Unable to allocate TAP interface :", err)
	}
	log.Println("Interface allocated:", iface.Name())
	if *localIP == "" {
		log.Fatalln("ip is required")
	}
	// set interface parameters
	runIP("link", "set", "dev", iface.Name(), "mtu", mtu)
	runIP("addr", "add", *localIP, "dev", iface.Name())
	runIP("link", "set", "dev", iface.Name(), "up")

	// and one more loop
	rawData := make([]byte, bufferSize)
	for {
		plen, err := iface.Read(rawData)
		if err != nil {
			log.Printf("read tap err: %s", err)
			break
		}
		reader := bytes.NewReader(rawData[:plen])
		frame := packet.NewEthernetFrame()
		// frame := net.NewIcmpPacket()
		err = frame.Read(kaitai.NewStream(reader), nil, nil)
		if err != nil {
			log.Printf("parse pakcet err: %s", err)
			break
		}
		log.Println(frame.SrcMac)
		ip4 := frame.Body.(*packet.Ipv4Packet)
		icmp := ip4.Body.Body.(*packet.IcmpPacket)
		log.Println(icmp.Echo.SeqNum)
	}
}
