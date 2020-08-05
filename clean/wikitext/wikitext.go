/*
 * @Version: 0.0.1
 * @Author: ider
 * @Date: 2020-08-04 19:03:07
 * @LastEditors: ider
 * @LastEditTime: 2020-08-05 10:50:30
 * @Description:清理 wikiText,留存正常页面
 */
package wikitext

import (
	"regexp"
	"strings"
)

func Clean(source *string) *string {
	/**
	 * @description: 清理标签至可读
	 * @param {type}
	 * @return {type}
	 */

	source = remove_comment(source)
	source = remove_refer(source)
	source = clean_blood(source)
	source = clean_head(source)
	source = clean_link(source)
	source = clean_empty_line(source)

	return source
}

func clean_blood(source *string) *string {
	re := regexp.MustCompile(`'{2,}?(.*?)'{2,}?`)
	for _, value := range re.FindAllString(*source, -1) {
		*source = strings.ReplaceAll(*source, value, strings.Trim(value, "'"))
	}
	return source
}

func clean_head(source *string) *string {
	re := regexp.MustCompile(`={2,}(.*?)={2,}`)
	for _, value := range re.FindAllString(*source, -1) {
		*source = strings.ReplaceAll(*source, value, strings.Trim(value, "="))
	}
	return source
}

func clean_link(source *string) *string {
	re := regexp.MustCompile(`\[{2}([^\[]*?)\]{2}`)
	flag := 0
	for _, value := range re.FindAllString(*source, -1) {
		flag = 1
		arrayS := strings.Split(strings.Trim(value, "[]"), "|")
		var repString = arrayS[0]
		if len(arrayS) > 1 && arrayS[len(arrayS)-1] != "" {
			repString = arrayS[len(arrayS)-1]
		}
		*source = strings.ReplaceAll(*source, value, repString)
	}
	if flag == 1 {
		source = clean_link(source)
	}
	return source
}

func render_b(source *string) *string {
	re := regexp.MustCompile(`'{3,}(.*?)'{3,}`)
	*source = re.ReplaceAllString(*source, "")
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

func remove_comment(source *string) *string {
	re := regexp.MustCompile(`&lt;!--(.*?)--&gt;`)
	*source = re.ReplaceAllString(*source, "")
	// 重复
	re = regexp.MustCompile(`<!--(.*?)-->`)
	*source = re.ReplaceAllString(*source, "")

	re = regexp.MustCompile(`\{\{(.*?)\}\}`)
	*source = re.ReplaceAllString(*source, "")
	return source
}

func remove_refer(source *string) *string {
	re := regexp.MustCompile(`&lt;ref(.*?)/&gt;`)
	*source = re.ReplaceAllString(*source, "")
	re = regexp.MustCompile(`&lt;ref(.*?)ref&gt;`)
	*source = re.ReplaceAllString(*source, "")
	// 重复
	re = regexp.MustCompile(`<ref(.*?)/>`)
	*source = re.ReplaceAllString(*source, "")
	re = regexp.MustCompile(`<ref(.*?)ref>`)
	*source = re.ReplaceAllString(*source, "")
	return source
}

func clean_empty_line(source *string) *string {
	re := regexp.MustCompile(`(?m)\*`)
	*source = re.ReplaceAllString(*source, "")

	re = regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	*source = re.ReplaceAllString(*source, "")

	return source
}
