package controllers

import (
    "github.com/revel/revel"
    "encoding/json"
    "io/ioutil"
    "github.com/ip4368/go-password"
    "github.com/ip4368/go-userprofile"
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
    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = nil
    data["message"] = "Successfully Registered"

    // read request body to byte
    var r RegisterProfile
    var bodyBytes []byte
    if c.Request.Body != nil {
        bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
    }
    json.Unmarshal([]byte(bodyBytes), &r)

    // validate all email, username and password
    if !userprofile.ValidateEmail(r.Email) {
        data["error"] = 1
        data["message"] = "Invalid Email"
        return c.RenderJson(data)
    }
    if !userprofile.ValidateUsername(r.Username) {
        data["error"] = 2
        data["message"] = "Invalid Username"
        return c.RenderJson(data)
    }
    if !password.ValidatePassword(r.Password) {
        data["error"] = 3
        data["message"] = "Invalid Password"
        return c.RenderJson(data)
    }
    
    //hashed, salt, _ := password.HashAutoSalt(r.Password)
    return c.RenderJson(data)
}

func (c Users) Login() revel.Result {
    // start with initialise response interface
    data := make(map[string]interface{})
    data["error"] = nil
    data["message"] = "Successfully Logged In"

    // read request body to byte
    var r RegisterProfile
    var bodyBytes []byte
    if c.Request.Body != nil {
        bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
    }
    json.Unmarshal([]byte(bodyBytes), &r)

    userprofile.ValidateEmail(r.Email)
    return c.RenderJson(data)
}
