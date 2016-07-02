package controllers

import "github.com/revel/revel"

type App struct {
	*revel.Controller
}

type Session map[string]string

func (c App) Index() revel.Result {
	return c.Render()
}
