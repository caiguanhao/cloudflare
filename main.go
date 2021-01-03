package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/cloudflare/cloudflare-go"
)

var (
	api *cloudflare.API
)

func init() {
	var err error
	api, err = cloudflare.NewWithAPIToken(getToken())
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "ls":
		if flag.NArg() != 1 {
			showHelp()
		}
		listZones()
	case "records":
		if flag.NArg() != 2 {
			showHelp()
		}
		listZoneRecords(flag.Arg(1))
	case "addrecord":
		if flag.NArg() != 5 {
			showHelp()
		}
		createDNSRecord(flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4))
	case "delrecord":
		if flag.NArg() != 3 {
			showHelp()
		}
		deleteDNSRecord(flag.Arg(1), flag.Arg(2))
	default:
		showHelp()
	}
}

func getToken() (token string) {
	var err error
	token = os.Getenv("CF_TOKEN")
	if token != "" {
		return
	}
	var home string
	home, err = os.UserHomeDir()
	if err != nil {
		return
	}
	file := filepath.Join(home, ".cloudflare.json")
	var fileContent []byte
	fileContent, err = ioutil.ReadFile(file)
	if err != nil {
		return
	}
	var config struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return
	}
	token = config.Token
	return
}

func showHelp() {
	fmt.Println("commands:")
	fmt.Println("  ls")
	fmt.Println("  records <domain>")
	fmt.Println("  addrecord <domain> <name> <type> <value>")
	fmt.Println("  delrecord <domain> <recordid>")
	os.Exit(1)
}

func listZones() {
	zones, err := api.ListZones()
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	err = enc.Encode(zones)
	if err != nil {
		log.Fatal(err)
	}
}

func listZoneRecords(name string) {
	zoneID, err := api.ZoneIDByName(name)
	if err != nil {
		log.Fatal(err)
	}
	recs, err := api.DNSRecords(zoneID, cloudflare.DNSRecord{})
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	err = enc.Encode(recs)
	if err != nil {
		log.Fatal(err)
	}
}

func createDNSRecord(name, dname, dtype, dvalue string) {
	zoneID, err := api.ZoneIDByName(name)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := api.CreateDNSRecord(zoneID, cloudflare.DNSRecord{
		Name:    dname,
		Type:    dtype,
		Content: dvalue,
	})
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	err = enc.Encode(resp)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteDNSRecord(name, recordId string) {
	zoneID, err := api.ZoneIDByName(name)
	if err != nil {
		log.Fatal(err)
	}
	err = api.DeleteDNSRecord(zoneID, recordId)
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	err = enc.Encode(true)
	if err != nil {
		log.Fatal(err)
	}
}
