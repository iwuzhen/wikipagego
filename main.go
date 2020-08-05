package main

import (
	"compress/bzip2"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"os/user"

	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-wikiparse"
	log "github.com/sirupsen/logrus"
)

// var Mode = flag.Int("m", 0, "Input mode 1,")

// func main() {

// 	switch *Mode {
// 	case 1:
// 		log.Println("test")
// 		mag.MagNetrans()
// 	default:
// 		log.Println("do nothing")

// 	}
// }
var (
	DataPath string
)

func calData(chanout chan []string) {
	for receive := range chanout {
		id := receive[0]
		// title := receive[1]
		text := receive[2]
		fmt.Println("原始单词数", len(strings.Fields(text)))
		fmt.Println("原始字母数", len(text))
		wikiPath := DataPath + id + ".wiki"
		htmlPath := DataPath + id + ".html"
		err := ioutil.WriteFile(wikiPath, []byte(text), 0666)
		if err != nil {
			fmt.Println(err)
		}

		cmd := exec.Command("/usr/bin/pandoc", "-f", "mediawiki", "-t", "html5", "-s", wikiPath, "-o", htmlPath)
		// 执行
		err = cmd.Run()
		if err != nil {
			log.Println(wikiPath, err)
		}
		// fmt.Println(err)
		filer, _ := os.Open(htmlPath)
		doc, err := goquery.NewDocumentFromReader(filer)
		if err != nil {
			log.Fatal(err)
		}
		text = doc.Text()
		fmt.Println("简化后字母数", len(strings.Fields(text)))

		filer.Close()
		os.Remove(wikiPath)
		os.Remove(htmlPath)
		// fmt.Println(text)
	}

}

func main() {

	userInfo, _ := user.Current()
	DataPath = "/run/user/" + userInfo.Uid + "/wiki/"
	os.Mkdir(DataPath, 0755)

	filePath := "/data/ssdj/download/enwiki-20200501-pages-articles.xml.bz2"
	// filePath :="/data/ssdj/download/t.txt.gz"

	fi, err := os.Open(filePath)
	if err != nil {
		log.Fatal("ERROR:", err)
	}

	defer fi.Close()

	fz := bzip2.NewReader(fi)

	p, err := wikiparse.NewParser(fz)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	chanin := make(chan []string)

	go calData(chanin)
	go calData(chanin)
	go calData(chanin)
	go calData(chanin)
	for err == nil {
		var page *wikiparse.Page
		page, err = p.Next()
		if err == nil {
			for _, rev := range page.Revisions {
				chanin <- []string{strconv.Itoa(int(page.ID)), page.Title, rev.Text}
			}
		} else {
			log.Println(err)
		}
		// time.Sleep(2 * time.Second)
	}
}
