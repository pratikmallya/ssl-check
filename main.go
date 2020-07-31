package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/wpalmer/gozone"
)

func main() {
	app := cli.App{
		Name:  "ssl-check",
		Usage: "get expiration dates for your domains from a zonefile",
		Action: func(c *cli.Context) error {
			getExpirationFromZoneFile(c.Args().Get(0))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getExpirationFromZoneFile(filename string) {
	stream, _ := os.Open(filename)
	var record gozone.Record
	scanner := gozone.NewScanner(stream)

	for {
		err := scanner.Next(&record)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
			break
		}

		cfg := tls.Config{}
		conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", record.DomainName), &cfg)
		if err != nil {
			fmt.Printf("%s. Error: %s \n", record.DomainName, err)
			continue
		}
		certChain := conn.ConnectionState().PeerCertificates
		cert := certChain[len(certChain)-1]
		fmt.Printf("%s: %s\n", record.DomainName, cert.NotAfter)
	}
}
