package main

import (
	"net/http"
	"path"

	"github.com/inter-action/myblog/server/articles"
	"github.com/inter-action/myblog/server/db"
	"github.com/labstack/echo"
	"github.com/olebedev/config"
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
	group.GET("/meta", self.MetaHandler)
}

// ConfHandler handle the app config, for example
func (self *API) ConfHandler(c echo.Context) error {
	return c.JSON(200, app.Conf.Root)
}

func (self *API) ArticlesHandler(c echo.Context) error {
	db.LoadArticles(app.Conf.UString("mdroot"))
	return c.JSON(200, db.Articles)
}

func (self *API) ArticleHandler(c echo.Context) error {
	db.LoadArticles(app.Conf.UString("mdroot"))
	result := db.ArticleMap[c.Param("slug")]
	if result != nil {
		if contents, err := result.LoadContent(); err != nil {
			c.Error(err)
		} else {
			data := struct {
				Article *articles.Article `json:"article"`
				Content string            `json:"content"`
			}{
				Article: result,
				Content: contents[1],
			}
			return c.JSON(200, data)
		}
	} else {
		return c.String(http.StatusNotFound, "")
	}

	return nil
}

func (self *API) MetaHandler(c echo.Context) error {
	mdroot := app.Conf.UString("mdroot")
	metaPath := path.Join(mdroot, "meta.yaml")
	if yamlConfig, err := config.ParseYamlFile(metaPath); err != nil {
		c.Echo().Logger.Errorf("meta.json, %s", err.Error())
		return c.String(http.StatusInternalServerError, "meta.json: "+err.Error())
	} else {
		// :bm, golang cant convert []interface{} to []string
		return c.JSON(200, struct {
			Me     string        `json:"me"`
			Images []interface{} `json:"images"`
		}{
			Me:     yamlConfig.UString("me"),
			Images: yamlConfig.UList("images"),
		})
	}
}
