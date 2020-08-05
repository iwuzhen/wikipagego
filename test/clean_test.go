package test

import (
	"io/ioutil"
	"testing"

	"github.com/iwuzhen/wikipagego/clean/wikitext"
)

func Test_Clean(t *testing.T) {
	bytes, err := ioutil.ReadFile("./data.wiki")
	if err != nil {
		t.Log(err)
	}
	source := string(bytes)
	ret := wikitext.Clean(&source)
	t.Error(*ret)
	// t.Error("error")
}
