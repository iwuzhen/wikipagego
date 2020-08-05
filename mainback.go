package main

import (
	"os"
	// "database/sql"
	// "io/ioutil"
	"io"
	// "time"
	"compress/bzip2"
	"encoding/xml"

	log "github.com/sirupsen/logrus"
	// "fmt"
)

type contributor struct {
	Username string `xml:"username"`
	Id       int    `xml:"id"`
}

type redirect struct {
	Title string `xml:"title,attr"`
}

type revision struct {
	Id          int         `xml:"id"`
	Parentid    int         `xml:"parentid"`
	Timestamp   string      `xml:"timestamp"`
	Contributor contributor `xml:"revision"`
	Comment     string      `xml:"comment"`
	Model       string      `xml:"model"`
	Format      string      `xml:"format"`
	Text        string      `xml:"text"`
	Sha1        string      `xml:"sha1"`
}

type Page struct {
	Title    string   `xml:"title"`
	Ns       int      `xml:"ns"`
	Id       int      `xml:"id"`
	Redirect redirect `xml:"redirect"`
	Revision revision `xml:"revision"`
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//   log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func test() {
	filePath := "/data/ssdj/download/enwiki-20200501-pages-articles.xml.bz2"
	// filePath :="/data/ssdj/download/t.txt.gz"

	fi, err := os.Open(filePath)
	if err != nil {
		log.Fatal("ERROR:", err)
	}

	defer fi.Close()

	fz := bzip2.NewReader(fi)
	// defer fz.Close()

	d := xml.NewDecoder(fz)
	count := 0
	for {
		count += 1
		t, tokenErr := d.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			}
			// handle error
		}
		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Local == "page" {
				var page Page
				if err := d.DecodeElement(&page, &t); err != nil {
					log.Warning(err)
				}
				// log.Info("start ",t)
				// // log.Info("page ",page)
				log.Info("title ", page.Title)
				// log.Info("Redirect ",page.Redirect)
				// time.Sleep(2*time.Second)

			}
			// else if  t.Name.Local == "title" {
			// 	innerText,_ := d.Token()
			// 	log.Info("title ",string(innerText.(xml.CharData)))
			// } else if  t.Name.Local == "redirect" {
			// 	log.Info("redirect ",t)
			// }

		}
	}
	log.Info("over ", count)
}
