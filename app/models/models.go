package models

import (
    "gopkg.in/mgo.v2"
    //"gopkg.in/mgo.v2/bson"
    "github.com/revel/revel"
    "github.com/ip4368/go-userprofile"
    "github.com/ip4368/go-password"
)

type Users_t struct {
    Identifier [2]string
    Passwd string
    Salt string
    PrevPasswd []string
    Suspended bool
}

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

    hashed, salt, _ := password.HashAutoSalt(passwd)
    doc := Users_t{
        [2]string{email, username},
        hashed,
        salt,
        []string{},
        false,
    }

    // get collection from mgo session
    users := s.DB("colonnade").C("users")
    err := users.Insert(doc)

    if err != nil { return 4 }
    return 0
    //
}

func LoginHandler(s *mgo.Session, email, passwd string) int {
    // validate password and check against database
    userprofile.ValidateEmail(email)
    password.ValidatePassword(passwd)

    return 0
}