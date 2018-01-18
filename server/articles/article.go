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

type Article struct {
	Title       string    `json:"title"`
	Path        string    `json:"-"`
	DateCreated time.Time `json:"created"`
	DateUpdated time.Time `json:"updated"`
	Tags        []string  `json:"tags"`
	Slug        string    `json:"slug"`
	RawContent  string    `json:"-"`
	MetaContent string    `json:"-"`
	cacheTime   time.Time
}

// LoadContent return []string of two, 1st is the meta 2st is content of the file
func (self *Article) LoadContent() ([]string, error) {
	if self.RawContent == "" || utils.IsOutCache(self.cacheTime) {
		fileContent, err := dry.FileGetString(self.Path)
		if err != nil {
			return nil, err
		}
		meta, content := splitMarkdown(fileContent)
		self.RawContent = content
		self.MetaContent = meta
		self.cacheTime = time.Now()
	}
	return []string{self.MetaContent, self.RawContent}, nil
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

func splitMarkdown(str string) (string, string) {
	regx := regexp.MustCompile(`(?s)=+([^=]*)=+(.*)`)
	matchs := regx.FindAllStringSubmatch(str, 1)
	if len(matchs) == 1 {
		return strings.TrimSpace(matchs[0][1]), strings.TrimSpace(matchs[0][2])
	}
	return "", str
}
