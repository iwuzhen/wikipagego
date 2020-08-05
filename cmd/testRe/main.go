/*
 * @Version: 0.0.1
 * @Author: ider
 * @Date: 2020-05-12 13:46:05
 * @LastEditors: ider
 * @LastEditTime: 2020-05-12 15:31:56
 * @Description: 
 */
package main
import (
	"fmt"
	 "regexp"
	 "strings"
)


func render_b(source *string) *string {
	re := regexp.MustCompile(`'{3,}(.*?)'{3,}`)
	for _,value  := range re.FindAllString(*source, -1){
		fmt.Println(value)
		*source = strings.ReplaceAll(*source,value,"<b>"+ strings.Trim(value, "'") + "</b>")
	}
	return source
}

func render_link(source *string) *string {
	re := regexp.MustCompile(`\[\[(.*?)\]\]`)
	for _,value  := range re.FindAllString(*source, -1){
		rep_value := strings.Trim(value, "[]")
		if strings.Index(rep_value, "|") > -1 {
			ret_list := strings.Split(rep_value, `|`)
			rep_value = ret_list[len(ret_list)-1]
		}
		*source = strings.ReplaceAll(*source,value,"<a href='#'>"+ rep_value + "</a>")
	}
	return source
}


func remove_refer(source *string) *string {
	re := regexp.MustCompile(`<ref>(.*?)</ref>`)
	for _,value  := range re.FindAllString(*source, -1){
		*source = strings.ReplaceAll(*source,value,"")
	}
	return source
}

func main (){
	s := `
'''Alabama''' ({{IPAc-en|,|æ|l|ə|'|b|æ|m|ə|}}) is a [[U.S. state|state]] in the [[Southern United States|southeastern region]] of the [[United States]]. It is bordered by [[Tennessee]] to the north, [[Georgia (U.S. state)|Georgia]] to the east, [[Florida]] and the [[Gulf of Mexico]] to the south, and [[Mississippi]] to the west. Alabama is the [[List of U.S. states and territories by area|30th largest by area]] and the [[List of U.S. states and territories by population|24th-most populous]] of the [[List of U.S. states|U.S. states]]. With a total of {{convert|1500|mi|km}} of [[inland waterway]]s, Alabama has among the most of any state.<ref>{{cite web |title=Alabama Transportation Overview |url=https://www.edpa.org/wp-content/uploads/Alabama-Transportation-Overview-1.pdf |publisher=Economic Development Partnership of Alabama |accessdate=January 21, 2017 |archive-url=https://web.archive.org/web/20181113075704/https://www.edpa.org/wp-content/uploads/Alabama-Transportation-Overview-1.pdf |archive-date=November 13, 2018 |url-status=dead }}</ref>
`
	fmt.Println(s)
	render_b(&s)
	render_link(&s)
	remove_refer(&s)
	fmt.Println(s)
	// strings.Trim("¡¡¡Hello, Gophers!!!", "!¡")

}