package controllers

import (
    "github.com/revel/revel"
    "github.com/janekolszak/revmgo"
    "github.com/ip4368/colonnade/app/models"
)

type Coordinators struct {
	*revel.Controller
	revmgo.MongoController
}

func (c Coordinators) AddStages(Id string) revel.Result {
	var stages struct{
		Description []string `json:"d"`
	}
	models.ParseBody(c.Request.Body, &stages)

	models.CoordinatorAddStages(
		c.MongoSession,
		models.User_t{
			Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
		},
		Id,
		stages.Description,
	)
	return c.RenderJson([]int{})
}

func (c Coordinators) AddTasks(Id string) revel.Result {
	var tasks struct{
		Description []string `json:"d"`
	}
	models.ParseBody(c.Request.Body, &tasks)

	models.CoordinatorAddTasks(
		c.MongoSession,
		models.User_t{
			Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
		},
		Id,
		tasks.Description,
	)
	return c.RenderJson([]int{})
}