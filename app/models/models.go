package models

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/revel/revel"
    "github.com/ip4368/go-userprofile"
    "github.com/ip4368/go-password"
    "strings"
    "time"
)

const COODRINATORS = 0
const TUTORS = 1
const STUDENTS = 2

type User_t struct {
    Identifier [2]string `bson:"identifier"`
    Passwd string        `bson:"passwd"`
    Salt string          `bson:"salt"`
    PrevPasswd []string  `bson:"prevpasswd,omitempty"`
    Suspended bool       `bson:"suspended"`
    Name string          `bson:"name"`
    Id    bson.ObjectId  `bson:"_id,omitempty"`
}

type Course_t struct {
    Name string           `bson:"name"`
    Description string    `bson:"description"`
    Suspended bool        `bson:"suspended,omitempty"`
    Users map[string]int  `bson:"users,omitempty"`
    TimeCreated time.Time `bson:"timeCreated,omitempty"`
    Id    bson.ObjectId   `bson:"_id,omitempty"`
}

type Admin_t struct {
    UserId bson.ObjectId `bson:"uid"`
    Id     bson.ObjectId `bson:"_id,omitempty"`
}

func GuardUsers() {
    localDBSession, err := mgo.Dial("mongodb://localhost/colonnade")
    if err != nil {
        // Only warn since we'll retry later for each request
        revel.WARN.Printf("Could not connect to Mongo DB. Error: %s", err)
    } else {
        users := usersCollection(localDBSession)
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

func GuardAdmins() {
    localDBSession, err := mgo.Dial("mongodb://localhost/colonnade")
    if err != nil {
        // Only warn since we'll retry later for each request
        revel.WARN.Printf("Could not connect to Mongo DB. Error: %s", err)
    } else {
        users := adminsCollection(localDBSession)
        index := mgo.Index{
            Key: []string{"uid"},
            Unique: true,
            DropDups: true,
            Background: true,
            Sparse: true,
        }
        users.EnsureIndex(index)
        localDBSession.Close()
    }
}

func usersCollection(s *mgo.Session) *mgo.Collection {
    return s.DB("colonnade").C("users")
}

func coursesCollection(s *mgo.Session) *mgo.Collection {
    return s.DB("colonnade").C("courses")
}

func adminsCollection(s *mgo.Session) *mgo.Collection {
    return s.DB("colonnade").C("admins")
}

func RegisterHandler(s *mgo.Session, email, username, passwd, name string) int {
    // validate all email, username and password
    if !userprofile.ValidateEmail(email) { return 1 }
    if !userprofile.ValidateUsername(username) { return 2 }
    if !password.ValidatePassword(passwd) { return 3 }
    if !userprofile.ValidateName(name) { return 4 }

    hashed, salt, _ := password.HashAutoSalt(passwd)
    doc := User_t{
        Identifier: [2]string{email, username},
        Passwd: hashed,
        Salt: salt,
        PrevPasswd: []string{},
        Suspended: false,
        Name: name,
    }

    err := usersCollection(s).Insert(doc)

    if err != nil { return 5 }
    return 0
    //
}

func LoginHandler(s *mgo.Session, email, passwd string) (int, [2]string, string, string) {
    // validate password and check against database
    if !userprofile.ValidateEmail(email) { return 1, [2]string{"", ""}, "", "" }
    if !password.ValidatePassword(passwd) { return 1, [2]string{"", ""}, "", "" }

    u := new(User_t)

    usersCollection(s).Find(bson.M{
        "identifier": email,
    }).Select(bson.M{
        "prevpasswd": 0,
    }).One(u)

    if strings.Compare(u.Identifier[0], email) != 0 { return 2, [2]string{"", ""}, "", "" }
    if u.Suspended { return 3, [2]string{"", ""}, "", "" }

    hashed, _, _ := password.HashWithSalt(passwd, u.Salt)
    if strings.Compare(u.Passwd, hashed) != 0 { return 4, [2]string{"", ""}, "", "" }
    return 0, u.Identifier, u.Id.Hex(), u.Name
}

func LoginStatus(email, username, name, userId string) int {
    if(email == "" && username == "" && name == "" && userId == "") { return 1 }
    return 0;
}

func LogoutHandler(email, username, name, userId string) int {
    return LoginStatus(email, username, name, userId)
}

func CoursesForUser(s *mgo.Session, UserIdHex string) (int, []Course_t, []Course_t, []Course_t) {
    isValidId := bson.IsObjectIdHex(UserIdHex)
    if !isValidId { return 2, []Course_t{}, []Course_t{}, []Course_t{} }
    UserId := bson.ObjectIdHex(UserIdHex)
    UserIdStr := UserId.String()

    var result []Course_t
    err := coursesCollection(s).Find(bson.M{
        "$and": []bson.M{
            bson.M{UserIdStr: bson.M{"$exists": true}},
            bson.M{"suspended": true},
        },
    }).Select(bson.M{
        "name": 1,
        "description": 1,
        "users": 1,
        "_id": 1,
    }).All(&result)
    if err != nil { return 3, []Course_t{}, []Course_t{}, []Course_t{} }

    var groups [3][]Course_t
    for _, value := range result {
        groups[value.Users[UserIdStr]] = append(groups[value.Users[UserIdStr]], Course_t{
                Name: value.Name,
                Description: value.Description,
                Id: value.Id,
            })
    }

    return 0, groups[COODRINATORS], groups[TUTORS], groups[STUDENTS]
}

func CheckAdmin(s *mgo.Session, UserIdHex string) int {
    isValidId := bson.IsObjectIdHex(UserIdHex)
    if !isValidId { return 3 }
    UserId := bson.ObjectIdHex(UserIdHex)

    var result Admin_t
    err := adminsCollection(s).Find(bson.M{"uid": UserId}).One(&result)
    if err != nil { return 2 }

    return 0
}
