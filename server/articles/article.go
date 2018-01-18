package articles

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/inter-action/myblog/server/utils"
	"github.com/olebedev/config"
	dry "github.com/ungerik/go-dry"
)

type Article struct {
	Title       string
	Path        string
	DateCreated time.Time
	DateUpdated time.Time
	Tags        []string
	Slug        string
}

func ParseArticle(path string) (*Article, error) {
	fileContent, err := dry.FileGetString(path)
	if err != nil {
		return nil, err
	}
	ar := Article{}
	meta, _ := splitMarkdown(fileContent)
	if err := readMetaDataInto(meta, &ar); err != nil {
		return nil, err
	}
	if fileInfo, err := os.Stat(path); err != nil {
		return nil, err
	} else {
		ar.DateUpdated = fileInfo.ModTime()
	}
	if path, err := filepath.Abs(path); err != nil {
		return nil, err
	} else {
		ar.Path = path
	}
	return &ar, nil
}

// LoadContent return []string of two, 1st is the meta 2st is content of the file
func (ar *Article) LoadContent() ([]string, error) {
	fileContent, err := dry.FileGetString(ar.Path)
	if err != nil {
		return nil, err
	}
	meta, content := splitMarkdown(fileContent)
	return []string{meta, content}, nil
}

func splitMarkdown(str string) (string, string) {
	regx := regexp.MustCompile(`(?s)=+([^=]*)=+(.*)`)
	matchs := regx.FindAllStringSubmatch(str, 1)
	if len(matchs) == 1 {
		return strings.TrimSpace(matchs[0][1]), strings.TrimSpace(matchs[0][2])
	}
	return "", str
}

func readMetaDataInto(str string, ar *Article) error {
	conf, err := config.ParseYaml(str)
	if err != nil {
		return err
	}

	ar.Title = strings.TrimSpace(conf.UString("title"))
	ar.Slug = strings.TrimSpace(conf.UString("slug", slug.Make(ar.Title)))
	if conf.UString("created") != "" {
		if ar.DateCreated, err = utils.ParseTime(conf.UString("created")); err != nil {
			return err
		}
	}
	if conf.UString("tags") != "" {
		ar.Tags = strings.Split(conf.UString("tags"), ",")
	}
	return nil
}
