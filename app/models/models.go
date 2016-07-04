package models

import (
    "gopkg.in/mgo.v2"
    //"gopkg.in/mgo.v2/bson"
    "github.com/revel/revel"
    "github.com/ip4368/go-userprofile"
    "github.com/ip4368/go-password"
)

func GuardUsers() {
    localDBSession, err := mgo.Dial("mongodb://localhost/colonnade")
    if err != nil {
        // Only warn since we'll retry later for each request
        revel.WARN.Printf("Could not connect to Mongo DB. Error: %s", err)
    } else {
        colonnadeDB := localDBSession.DB("colonnade")
        users := colonnadeDB.C("users")
        index := mgo.Index{
            Key: []string{"identifier"},
            Unique: true,
            DropDups: true,
            Background: true,
            Sparse: true,
        }
        users.EnsureIndex(index)
        localDBSession.Close()
    }
}

func RegisterHandler(s *mgo.Session, email, username, passwd string) int {
    // validate all email, username and password
    if !userprofile.ValidateEmail(email) { return 1 }
    if !userprofile.ValidateUsername(username) { return 2 }
    if !password.ValidatePassword(passwd) { return 3 }
    return 0
    //hashed, salt, _ := password.HashAutoSalt(r.Password)
}

func LoginHandler(s *mgo.Session, email, passwd string) int {
    // validate password and check against database
    userprofile.ValidateEmail(email)
    password.ValidatePassword(passwd)

    return 0
}