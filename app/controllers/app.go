package controllers

import (
	"bitlyclone/app"
	"bitlyclone/app/models"
	"bitlyclone/app/services"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) RegisterRoute(URL string) revel.Result {
	c.Validation.Required(URL)
	c.Validation.URL(URL)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect("/")
	}

	c.Log.Info(URL)

	route := models.Route{
		OriginalUrl: URL,
		ShortPath:   services.BuildRandomString(),
	}

	err := app.DB.Set(route.ShortPath, route, 0).Err()

	if err != nil {
		c.Log.Error(err.Error())
	}

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
