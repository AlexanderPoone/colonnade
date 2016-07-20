var app = angular.module("Colonnade", ['ngRoute', 'ngCookies'])

app
.config(function ($routeProvider, $locationProvider, $httpProvider){
    $httpProvider.useApplyAsync(true);
    $routeProvider
    .when('/',{
        templateUrl:'public/template/main.html',
        controller:'mainCtrl'})
    .when('/dashboard/',{
        templateUrl:'public/template/dashboard.html',
        controller:'dashboardCtrl'})
    .when('/login/',{
        templateUrl:'public/template/login.html',
        controller:'loginCtrl'})
    .when('/admin/',{
        templateUrl:'public/template/admin.html',
        controller:'adminCtrl'})
    .when('/admin/courses/',{
        templateUrl:'public/template/adminCourses.html',
        controller:'adminCoursesCtrl'})
    .when('/admin/course/new',{
        templateUrl:'public/template/adminNewCourse.html',
        controller:'adminNewCourseCtrl'})
    .when('/admin/course/:Id',{
        templateUrl:'public/template/adminCourse.html',
        controller:'adminCourseCtrl'})
    .when('/admin/users/',{
        templateUrl:'public/template/adminUsers.html',
        controller:'adminUsersCtrl'})
    .when('/admin/user/:Id',{
        templateUrl:'public/template/adminUser.html',
        controller:'adminUserCtrl'})
    .otherwise({
        templateUrl:'public/template/404.html',
        controller:'404Ctrl'});
})
.controller("GlobCtrl", function($scope, $http, $location, login){
    login.init($scope);
    $scope.logout = function(){
        login.logout(function(data){
            if(data.error == 0){
                $location.url("/");
            }
        });
    }
})
.controller("mainCtrl", function($scope, $http, login){
})
.controller("dashboardCtrl", function($scope, $http, $location, login, user){
    if(!login.loggedIn()){
        $location.url("/login/");
    }else{
        $scope.courses = {};
        user.getCoursesForUser(function(res){
            $scope.courses.coordinator = res.data.asCoordinator;
            $scope.courses.tutor       = res.data.asTutor;
            $scope.courses.student     = res.data.asStudent;
        })
    }
})
.controller("loginCtrl", function($scope, $location, REGEX, login, register){
    $scope.login = function() {
        login.launch($scope.loginInfo.email, $scope.loginInfo.password, function(data){
            if(data.error == 0){
                $location.url("/dashboard/");
                $scope.loginErrMsg = "";
            }else{
                $scope.loginErrMsg = data.message;
            }
        });
    };
    $scope.register = function() {
        var ri = $scope.registerInfo;
        var invalid = {}
        if(ri){
            invalid.email = !Boolean(ri.email ? ri.email.match(REGEX.email) : false);
            invalid.username = !Boolean(ri.username ? ri.username.match(REGEX.username) : false);
            invalid.name = !Boolean(ri.name ? ri.name.match(REGEX.name) : false);
            invalid.password = !Boolean(ri.password ? ri.password.match(REGEX.password) : false);
        }else{
            invalid.email = invalid.username = invalid.name = invalid.password = true;
        }

        if (!invalid.email && !invalid.username && !invalid.name && !invalid.password) {
            register(
                ri.email,
                ri.username,
                ri.name,
                ri.password,
                function(data){
                    if(data.error == 0){
                        $scope.registerInfo.email = "";
                        $scope.registerInfo.username = "";
                        $scope.registerInfo.name = "";
                        $scope.registerInfo.password = "";
                        $scope.regSucc = true;
                        $scope.regErrMsg = "";
                    }else{
                        $scope.regErrMsg = data.message;
                    }
                }
            );
        }
        $scope.regInvalid = invalid;
    };
})
.controller("adminCtrl", function($scope){
})
.controller("adminCoursesCtrl", function($scope, $routeParams, login, admin){
    var page = $routeParams.p ? $routeParams.p : 0 ;
    admin.getAllCourses(page, function(response){
        $scope.courses = response.data.courses;
        for(var i in $scope.courses){
            var tempDate = new Date($scope.courses[i].TimeCreated);
            $scope.courses[i].newDate = tempDate.getDate().toString() + '/' +
                                        (tempDate.getMonth() + 1).toString() + '/' +
                                        tempDate.getFullYear().toString();
        }
    });
})
.controller("adminNewCourseCtrl", function($scope, $routeParams, login, admin){
    $scope.foundUsers = [];
    $scope.changeData = function(findUserData){
        if(findUserData.query.length >= 3){
            function add2Options(response){
                $scope.foundUsers = response.data.users;
            }
            admin.findUserByIdentifier(findUserData.query, add2Options);
        }else{
            $scope.foundUsers = null;
        }
    }
    var newCourseId = null;
    $scope.step = {};
    $scope.step.style = {};
    $scope.step.style.info = {
        'active': true,
    };
    $scope.step.style.staff = {
        'active': false,
        'disabled': true,
    };
    $scope.step.style.done = {
        'active': false,
        'disabled': true,
    };
    $scope.step.current = 'info';
    $scope.step.handler = [];
    $scope.step.handler.push(function(courseInfo){
        $scope.step.style.form0 = {'loading': true };
        admin.createNewCourse(
            courseInfo.name,
            courseInfo.description,
            function(response){
                if(response.error == 0){
                    newCourseId = response.data.courseId;
                    $scope.step.style.form0 = {'loading': false };
                    $scope.step.current = 'staff';
                    $scope.step.style.info = {
                        'active': false,
                    };
                    $scope.step.style.staff = {
                        'active': true,
                        'disabled': false,
                    };
                }
            });
    });
    $scope.step.handler.push(function(users){
        var reqUsers = [];
        for(i in users){
            reqUsers.push({
                uid: users[i].Id,
                role: 0,
            });
        }

        admin.addUsers2Course(newCourseId, reqUsers, function(response){
            if(response.error == 0){
                $scope.step.current = 'done';
                $scope.step.style.staff = {
                    'active': false,
                };
                $scope.step.style.done = {
                    'active': true,
                    'disabled': false,
                };
            }
        })
    });
})
.controller("adminCourseCtrl", function($scope, $routeParams, login, admin, batchUpdate, ROLES){
    $scope.ROLES = ROLES;
    
    var courseId = $routeParams.Id;
    var originalDetail;
    $scope.editingTitleBlock = false;

    admin.getCourseDetail(courseId, function(result){
        if(result.error == 0){
            originalDetail = result.data.course;
            $scope.course = Object.assign({}, originalDetail);
            var tempDate = new Date($scope.course.TimeCreated);
            $scope.course.newDate = tempDate.getDate().toString() + '/' +
                                    (tempDate.getMonth() + 1).toString() + '/' +
                                    tempDate.getFullYear().toString();
        }
    });

    $scope.editTitleBlock = function(){
        $scope.editingTitleBlock = true;
    }

    $scope.editDescription = function(){
        $scope.editingDescription = true;
    }

    $scope.submitTitleBlock = function(){
        var details = [];
        if($scope.course.Name != originalDetail.Name){
            details.push({
                t: "Name",
                v: $scope.course.Name,
            })
        }
        if($scope.course.Suspended != originalDetail.Suspended){
            details.push({
                t: "Suspended",
                v: $scope.course.Suspended,
            })
        }

        if(details.length){  // only do work if somewhere need to be updated
            admin.updateCourse(courseId, details, function(res){
                if(res.error == 0){
                    batchUpdate($scope.course, originalDetail, details, res.data.results);
                }else{
                    $scope.course.Name = originalDetail.Name;
                    $scope.course.Suspended = originalDetail.Suspended;
                }
            });
        }
        $scope.editingTitleBlock = false;
    }

    $scope.submitDescription = function(){
        var details = [];
        if($scope.course.Description != originalDetail.Description){
            details.push({
                t: "Description",
                v: $scope.course.Description,
            })
        }
        if(details.length){  // only do work if somewhere need to be updated
            admin.updateCourse(courseId, details, function(res){
                if(res.error == 0){
                    batchUpdate($scope.course, originalDetail, details, res.data.results);
                }else{
                    $scope.course.Description = originalDetail.Description;
                }
            });
        }
        $scope.editingDescription = false;
    }

    $scope.updateCoordinators = function(){
        var userAdding = [];
        for(i in $scope.coordinatorsData.added){
            userAdding.push({uid: $scope.coordinatorsData.added[i].Id, role: ROLES.COORDINATOR});
        }
        var userRemoving = [];
        for(i in $scope.coordinatorsData.removed){
            userRemoving.push($scope.coordinatorsData.removed[i].Id);
        }

        if(userAdding.length) admin.addUsers2Course(courseId, userAdding);
        if(userRemoving.length) admin.removeUsersFromCourse(courseId, userRemoving);
    }

    $scope.updateTutors = function(){
        var userAdding = [];
        for(i in $scope.tutorsData.added){
            userAdding.push({uid: $scope.tutorsData.added[i].Id, role: ROLES.TUTOR});
        }
        var userRemoving = [];
        for(i in $scope.tutorsData.removed){
            userRemoving.push($scope.tutorsData.removed[i].Id);
        }

        if(userAdding.length) admin.addUsers2Course(courseId, userAdding);
        if(userRemoving.length) admin.removeUsersFromCourse(courseId, userRemoving);
    }

    $scope.updateStudents = function(){
        var userAdding = [];
        for(i in $scope.studentsData.added){
            userAdding.push({uid: $scope.studentsData.added[i].Id, role: ROLES.STUDENT});
        }
        var userRemoving = [];
        for(i in $scope.studentsData.removed){
            userRemoving.push($scope.studentsData.removed[i].Id);
        }

        if(userAdding.length) admin.addUsers2Course(courseId, userAdding);
        if(userRemoving.length) admin.removeUsersFromCourse(courseId, userRemoving);
    }
})
.controller("adminUsersCtrl", function($scope, $routeParams, login, admin){
    var page = $routeParams.p ? $routeParams.p : 0 ;
    admin.getAllUsers(page, function(response){
        $scope.users = response.data.users;
    });
})
.controller("adminUserCtrl", function($scope, $routeParams, login, admin, batchUpdate){
    var userId = $routeParams.Id;
    var originalDetail;
    $scope.editingNameBlock = false;

    admin.getUserDetail(userId, function(result){
        originalDetail = result.data.user;
        $scope.user = Object.assign({}, originalDetail);
    });

    $scope.editNameBlock = function(){
        $scope.editingNameBlock = true;
    }

    $scope.submitNameBlock = function(){
        var details = [];
        if($scope.user.Name != originalDetail.Name){
            details.push({
                t: "Name",
                v: $scope.user.Name,
            })
        }
        if($scope.user.Suspended != originalDetail.Suspended){
            details.push({
                t: "Suspended",
                v: $scope.user.Suspended,
            })
        }

        if(details.length){  // only do work if somewhere need to be updated
            admin.updateUser(userId, details, function(res){
                if(res.error == 0){
                    batchUpdate($scope.user, originalDetail, details, res.data.results);
                }else{
                    $scope.user.Name = originalDetail.Name;
                    $scope.user.Suspended = originalDetail.Suspended;
                }
            });
        }
        $scope.editingNameBlock = false;
    }
})
.controller("404Ctrl", function($scope, $http, login){
})
.factory('login', function($http, $location, API){
    var user = {};
    user.loggedIn = false;
    user.name = "";
    user.email = "";
    user.admin = false;
    var globScope = null;
    var checkLogin = function(callback){
        $http.get(API.url + '/user/loginInfo', {
            withCredentials: true,
        }).then(function successCallback(response) {
            if(response.data.error == 0){
                user.loggedIn = true;
                user.name = response.data.data.name;
                user.email = response.data.data.email;
                user.admin = globScope.admin = response.data.data.admin;
                globScope.login = true;
            }else{
                user.loggedIn = false;
                user.name = "";
                user.email = "";
                user.admin = globScope.admin = false;
                globScope.login = false;
            }
            if(callback) callback(response.data);
        }, function errorCallback(response) {
            console.log("error");
            if(callback) callback(response.data);
        });
    }
    return {
        init: function(scope, callback){
            globScope = scope;
            checkLogin();
        },
        launch: function(email, password, callback){
            $http.post(API.url + '/user/login', {
                email: email,
                password: password
            }, {
                withCredentials: true,
            }).then(function successCallback(response) {
                if(response.data.error == 0){
                    user.loggedIn = true;
                    user.name = response.data.data.name;
                    user.email = email;
                    user.admin = globScope.admin = response.data.data.admin;
                    globScope.login = true;
                }
                if(callback) callback(response.data);
            }, function errorCallback(response) {
                console.log("error");
                if(callback) callback(response.data);
            });
        },
        getUser: function(){
            return user;
        },
        checkLogin: checkLogin,
        logout: function(callback){
            $http.get(API.url + '/user/logout', {
                withCredentials: true,
            }).then(function successCallback(response) {
                if(response.data.error == 0){
                    user.loggedIn = false;
                    user.name = "";
                    user.email = "";
                    user.admin = globScope.admin = false;
                    globScope.login = false;
                }
                if(callback) callback(response.data);
            }, function errorCallback(response) {
                console.log("error");
                if(callback) callback(response.data);
            });
        },
        loggedIn: function(){
            return user.loggedIn;
        },
    }
})
.factory('register', function($http, API){
    return function(email, username, name, password, callback){
        $http.post(API.url + '/user/register', {
            email: email,
            password: password,
            name: name,
            username: username
        }, {
            withCredentials: true,
        }).then(function successCallback(response) {
            if(callback) callback(response.data);
        }, function errorCallback(response) {
            console.log("error");
            if(callback) callback(response.data);
        });
    }
})
.factory('user', function($http, API){
    return {
        getCoursesForUser: function(callback){
            $http.get(API.url + "/courses", {
                withCredentials: true,
            }).then(function successCallback(response){
                callback(response.data);
            }, function errorCallback(response){
                console.log("error");
                if(callback) callback(response.data);
            })
        },
    }
})
.factory('admin', function($http, API, login){
    return {
        getAllCourses: function(p, callback){
            $http.get(API.url + "/admin/courses", {
                params: {p: p},
                withCredentials: true,
            }).then(function successCallback(response){
                if(response.data.error == 0) if(callback) callback(response.data);
            },function errorCallback(response){
                console.log("error");
                if(callback) callback(response.data);
            });
        },
        getCourseDetail: function(courseId, callback){
            $http.get(API.url + "/admin/course/" + courseId, {
                withCredentials: true,
            }).then(function successCallback(response){
                if(response.data.data.course.Users === null){
                    response.data.data.course.Users = [];
                }
                if(callback) callback(response.data);
            }, function errorCallback(response){
                console.log("error");
                if(callback) callback(response.data);
            });
        },
        createNewCourse: function(name, description, callback){
            $http.post(API.url + "/admin/course/new", {
                name: name,
                description: description,
            }, {
                withCredentials: true,
            }).then(function successCallback(response){
                if(callback) callback(response.data);
            }, function errorCallback(response){
                console.log("error");
                if(callback) callback(response.data);
            })
        },
        getAllUsers: function(p, callback){
            $http.get(API.url + "/admin/users", {
                params: {p: p},
                withCredentials: true,
            }).then(function successCallback(response){
                if(callback) callback(response.data);
            }, function errorCallback(response){
                console.log("error");
                if(callback) callback(response.data);
            });
        },
        getUserDetail: function(userId, callback){
            $http.get(API.url + "/admin/user/" + userId, {
                withCredentials: true,
            }).then(function successCallback(response){
                if(callback) callback(response.data);
            }, function errorCallback(response){
                console.log("error");
                if(callback) callback(response.data);
            });
        },
        addUsers2Course: function(courseId, users, callback){
            $http.post(API.url + "/admin/course/" + courseId + "/addUsers",{
                users: users,
            }, {
                withCredentials: true,
            }).then(function successCallback(response){
                if(callback) callback(response.data);
            }, function errorCallback(response){
                console.log("error");
                if(callback) callback(response.data);
            })
        },
        removeUsersFromCourse: function(courseId, users, callback){
            $http.post(API.url + "/admin/course/" + courseId + "/removeUsers",{
                users: users,
            }, {
                withCredentials: true,
            }).then(function successCallback(response){
                if(callback) callback(response.data);
            }, function errorCallback(response){
                console.log("error");
                if(callback) callback(response.data);
            })
        },
        findUserByIdentifier: function(identifier, callback){
            $http.get(API.url + "/admin/findUser", {
                params: {q: identifier},
                withCredentials: true,
            }).then(function successCallback(response){
                if(response.data.error == 0) if(callback) callback(response.data);
            }, function errorCallback(response){
                console.log("error");
                if(callback) callback(response.data);
            });
        },
        updateCourse: function(courseId, details, callback){
            $http.post(API.url + "/admin/course/" + courseId + "/update", {
                d: details,
            },{
                withCredentials: true,
            }).then(function successCallback(response){
                if(callback) callback(response.data);
            }, function errorCallback(response){
                console.log("error");
                if(callback) callback(response.data);
            })
        },
        updateUser: function(userId, details, callback){
            $http.post(API.url + "/admin/user/" + userId + "/update", {
                d: details,
            },{
                withCredentials: true,
            }).then(function successCallback(response){
                if(callback) callback(response.data);
            }, function errorCallback(response){
                console.log("error");
                if(callback) callback(response.data);
            })
        },
    }
})
.factory("batchUpdate", function(){
    return function(scopeData, origData, submittedData, isSuccessArray){
        for(i in isSuccessArray){
            var Type = submittedData[i].t;
            if(isSuccessArray[i] == 0){
                origData[Type] = scopeData[Type]
            }else{
                scopeData[Type] = origData[Type]
            }
        }
    }
})
.directive("findUser", function(){
    function link(scope, elem, attrs, ngModelCtrl){
        var cooe = 40;
        var inputCooe = 8;
        scope.inputWidth = {
            width: "18px",
        }
        if(!scope.ngModel){
            scope.ngModel = {
                chosen : [],
                query  : "",
            };
        }
        function calcHeight(list, cooe){
            try{
                if(list.length<=1) return (1 * cooe + 1).toString() + "px";
                if(list.length>=5) return (5 * cooe + 1).toString() + "px";
                return (list.length * cooe + 1).toString() + "px";
            } catch (err) {
                return (cooe + 1).toString() + "px";;
            }
        }
        scope.dropdownMain = {};
        scope.dropdownContainer = {};
        var open = false;
        function openSearching(){
            open = true;
            $('.dropdown-search-input').focus();
            scope.dropdownContainer = {
                height: calcHeight(scope.options, cooe),
                "border-width": "1px",
            };
            scope.dropdownMain = {
                "border-radius": "4px 4px 0 0",
            };
        }
        function closeSearching(){
            open = false;
            scope.dropdownContainer = {
                height: "0px",
                "border-width": "0px",
            };
            scope.dropdownMain = {
                "border-radius": "4px",
            };
        }
        function updateQueryBox(query){
            scope.inputWidth = {
                width: (query.length * inputCooe + 12).toString() + "px",
            }
        }
        scope.openSearching = function(){openSearching();}
        scope.toggleSearching = function(){
            if(!open) openSearching();
            else closeSearching();
        }
        scope.add = function(user){
            new_user = {};
            new_user.Id = user.Id;
            new_user.Identifier = user.Identifier;
            new_user.Name = user.Name;
            new_user.Suspended = user.Suspended;
            scope.ngModel.chosen.push(new_user);
            scope.ngModel.query = "";
            updateQueryBox(scope.ngModel.query);
            closeSearching();
        }
        scope.remove = function(user){
            var i = scope.ngModel.chosen.indexOf(user);
            if(i > -1) scope.ngModel.chosen.splice(i, 1);
        }
        scope.queryChange = function(query){
            updateQueryBox(query);
        }
        scope.$watch('ngModel.chosen+ngModel.query', function() {
            ngModelCtrl.$setViewValue({chosen: scope.ngModel.chosen, query: scope.ngModel.query});
        });
        scope.$watch('options', function(){
            if(open){
                scope.dropdownMain = {
                    height: calcHeight(scope.options, cooe),
                };
            }
        }, true);
    }
    return {
        restrict: "E",
        require: "ngModel",
        templateUrl: 'public/template/findUser.html',
        scope: {
            options: "=",
            ngModel: "=",
        },
        link: link,
    }
})
.directive('listEditUser', function(admin){
    function link(scope, elem, attrs, ngModelCtrl){
        var pending = {
            removed: [],
            added  : [],
        };

        scope.edit = false;
        scope.findUserData = {chosen: [], query: ""};
        scope.editMode = function(){
            scope.edit=true;
        }
        scope.removeUser = function(user){
            var tempUser = {
                Id         : user.Detail.Id,
                Identifier : user.Detail.Identifier,
                Name       : user.Detail.Name,
                Suspended  : user.Detail.Suspended,
            };

            pending.removed.push(tempUser);
            ngModelCtrl.$setViewValue(pending);
            scope.users.pop(user);
        }
        scope.addUser = function(users){
            for(i in users){
                var tempUser = {
                    Id         : users[i].Id,
                    Identifier : users[i].Identifier,
                    Name       : users[i].Name,
                    Suspended  : users[i].Suspended,
                };

                pending.added.push(tempUser);
                scope.users.push({
                    Role   : scope.role,
                    Detail : users[i],
                });
            }
            scope.findUserData.chosen = [];

            ngModelCtrl.$setViewValue(pending);
        }
        scope.submitChange = function(){
            scope.update();
            scope.findUserData.chosen = [];
            scope.findUserData.query  = "";
            viewMode();
        }
        scope.changeData = function(findUserData){
            if(findUserData.query.length >= 3){
                function add2Options(response){
                    scope.foundUsers = response.data.users;
                }
                admin.findUserByIdentifier(findUserData.query, add2Options);
            }else{
                scope.foundUsers = null;
            }
        }

        function viewMode(){
            scope.edit=false;
        }
    }

    return {
        restrict: "E",
        require: "ngModel",
        templateUrl: "public/template/listEditUser.html",
        scope: {
            users: "=",
            role: "=",
            title: "@",
            update: "&",
            ngModel: "=",
        },
        link: link,
    }
})
.constant('ROLES', {COORDINATOR: 0, TUTOR: 1, STUDENT: 2})
.constant('API', {url: './api'})
.constant('REGEX', {
    email    : /^[a-z0-9._%+-]+@(?:[a-z0-9-]+\.)+[a-z]{2,4}$/,
    password : /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,20}$/,
    username : /^[a-zA-Z0-9]{2,12}$/,
    name     : /^.{3,}$/,
});

$('#menu-button').click(function() {
    $('.ui.sidebar').sidebar('toggle');
});

$('.side-menu-button').click(function() {
    $('.ui.sidebar').sidebar('toggle');
});
