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

	result := models.CoordinatorAddStages(
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

	data := make(map[string]interface{})
	data["error"] = result
	switch result{
		case 0:
			data["message"] = "Stages have been added"
		case 1:
			data["message"] = "Course ID is invalid"
		case 2:
			data["message"] = "Your login session is not valid"
		case 3:
			data["message"] = "Course is not found"
	}
	return c.RenderJson(data)
}

func (c Coordinators) AddTasks(Id string, Stage int) revel.Result {
	var tasks struct{
		Description []string `json:"d"`
	}
	models.ParseBody(c.Request.Body, &tasks)

	result := models.CoordinatorAddTasks(
		c.MongoSession,
		models.User_t{
			Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
		},
		Id,
		Stage,
		tasks.Description,
	)
	
	data := make(map[string]interface{})
	data["error"] = result
	switch result{
		case 0:
			data["message"] = "Stages have been added"
		case 1:
			data["message"] = "Course ID is invalid"
		case 2:
			data["message"] = "Your login session is not valid"
		case 3:
			data["message"] = "Course is not found"
	}
	return c.RenderJson(data)
}