package controllers

import (
	"github.com/revel/revel"
    //"encoding/json"
    //"io/ioutil"
    "github.com/janekolszak/revmgo"
    "github.com/ip4368/colonnade/app/models"
)

func init() {
    revmgo.ControllerInit()
    models.GuardAdmins()
}

type Admins struct {
    *revel.Controller
    revmgo.MongoController
}

func (c Admins) CheckAdmin() revel.Result {
    result := models.CheckAdmin(
		c.MongoSession,
		c.Session["email"],
        c.Session["username"],
        c.Session["name"],
        c.Session["userId"],
	)

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "User is admin"
            c.Session["admin"] = "t"
        case 1 :
            data["message"] = "User is not logged in"
        case 2 :
            data["message"] = "User is not admin"
        case 3 :
            data["message"] = "Invalid User ID"
    }
    return c.RenderJson(data)
}

func (c Admins) IsAdmin() revel.Result {
	result := models.IsAdmin(
		c.Session["email"],
        c.Session["username"],
        c.Session["name"],
        c.Session["userId"],
		c.Session["admin"],
	)

	// start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "User is admin"
        case 1 :
            data["message"] = "User is not logged in"
        case 2 :
            data["message"] = "User is not admin"
    }
    return c.RenderJson(data)
}

/*func (c Admins) Courses() revel.Result {
	
}*/
