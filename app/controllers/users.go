package controllers

import (
    "github.com/revel/revel"
    "github.com/janekolszak/revmgo"
    "github.com/ip4368/colonnade/app/models"
)

func init() {
    revmgo.ControllerInit()
    models.GuardUsers()
}

type Users struct {
    *revel.Controller
    revmgo.MongoController
}

type RegisterProfile struct {
    Username string  `json:"username"`
    Password string  `json:"password"`
    Email    string  `json:"email"`
    Name     string  `json:"name"`
}

func (c Users) Register() revel.Result {
    // read request body to byte
    var r RegisterProfile
    models.ParseBody(c.Request.Body, &r)

    result := models.RegisterHandler(c.MongoSession, r.Email, r.Username, r.Password, r.Name)

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "Successfully Registered"
        case 1 :
            data["message"] = "Invalid Email"
        case 2 :
            data["message"] = "Invalid Username"
        case 3 :
            data["message"] = "Invalid Password"
        case 4 :
            data["message"] = "Invalid Name"
        case 5 :
            data["message"] = "Username/Password has been used"
    }
    return c.RenderJson(data)
}

func (c Users) Login() revel.Result {
    // read request body to byte
    var r RegisterProfile
    models.ParseBody(c.Request.Body, &r)

    result, identifier, id, name := models.LoginHandler(c.MongoSession, r.Email, r.Password)
    admin := models.CheckAdmin(
        c.MongoSession,
        models.User_t{
            Email: identifier[0],
            Username: identifier[1],
            Name: name,
            UserIdHex: id,
        })

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "Successfully Logged In"
            data["data"] = make(map[string]interface{})
            data["data"].(map[string]interface{})["name"] = name
            c.Session["email"] = identifier[0]
            c.Session["username"] = identifier[1]
            c.Session["name"] = name
            c.Session["userId"] = id
            if admin == 0 {
                c.Session["admin"] = "t"
                data["data"].(map[string]interface{})["admin"] = true
            }
        case 1 :
            data["message"] = "Invalid Login Details"
        case 2 :
            data["message"] = "Email is not registered"
        case 3 :
            data["message"] = "User has been suspended"
        case 4 :
            data["message"] = "Password incorrect"
    }
    return c.RenderJson(data)
}

func (c Users) Logout() revel.Result {
    result := models.LogoutHandler(
        models.User_t{
            Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
        })

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "Successfully Logged Out"
            if c.Session["email"] != "" {c.Session["email"] = "" }
            if c.Session["username"] != "" { c.Session["username"] = "" }
            if c.Session["name"] != "" { c.Session["name"] = "" }
            if c.Session["userId"] != "" { c.Session["userId"] = "" }
            if c.Session["admin"] != "" { c.Session["admin"] = "" }
        case 1 :
            data["message"] = "Not Logged In"
    }
    return c.RenderJson(data)
}

func (c Users) LoginInfo() revel.Result {
    result := models.LoginStatus(
        models.User_t{
            Email: c.Session["email"],
            Username: c.Session["username"],
            Name: c.Session["name"],
            UserIdHex: c.Session["userId"],
        })

    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "Logged In"
            data["data"] = make(map[string]interface{})
            data["data"].(map[string]interface{})["name"] = c.Session["name"]
            data["data"].(map[string]interface{})["email"] = c.Session["email"]
            if c.Session["admin"] == "t" {
                data["data"].(map[string]interface{})["admin"] = true
            }
        case 1 :
            data["message"] = "Not Logged In"
    }
    return c.RenderJson(data)
}
