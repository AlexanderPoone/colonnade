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

func (c Admins) IsAdmin() revel.Result {
	loginStat := models.LoginStatus(
        c.Session["email"],
        c.Session["username"],
        c.Session["name"],
        c.Session["userId"],
    )

    var result int = 0
    if loginStat == 0 {
    	result = models.CheckAdmin(c.MongoSession, c.Session["userId"])
    } else { result = 1 }

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
        case 3 :
            data["message"] = "Invalid User ID"
    }
    return c.RenderJson(data)
}
