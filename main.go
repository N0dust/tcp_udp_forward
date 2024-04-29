package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	var params arrayFlags
	flag.Var(&params, "p", "forward params")
	flag.Parse()

	for _, param := range params {
		fmt.Println(param)
		if strings.HasSuffix(param, "/udp") {
			param = strings.TrimSuffix(param, "/udp")
			from, to, err := ParseFromTo(param)
			if err != nil {
				log.Fatalln("UDP", param, err)
			}
			go UdpForward(from, to)
		}
		if strings.HasSuffix(param, "/tcp") {
			param = strings.TrimSuffix(param, "/tcp")
			from, to, err := ParseFromTo(param)
			if err != nil {
				log.Fatalln("TCP", param, err)
			}
			go TcpForward(from, to)
		}
	}
	select {}
}

func TcpForward(from string, to string) {
	tcpListen, err := net.Listen("tcp", from)
	if err != nil {
		fmt.Println(err, err.Error())
		os.Exit(0)
	}
	dIpaddr, _ := net.ResolveTCPAddr("tcp4", to)

	for {
		sConn, err := tcpListen.Accept()
		if err != nil {
			continue
		}

		dConn, err := net.DialTCP("tcp", nil, dIpaddr)
		if err != nil {
			fmt.Println(err)
			sConn.Write([]byte("connect fail" + to))
			sConn.Close()
			continue
		}
		go io.Copy(sConn, dConn)
		go io.Copy(dConn, sConn)
	}
}

func UdpForward(from string, to string) {

	srcAddr, _ := net.ResolveUDPAddr("udp", from)
	dstAddr, _ := net.ResolveUDPAddr("udp", to)

	conn, _ := net.ListenUDP("udp", srcAddr)
	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		// read data from src
		n, src, _ := conn.ReadFromUDP(buf)
		// write to dst
		_, _ = conn.WriteToUDP(buf[:n], dstAddr)

		// most demo forget about this part
		// read reply from dst
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Failed to read reply from remote UDP:", err)
			continue
		}
		// write reply to src
		_, err = conn.WriteToUDP(buf[:n], src)
		if err != nil {
			log.Println("reply write UDP", err)
		}
	}
}

func ParseFromTo(s string) (string, string, error) {
	fromTo := strings.Split(s, ":")
	if len(fromTo) != 3 {
		return "", "", errors.New("invalid from to")
	}
	from := "0.0.0.0:" + fromTo[0]
	to := fromTo[1] + ":" + fromTo[2]
	if !ValidIpPort(fromTo[0]) {
		return "", "", errors.New("invalid port")
	}
	if !ValidIPAddress(fromTo[1]) {
		return "", "", errors.New("invalid ip")
	}
	if !ValidIpPort(fromTo[2]) {
		return "", "", errors.New("invalid port")
	}

	return from, to, nil
}

func ValidIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	}
	return true
}

func ValidIpPort(port string) bool {
	p, err := strconv.Atoi(port)
	if err != nil {
		log.Println(err)
		return false
	}
	if p < 0 || p > 65535 {
		return false
	}

	return true
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprint(*i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
