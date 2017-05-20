package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ericyan/shield/transport/udp"
	"github.com/ericyan/shield/tuntap"
)

func main() {
	bind := flag.String("bind", "0.0.0.0:1984", "ip and port to bind")
	peer := flag.String("peer", "", "ip and port of the peer")
	ip := flag.String("ip", "172.16.0.1/31", "ip of the tun device (in CIDR)")
	mtu := flag.Int("mtu", 1350, "MTU of the tun device")
	bufsize := flag.Int("bufsize", 1500, "size of the buffer")
	flag.Parse()

	// Set up the TUN device
	tun, err := tuntap.NewTUN()
	if err != nil {
		log.Fatal(err)
	}
	tun.SetMTU(*mtu)
	addr, err := tuntap.ParseCIDR(*ip)
	if err != nil {
		log.Fatal(err)
	}
	tun.AddAddr(addr)
	tun.Up()

	// Set up the UDP transport
	udp, err := udp.New(*bind, *peer)
	if err != nil {
		log.Fatal(err)
	}

	// Anything received from UDP goes to TUN
	go func() {
		buf := make([]byte, *bufsize)
		_, err := io.CopyBuffer(tun, udp, buf)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Anything received from TUN goes to UDP
	go func() {
		buf := make([]byte, *bufsize)
		_, err := io.CopyBuffer(udp, tun, buf)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("Listening on %s with peer %s", udp.LocalAddr(), udp.RemoteAddr())
	addrs, err := tun.Addrs()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s: %v", tun.Attrs().Name, addrs)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Fatalf("Signal (%v) received, stopping\n", s)
}
