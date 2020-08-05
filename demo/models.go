package demo

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
