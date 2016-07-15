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

func (c Admins) Courses(p int) revel.Result {
    result, courses := models.AdminCourses(
        c.MongoSession,
        models.User_t{
            Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
        },
        c.Session["admin"],
        p,
    )

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "Success"
            data["data"] = make(map[string]interface{})
            data["data"].(map[string]interface{})["courses"] = courses
        case 1 :
            data["message"] = "User is not admin"
        case 2 :
            data["message"] = "Unexpected Error in Database"
    }
    return c.RenderJson(data)
}

func (c Admins) Users(p int) revel.Result {
    result, users := models.AdminUsers(
        c.MongoSession,
        models.User_t{
            Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
        },
        c.Session["admin"],
        p,
    )

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "Success"
            data["data"] = make(map[string]interface{})
            data["data"].(map[string]interface{})["users"] = users
        case 1 :
            data["message"] = "User is not admin"
        case 2 :
            data["message"] = "Unexpected Error in Database"
    }
    return c.RenderJson(data)
}

func (c Admins) Course(Id string) revel.Result {
    result, course := models.AdminCourse(
        c.MongoSession,
        models.User_t{
            Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
        },
        c.Session["admin"],
        Id,
    )

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "Success"
            data["data"] = make(map[string]interface{})
            data["data"].(map[string]interface{})["course"] = course
        case 1 :
            data["message"] = "User is not admin"
        case 2 :
            data["message"] = "Course Id is not valid"
        case 3 :
            data["message"] = "Unexpected Error in Database"
    }
    return c.RenderJson(data)
}

func (c Admins) User(Id string) revel.Result {
    result, user := models.AdminUser(
        c.MongoSession,
        models.User_t{
            Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
        },
        c.Session["admin"],
        Id,
    )

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "Success"
            data["data"] = make(map[string]interface{})
            data["data"].(map[string]interface{})["user"] = user
        case 1 :
            data["message"] = "Request user is not admin"
        case 2 :
            data["message"] = "User Id is not valid"
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

func (c Admins) AddUsers2Course(Id string) revel.Result {
    var users struct{
        Users  []models.UserInCourse_t  `json:"users"`
    }
    models.ParseBody(c.Request.Body, &users)

    status, usersStatus := models.AdminAddUser2Course(
        c.MongoSession,
        models.User_t{
            Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
        },
        c.Session["admin"],
        Id,
        users.Users)

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = status
    switch status {
        case 0 :
            data["message"] = "All/Part of the users have been added"
            data["data"] = make(map[string]interface{})
            data["data"].(map[string]interface{})["userStatus"] = usersStatus
        case 1 :
            data["message"] = "User is not admin"
        case 2 :
            data["message"] = "Course id is invalid"
        case 3 :
            data["message"] = "No valid users can be added"
        case 4 :
            data["message"] = "Unexpected Error in Database"
    }
    return c.RenderJson(data)
}

func (c Admins) UpdateCourse(Id string) revel.Result {
    var details models.Details_t
    models.ParseBody(c.Request.Body, &details)
    result, allRes := models.AdminUpdateCourse(
        c.MongoSession,
        models.User_t{
            Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
        },
        c.Session["admin"],
        Id,
        details,
    )

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "All/Part of the details have been updated"
            data["data"] = make(map[string]interface{})
            data["data"].(map[string]interface{})["results"] = allRes
        case 1 :
            data["message"] = "User is not admin"
        case 2 :
            data["message"] = "Course Id is not valid"
        case 3 :
            data["message"] = "None of the detail has been updated due to unexpected Error in Database"
    }
    return c.RenderJson(data)
}

func (c Admins) FindUserByIdentifier(q, s string) revel.Result {
    allowSuspend := true
    if s == "f" {allowSuspend = false}
    status, result := models.AdminGetUserByIdentifier(
        c.MongoSession,
        models.User_t{
            Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
        },
        c.Session["admin"],
        q,
        allowSuspend)

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = status
    switch status {
        case 0 :
            data["message"] = "Result found"
            data["data"] = make(map[string]interface{})
            data["data"].(map[string]interface{})["users"] = result
        case 1 :
            data["message"] = "User is not admin"
        case 2 :
            data["message"] = "Unexpected Error in Database"
    }
    return c.RenderJson(data)
}
