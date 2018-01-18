package articles

import (
	"strings"
	"testing"
)

func TestParseArticle(t *testing.T) {
	_, err := ParseArticle("./__test__/article.md")
	if err != nil {
		t.Error("failed miserbley:", err)
	}
}

func TestSplitContent(t *testing.T) {
	fileContent := `
===
title: you dont know js
created: 2018-01-18 14:20:06
===

#title
    `
	meta, content := splitMarkdown(fileContent)

	if meta != strings.TrimSpace(`
title: you dont know js
created: 2018-01-18 14:20:06
    `) || content != "#title" {
		t.Errorf("test failed, meta: %s, content: %s", meta, content)
	}
}

func TestReadMetaData(t *testing.T) {
	yaml := `
title: you dont know js
created: 2018-01-18 14:20:06
tags: machine_learning, es
`
	ar := Article{}
	err := readMetaDataInto(yaml, &ar)
	if err != nil {
		t.Error(err)
	}
	if ar.Title != "you dont know js" {
		t.Error("err: ", ar.Title)
	}

	if ar.Slug != "you-dont-know-js" {
		t.Error("err: ", ar.Slug)
	}

	if len(ar.Tags) != 2 {
		t.Error("err: ", ar.Tags)
	}
}
