package main

import (
	"flag"
	"log"
	"net/http"
	"regexp"
)

func siteChecker(url string) (bool, string) {
	r, err := regexp.Compile(`^200 OK$`)
	if err != nil {
		log.Fatal("FATAL: Unable to compile the regexp")
	}
	client := &http.Client{}
	rsp, err := client.Get(url)
	if err != nil {
		log.Printf("ERR: %v\n", err)
	}

	return r.MatchString(rsp.Status), rsp.Status
}

func main() {
	siteUrl := flag.String("url", "http://www.example.com", "Url name")
	flag.Parse()

	ok, status := siteChecker(*siteUrl)

	if ok {
		log.Printf("INFO: %s is up. Status: %s\n", *siteUrl, status)
	} else {
		log.Printf("INFO: %s is down. Status: %s\n", *siteUrl, status)
	}
}
