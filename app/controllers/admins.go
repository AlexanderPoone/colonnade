package controllers

import (
    "github.com/revel/revel"
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
        models.User_t{
            Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
        },
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
        models.User_t{
            Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
        },
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

func (c Admins) Courses() revel.Result {
    var courseIdHex struct {Id string `json:"Id,omitempty"`}
    models.ParseBody(c.Request.Body, &courseIdHex)

    result, courses := models.AdminCourses(
        c.MongoSession,
        models.User_t{
            Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
        },
        c.Session["admin"],
        courseIdHex.Id,
    )

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "Success"
            data["data"] = make(map[string]interface{})
            data["data"].(map[string]interface{})["data"] = courses
        case 1 :
            data["message"] = "User is not admin"
        case 2 :
            data["message"] = "Course Id is not valid"
        case 3 :
            data["message"] = "Unexpected Error in Database"
    }
    return c.RenderJson(data)
}

func (c Admins) NewCourse() revel.Result {
    // read request body to byte
    var course models.Course_t
    models.ParseBody(c.Request.Body, &course)

    result, idHex := models.AdminNewCourse(
        c.MongoSession,
        models.User_t{
            Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
        },
        c.Session["admin"],
        course)

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "Course has been created"
            data["data"] = make(map[string]interface{})
            data["data"].(map[string]interface{})["courseId"] = idHex
        case 1 :
            data["message"] = "User is not admin"
        case 2 :
            data["message"] = "Unexpected Error in Database"
    }
    return c.RenderJson(data)
}
