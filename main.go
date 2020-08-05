package main

import (
	"compress/bzip2"
	"fmt"
	"os"
	"strconv"
	"strings"

	"os/user"

	"github.com/dustin/go-wikiparse"
	log "github.com/sirupsen/logrus"

	"github.com/iwuzhen/wikipagego/clean/wikitext"
	"github.com/iwuzhen/wikipagego/config"
	pgc "github.com/iwuzhen/wikipagego/database/postgres"
	"gorm.io/gorm/clause"
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
		// id := receive[0]
		title := receive[1]
		text := receive[2]
		// fmt.Println("原始单词数", len(strings.Fields(text)))
		// fmt.Println("原始字母数", len(text))
		CharCount := int32(len(text))
		worldCount := int32(len(strings.Fields(text)))
		newcur := wikitext.Clean(&text)
		// fmt.Println("精简后的", len(strings.Fields(*newcur)))
		// fmt.Println("精简后的", *newcur)
		doc := pgc.WikiWordCount{
			Title:          title,
			CharCount:      CharCount,
			WordCount:      worldCount,
			CleanWordCount: int32(len(strings.Fields(*newcur))),
		}
		// pgc.WikiDB.Save(&doc)
		pgc.WikiDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&doc)

	}
}

func main() {

	cfg := config.GetConfig()
	pgc.NewWikiDBConn(cfg.PgConn)

	userInfo, _ := user.Current()
	DataPath = "/run/user/" + userInfo.Uid + "/wiki/"
	os.Mkdir(DataPath, 0755)

	// filePath := "/data/ssdj/download/enwiki-20200501-pages-articles.xml.bz2"
	filePath := "/data/ssdj/download/enwiki-20200801-pages-articles.xml.bz2"

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

	for i := 0; i < 20; i++ {
		go calData(chanin)
	}
	i := 0
	for err == nil {
		var page *wikiparse.Page
		page, err = p.Next()
		if err == nil {
			for _, rev := range page.Revisions {
				i += 1
				chanin <- []string{strconv.Itoa(int(page.ID)), page.Title, rev.Text}
			}
		} else {
			log.Println(err)
		}
		// time.Sleep(1 * time.Second)
		if i%1000 == 0 {
			fmt.Println(i)
		}
	}
	close(chanin)
}
