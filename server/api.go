package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/inter-action/myblog/server/articles"
	"github.com/inter-action/myblog/server/utils"
	"github.com/labstack/echo"
)

// API is a defined as struct bundle
// for api. Feel free to organize
// your app as you wish.
type API struct{}

// Bind attaches api routes
func (self *API) Bind(group *echo.Group) {
	group.GET("/v1/conf", self.ConfHandler)
	group.GET("/articles", self.ArticlesHandler)
	group.GET("/articles/:slug", self.ArticleHandler)
}

// ConfHandler handle the app config, for example
func (self *API) ConfHandler(c echo.Context) error {
	return c.JSON(200, app.Conf.Root)
}

var _articles articles.Articles
var _articlesTime time.Time
var _articleMap map[string]*articles.Article

func (self *API) ArticlesHandler(c echo.Context) error {
	loadArticles()
	return c.JSON(200, _articles)
}

func (self *API) ArticleHandler(c echo.Context) error {
	loadArticles()
	result := _articleMap[c.Param("slug")]
	if result != nil {
		if contents, err := result.LoadContent(); err != nil {
			c.Error(err)
		} else {
			return c.String(200, contents[1])
		}
	} else {
		return c.String(http.StatusNotFound, "")
	}

	return nil
}

func loadArticles() error {
	if _articles == nil || utils.IsOutCache(_articlesTime) {
		path := app.Conf.UString("mdroot")
		if path == "" {
			return errors.New("mdroot not config")
		}

		var err error
		_articles, err = articles.ParseArticles(path)
		if err != nil {
			return err
		}
		_articleMap = make(map[string]*articles.Article)
		for _, ar := range _articles {
			_articleMap[ar.Slug] = ar
		}
		_articlesTime = time.Now()
	}

	return nil
}
