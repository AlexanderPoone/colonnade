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
    models.GuardUsers()
}

type Users struct {
    *revel.Controller
    revmgo.MongoController
}

type RegisterProfile struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Email string `json:"email"`
}

func (c Users) Register() revel.Result {
    // read request body to byte
    var r RegisterProfile
    var bodyBytes []byte
    if c.Request.Body != nil {
        bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
    }
    json.Unmarshal([]byte(bodyBytes), &r)

    result := models.RegisterHandler(c.MongoSession, r.Email, r.Username, r.Password)

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
    }
    return c.RenderJson(data)
}

func (c Users) Login() revel.Result {
    // read request body to byte
    var r RegisterProfile
    var bodyBytes []byte
    if c.Request.Body != nil {
        bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
    }
    json.Unmarshal([]byte(bodyBytes), &r)

    result := models.LoginHandler(c.MongoSession, r.Email, r.Password)

    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = result
    switch result {
        case 0 :
            data["message"] = "Successfully Logged In"
        case 1 :
            data["message"] = "Invalid Log in"
    }
    return c.RenderJson(data)
}
