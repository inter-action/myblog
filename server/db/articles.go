package db

import (
	"errors"
	"sync"
	"time"

	"github.com/inter-action/myblog/server/articles"
	"github.com/inter-action/myblog/server/utils"
)

var Articles articles.Articles
var ArticlesTime time.Time
var ArticleMap map[string]*articles.Article

var mutex = &sync.Mutex{}

func LoadArticles(mdroot string) error {
	if Articles == nil || utils.IsOutCache(ArticlesTime) {
		mutex.Lock()
		defer mutex.Unlock()

		path := mdroot
		if path == "" {
			return errors.New("mdroot not config")
		}

		var err error
		Articles, err = articles.ParseArticles(path)
		if err != nil {
			return err
		}
		ArticleMap = make(map[string]*articles.Article)
		for _, ar := range Articles {
			ArticleMap[ar.Slug] = ar
		}
		ArticlesTime = time.Now()
	}

	return nil
}
