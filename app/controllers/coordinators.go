package controllers

import (
    "github.com/revel/revel"
    "github.com/janekolszak/revmgo"
    //"github.com/ip4368/colonnade/app/models"
)

type Coordinators struct {
	*revel.Controller
	revmgo.MongoController
}

func (c Coordinators) AddStages() revel.Result {
	return c.RenderJson([]int{})
}

func (c Coordinators) AddTasks() revel.Result {
	return c.RenderJson([]int{})
}