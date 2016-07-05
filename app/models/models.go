package models

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/revel/revel"
    "github.com/ip4368/go-userprofile"
    "github.com/ip4368/go-password"
    "strings"
)

type Users_t struct {
    Identifier [2]string `bson:"identifier"`
    Passwd string `bson:"passwd"`
    Salt string `bson:"salt"`
    PrevPasswd []string `bson:"prevpasswd"`
    Suspended bool `bson:"suspended"`
    Id    bson.ObjectId `bson:"_id,omitempty"`
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

func userCollection(s *mgo.Session) *mgo.Collection {
    return s.DB("colonnade").C("users")
}

func RegisterHandler(s *mgo.Session, email, username, passwd string) int {
    // validate all email, username and password
    if !userprofile.ValidateEmail(email) { return 1 }
    if !userprofile.ValidateUsername(username) { return 2 }
    if !password.ValidatePassword(passwd) { return 3 }

    hashed, salt, _ := password.HashAutoSalt(passwd)
    doc := Users_t{
        Identifier: [2]string{email, username},
        Passwd: hashed,
        Salt: salt,
        PrevPasswd: []string{},
        Suspended: false,
    }

    err := userCollection(s).Insert(doc)

    if err != nil { return 4 }
    return 0
    //
}

func LoginHandler(s *mgo.Session, email, passwd string) (int, [2]string, string) {
    // validate password and check against database
    if !userprofile.ValidateEmail(email) { return 1, [2]string{"", ""}, "" }
    if !password.ValidatePassword(passwd) { return 1, [2]string{"", ""}, "" }

    u := new(Users_t)

    userCollection(s).Find(bson.M{"identifier": email}).One(u)

    if strings.Compare(u.Identifier[0], email) != 0 { return 2, [2]string{"", ""}, "" }
    if u.Suspended { return 3, [2]string{"", ""}, "" }

    hashed, _, _ := password.HashWithSalt(passwd, u.Salt)
    if strings.Compare(u.Passwd, hashed) != 0 { return 4, [2]string{"", ""}, "" }
    return 0, u.Identifier, u.Id.String()
}
