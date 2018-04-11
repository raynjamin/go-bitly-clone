package controllers

import (
	"bitlyclone/app"
	"bitlyclone/app/models"
	"bitlyclone/app/services"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) RegisterRoute(url string) revel.Result {
	c.Validation.Required(url)
	c.Validation.URL(url)
	c.Validation.Keep()

	if c.Validation.HasErrors() {
		c.FlashParams()
		return c.Redirect("/")
	}

	c.Log.Info(url)

	shortString := GetUniqueShortPath()

	route := models.Route{
		OriginalUrl: url,
		ShortPath:   shortString,
	}

	app.DB.Get(route.ShortPath)
	err := app.DB.Set(route.ShortPath, route, 0).Err()

	if err != nil {
		c.Log.Error(err.Error())
	}

	newURL := c.Request.URL.Host + "/" + route.ShortPath

	c.Flash.Success(fmt.Sprintf("URL creation success: %s -> %s", url, newURL))
	return c.Redirect("/")
}

func (c App) RandoPath(randopath string) revel.Result {
	val, err := app.DB.Get(randopath).Result()

	if err != nil {
		c.Validation.Error("SHORTENED URL NOT FOUNDSIES")
		c.Validation.Keep()
		return c.Redirect("/")
	}

	r := &models.Route{}

	r.UnMarshalBinary([]byte(val))

	return c.Redirect(r.OriginalUrl)
}

// TODO: Clean this function up majorly.
func GetUniqueShortPath() string {
	shortUrl := services.BuildRandomString()

	_, err := app.DB.Get(shortUrl).Result()

	for err != redis.Nil {
		shortUrl := services.BuildRandomString()
		_, err = app.DB.Get(shortUrl).Result()
	}

	return shortUrl
}
