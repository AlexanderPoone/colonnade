package controllers

import (
    "github.com/revel/revel"
    "encoding/json"
    "io/ioutil"
    "github.com/janekolszak/revmgo"
    "github.com/ip4368/colonnade/app/models"
)

func init() {
    revmgo.ControllerInit()
}

type Courses struct {
    *revel.Controller
    revmgo.MongoController
}

func (c Courses) CoursesForUser() revel.Result {
    loginStat := models.LoginStatus(
        c.Session["email"],
        c.Session["username"],
        c.Session["name"],
        c.Session["userId"],
    )
    var result int = 0
    var retData []models.Course_t
    if loginStat == 0 {
    	result, retData = models.CoursesForUser(c.MongoSession, c.Session["userId"])
    } else { result = 1 }

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "Sucess"
            data["data"] = retData
        case 1 :
            data["message"] = "User has not logged in"
        case 2 :
            data["message"] = "Invalid User ID"
        case 3 :
            data["message"] = "Unexpected Error in Database"
    }
    return c.RenderJson(data)
}

func (c Courses) Course() revel.Result {
	var r RegisterProfile
    var bodyBytes []byte
    if c.Request.Body != nil {
        bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
    }
    json.Unmarshal([]byte(bodyBytes), &r)

    //result := models.RegisterHandler(c.MongoSession, r.Email, r.Username, r.Password, r.Name)
    result := 0

    // start with initialise response interface
    data := make(map[string]interface{})
    //data["error"] = result
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "Sucess"
    }
    return c.RenderJson(data)
}
