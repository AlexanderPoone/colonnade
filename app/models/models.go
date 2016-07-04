package models

import (
    "gopkg.in/mgo.v2"
    //"gopkg.in/mgo.v2/bson"
    "github.com/revel/revel"
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