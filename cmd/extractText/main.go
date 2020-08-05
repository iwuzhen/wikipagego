/*
 * @Version: 0.0.1
 * @Author: ider
 * @Date: 2020-05-11 18:19:04
 * @LastEditors: ider
 * @LastEditTime: 2020-08-04 21:37:37
 * @Description:
 */

package main

import (
	"os"
	"regexp"

	// "database/sql"
	// "io/ioutil"
	"io"
	// "time"
	"compress/bzip2"
	"encoding/xml"

	log "github.com/sirupsen/logrus"

	// "fmt"
	"strings"

	"github.com/caarlos0/env"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type config struct {
	PgUri string `env:"PGURI" envDefault:"postgres://postgres:postgres@192.168.1.220/wiki_knogen?sslmode=disable"`
}

type redirect struct {
	Title string `xml:"title,attr"`
}

type Page struct {
	Title    string   `xml:"title"`
	Redirect redirect `xml:"redirect"`
	Text     string   `xml:"revision>text"`
}

func connectPostgres(cfg config) *sqlx.DB {
	db, err := sqlx.Open("postgres", cfg.PgUri)
	if err != nil {
		log.Fatal("ERROR:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("ping 失败", err)
	}
	return db
}

func initDB(db *sqlx.DB, table_name string) {

	// 初始化表
	// DROP TABLE IF EXISTS ` + table_name + `;
	sql := `
CREATE TABLE ` + table_name + `(
	id SERIAL PRIMARY KEY,
	title VARCHAR UNIQUE ,
	text VARCHAR NOT NULL
);
comment on column ` + table_name + `.id is '主键ID，自增';
`
	// CREATE UNIQUE INDEX unique_title on ` + table_name + `(title);
	db.Exec(sql)

}

func render_b(source *string) *string {
	re := regexp.MustCompile(`'{3,}(.*?)'{3,}`)
	for _, value := range re.FindAllString(*source, -1) {
		// fmt.Println(value)
		*source = strings.ReplaceAll(*source, value, "<b>"+strings.Trim(value, "'")+"</b>")
	}
	return source
}

func render_link(source *string) *string {
	re := regexp.MustCompile(`\[\[(.*?)\]\]`)
	for _, value := range re.FindAllString(*source, -1) {
		rep_value := strings.Trim(value, "[]")
		if strings.Index(rep_value, "|") > -1 {
			ret_list := strings.Split(rep_value, `|`)
			rep_value = ret_list[len(ret_list)-1]
		}
		*source = strings.ReplaceAll(*source, value, "<a href='#'>"+rep_value+"</a>")
	}
	return source
}

func remove_refer(source *string) *string {
	re := regexp.MustCompile(`<ref>(.*?)</ref>`)
	for _, value := range re.FindAllString(*source, -1) {
		*source = strings.ReplaceAll(*source, value, "")
	}
	return source
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)
}

func handleText(text *string) *string {
	ret_row := ""
	no_dot_row := ""
	ret := strings.Split(*text, "\n")
	for _, value := range ret {
		if strings.HasPrefix(value, "'''") {
			ret_row = value
			break
		}
		if no_dot_row == "" && !strings.HasPrefix(value, "<") && !strings.HasPrefix(value, "=") && !strings.HasPrefix(value, "|") && !strings.HasPrefix(value, "[[") && !strings.HasPrefix(value, "]]") && !strings.HasPrefix(value, "{{") && !strings.HasPrefix(value, "}}") && value != "" {
			no_dot_row = value
		}
	}
	if ret_row != "" {
		return &ret_row
	}
	return &no_dot_row
}

func main() {
	table_name := "wiki_title"

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Warning("%+v\n", err)
	}
	// log.Info(cfg.TempFolder)
	db := connectPostgres(cfg)
	defer db.Close()
	initDB(db, table_name)

	sql := `insert into ` + table_name + ` (title, text)
	values ($1, $2)
	ON CONFLICT (title) 
	DO UPDATE SET text = EXCLUDED.text;
`
	tx, _ := db.Beginx()
	stmt, _ := tx.Prepare(sql)

	filePath := "/data/ssdj/download/enwiki-20200501-pages-articles.xml.bz2"

	fi, err := os.Open(filePath)
	if err != nil {
		log.Fatal("ERROR:", err)
	}
	defer fi.Close()

	fz := bzip2.NewReader(fi)

	d := xml.NewDecoder(fz)
	count := 0
	for {
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
				// log.Info("page ",)
				if page.Redirect.Title == "" {
					ret := handleText(&(page.Text))

					// log.Info("page ",*ret)
					if *ret == "" {
					} else {
						// remove_refer(ret)
						// render_b(ret)
						// render_link(ret)
						// log.Info("text: ",*ret)
						// log.Info("title: ",page.Title)
						_, err = stmt.Exec(page.Title, *ret)
						if err != nil {
							log.Fatal("ERROR:", err)
						}
						if count%1000 == 0 {
							tx.Commit()
							tx, _ = db.Beginx()
							stmt, _ = tx.Prepare(sql)
							log.Info(count)
						}
						count++
					}
				}
			}
		}
	}
	tx.Commit()
	log.Info("over ", count)
}
