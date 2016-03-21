package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/dockermec/etcdlib"
	"github.com/dockermec/g"
	"github.com/dockermec/http"
	"net"
	"os"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v1", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	//write the ip of the host
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops:" + err.Error())
		os.Exit(1)
	}
	ipstr := ""
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				os.Stdout.WriteString(ipnet.IP.String() + "\n")
				ipstr = ipstr + ipnet.IP.String()
			}
		}
	}
	fmt.Println(ipstr)
	//md5
	h := md5.New()
	h.Write([]byte(ipstr))
	fmt.Printf("%s\n", hex.EncodeToString(h.Sum(nil)))
	//etcd write
	endpoints := []string{"http://Master:2379"}
	startwork := etcdlib.NewWorker(hex.EncodeToString(h.Sum(nil)), hex.EncodeToString(h.Sum(nil)), endpoints)
	fmt.Println(startwork)
	//load the config file
	g.ParseConfig(*cfg)

	go http.Start()

	select {}
}
