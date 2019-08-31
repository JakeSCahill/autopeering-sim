package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"

	"github.com/pkg/errors"
	"github.com/wollac/autopeering/discover"
	"github.com/wollac/autopeering/logger"
	"github.com/wollac/autopeering/peer"
	"github.com/wollac/autopeering/transport"
)

const defaultZLC = `{
	"level": "info",
	"development": false,
	"outputPaths": ["stdout"],
	"errorOutputPaths": ["stderr"],
	"encoding": "console",
	"encoderConfig": {
	  "timeKey": "ts",
	  "levelKey": "level",
	  "nameKey": "logger",
	  "callerKey": "caller",
	  "messageKey": "msg",
	  "stacktraceKey": "stacktrace",
	  "lineEnding": "",
	  "levelEncoder": "",
	  "timeEncoder": "iso8601",
	  "durationEncoder": "",
	  "callerEncoder": ""
	}
  }`

func waitInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func parseMaster(s string) (*peer.Peer, error) {
	if len(s) == 0 {
		return nil, nil
	}

	parts := strings.Split(s, "@")
	if len(parts) != 2 {
		return nil, errors.New("parseMaster")
	}
	pubKey, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.Wrap(err, "parseMaster")
	}

	return peer.NewPeer(pubKey, parts[1]), nil
}

func main() {
	var (
		listenAddr = flag.String("addr", "127.0.0.1:14626", "listen address")
		masterNode = flag.String("master", "", "master node as 'pubKey@address' where pubKey is in Base64")

		err error
	)
	flag.Parse()

	logger := logger.NewLogger(defaultZLC, "debug")
	defer logger.Sync()

	priv, err := peer.GeneratePrivateKey()
	if err != nil {
		log.Fatalf("GeneratePrivateKey: %v", err)
	}

	addr, err := net.ResolveUDPAddr("udp", *listenAddr)
	if err != nil {
		log.Fatalf("ResolveUDPAddr: %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("ListenUDP: %v", err)
	}
	defer conn.Close()

	cfg := discover.Config{
		Log: logger.Named("discover"),
	}

	master, err := parseMaster(*masterNode)
	if err != nil {
		log.Printf("Ignoring master: %v\n", err)
	} else if master != nil {
		cfg.Bootnodes = []*peer.Peer{master}
	}

	// use the UDP connection for transport
	trans := transport.Conn(conn, func(network, address string) (net.Addr, error) { return net.ResolveUDPAddr(network, address) })
	defer trans.Close()

	// start the discovery on that connection
	disc := discover.Listen(trans, peer.NewLocal(priv, peer.NewMapDB(logger.Named("db"))), cfg)
	defer disc.Close()

	id := base64.StdEncoding.EncodeToString(disc.Local().PublicKey())
	fmt.Println("Discovery protocol started: ID=" + id + ", address=" + disc.LocalAddr())
	fmt.Println("Hit Ctrl+C to exit")

	// wait for Ctrl+c
	waitInterrupt()
}
