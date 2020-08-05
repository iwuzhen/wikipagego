package demo

import (
	"os"
	"compress/bzip2"
)

func CountText() {
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
				// log.Info("title ",page.Title)
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
