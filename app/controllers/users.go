package controllers

import (
    "github.com/revel/revel"
    "encoding/json"
    "io/ioutil"
    //"github.com/ip4368/go-password"
    "github.com/ip4368/go-userprofile"
)

type Users struct {
    *revel.Controller
}

type RegisterProfile struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Email string `json:"email"`
}

func (c *Users) Register() revel.Result {
	// read request body to byte
	var r RegisterProfile
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	json.Unmarshal([]byte(bodyBytes), &r)

	//password.ValidatePassword(r.Password)
	userprofile.ValidateEmail(r.Email)
	userprofile.ValidateUsername(r.Username)
	return c.Render()
}

func (c *Users) Login() revel.Result {
	// read request body to byte
	var r RegisterProfile
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	json.Unmarshal([]byte(bodyBytes), &r)

	userprofile.ValidateEmail(r.Email)
	return c.Render()
}
