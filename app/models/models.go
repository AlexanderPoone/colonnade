package models

import (
    "encoding/json"
    "io"
    "io/ioutil"
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

type User_db struct {
    Identifier  [2]string      `bson:"identifier"`
    Passwd      string         `bson:"passwd"`
    Salt        string         `bson:"salt"`
    PrevPasswd  []string       `bson:"prevpasswd,omitempty"`
    Suspended   bool           `bson:"suspended"`
    Name        string         `bson:"name"`
    Id          bson.ObjectId  `bson:"_id,omitempty"`
}

type User_t struct {
    Email      string
    Username   string
    Name       string
    UserIdHex  string
}

type UserInCourse struct {
    Id    bson.ObjectId  `bson:"uid"`
    Name  string         `bson:"name"`
    Role  int            `bson:"role"`
}

type Course_db struct {
    Name string                 `bson:"name"`
    Description string          `bson:"description"`
    Suspended   bool            `bson:"suspended,omitempty"`
    Users       []UserInCourse  `bson:"users,omitempty"`
    TimeCreated time.Time       `bson:"timeCreated,omitempty"`
    Id          bson.ObjectId   `bson:"_id,omitempty"`
}

type Course_t struct {
    Name         string  `json:"name"`
    Description  string  `json:"description"`
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

func ParseBody(body io.Reader, r interface{}) {
    var bodyBytes []byte
    if body != nil {
        bodyBytes, _ = ioutil.ReadAll(body)
    }
    json.Unmarshal([]byte(bodyBytes), &r)
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
    doc := User_db{
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

    u := new(User_db)

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

func LoginStatus(user User_t) int {
    if(user.Email == "" && 
       user.Username == "" &&
       user.Name == "" &&
       user.UserIdHex == "") { return 1 }
    return 0;
}

func LogoutHandler(user User_t) int {
    return LoginStatus(user)
}

func CoursesForUser(s *mgo.Session, UserIdHex string) (int, []Course_db, []Course_db, []Course_db) {
    isValidId := bson.IsObjectIdHex(UserIdHex)
    if !isValidId { return 2, []Course_db{}, []Course_db{}, []Course_db{} }
    UserId := bson.ObjectIdHex(UserIdHex)

    var result []Course_db
    err := coursesCollection(s).Find(bson.M{
        "$and": []bson.M{
            bson.M{"users.uid": UserId},
            bson.M{"suspended": false},
        },
    }).Select(bson.M{
        "name": 1,
        "description": 1,
        "users": bson.M{"$elemMatch":bson.M{"$eq":UserId}},
        "_id": 1,
    }).All(&result)
    if err != nil { return 3, []Course_db{}, []Course_db{}, []Course_db{} }

    var groups [3][]Course_db
    for _, value := range result {
        groups[value.Users[0].Role] = append(groups[value.Users[0].Role], Course_db{
                Name: value.Name,
                Description: value.Description,
                Id: value.Id,
            })
    }

    return 0, groups[COODRINATORS], groups[TUTORS], groups[STUDENTS]
}

func CheckAdmin(s *mgo.Session, user User_t) int {
    if LoginStatus(user) != 0 { return 1 }
    isValidId := bson.IsObjectIdHex(user.UserIdHex)
    if !isValidId { return 3 }
    UserId := bson.ObjectIdHex(user.UserIdHex)

    var result Admin_t
    err := adminsCollection(s).Find(bson.M{"uid": UserId}).One(&result)
    if err != nil { return 2 }

    return 0
}

func IsAdmin(user User_t, admin string) int {
    if LoginStatus(user) != 0 { return 1 }
    if admin != "t" {return 2}
    return 0
}

func AdminCourses(s *mgo.Session, user User_t, admin string, page int) (int, []Course_db) {
    if IsAdmin(user, admin) != 0 { return 1, []Course_db{} }

    var result []Course_db
    var limit = 20
    var skip = page * 20
    err := coursesCollection(s).Find(bson.M{}).Select(bson.M{
        "description": 1,
        "name": 1,
        "timeCreated": 1,
        "_id": 1,
    }).Skip(skip).Limit(limit).All(&result)

    if err != nil { return 3, []Course_db{} }
    return 0, result
}

func AdminNewCourse(s *mgo.Session, user User_t, admin string, course Course_t) (int, string) {
    if IsAdmin(user, admin) != 0 { return 1, "" }

    id := bson.NewObjectId()
    newCourse := Course_db{
        Name        : course.Name,
        Description : course.Description,
        Suspended   : false,
        Users       : []UserInCourse{},
        TimeCreated : time.Now(),
        Id          : id,
    }

    err := coursesCollection(s).Insert(newCourse)
    if err != nil { return 2, "" }

    return 0, id.Hex()
}
