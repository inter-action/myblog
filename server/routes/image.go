package routes

import (
	"net/http"
	"os"
	"path"

	"github.com/inter-action/myblog/server/db"
	"github.com/inter-action/myblog/server/utils"
	"github.com/labstack/echo"
)

func CreateImageRoutes(mdroot string) func(echo.Context) error {
	if mdroot == "" {
		panic("mdroot required")
	}

	return func(c echo.Context) error {
		db.LoadArticles(mdroot)

		article := db.ArticleMap[c.Param("slug")]
		if article == nil {
			return c.String(http.StatusNotFound, "no corresponding article with that slug")
		}

		request := c.Request()
		urlPath := request.URL.Path
		imageRelatviePath := urlPath[len("/images/"+c.Param("slug")+"/"):]
		imagePath := path.Join(path.Dir(article.Path), imageRelatviePath)
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			return c.String(http.StatusNotFound, "no file to path: "+imageRelatviePath)
		}
		if reader, err := utils.FileBufferedReader(imagePath); err != nil {
			c.Echo().Logger.Errorf("failed to send file, %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.Stream(http.StatusOK, "image/*", reader)
		}
	}
}
