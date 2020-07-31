package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/table"
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
	tw := table.NewWriter()
	tw.SetAutoIndex(true)
	tw.SetStyle(table.StyleLight)

	tw.AppendHeader(table.Row{"Domain", "Expiration"})

	for {
		err := scanner.Next(&record)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
			break
		}

		if !(record.Type == gozone.RecordType_A ||
			record.Type == gozone.RecordType_CNAME) {
			continue
		}

		cfg := tls.Config{}
		conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", record.DomainName), &cfg)
		if err != nil {
			tw.AppendRow(table.Row{record.DomainName, err})
			continue
		}
		certChain := conn.ConnectionState().PeerCertificates
		cert := certChain[len(certChain)-1]
		tw.AppendRow(table.Row{record.DomainName, cert.NotAfter})
	}
	fmt.Println(tw.Render())
}
