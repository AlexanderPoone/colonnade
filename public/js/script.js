var app = angular.module("Colonnade", ['ngRoute', 'ngCookies'])
var API_URL = './api';

var email_regex = /^[a-z0-9._%+-]+@(?:[a-z0-9-]+\.)+[a-z]{2,4}$/;
var password_regex = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,20}$/;
var username_regex = /^[a-zA-Z0-9]{2,12}$/;
var name_regex = /^.{3,}$/;

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
    .when('/admin/:Type?/:Id?',{
        templateUrl:'public/template/admin.html',
        controller:'adminCtrl'})
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
.controller("dashboardCtrl", function($scope, $http, login){
})
.controller("loginCtrl", function($scope, $location, login, register){
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
            invalid.email = !Boolean(ri.email ? ri.email.match(email_regex) : false);
            invalid.username = !Boolean(ri.username ? ri.username.match(username_regex) : false);
            invalid.name = !Boolean(ri.name ? ri.name.match(name_regex) : false);
            invalid.password = !Boolean(ri.password ? ri.password.match(password_regex) : false);
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
.controller("adminCtrl", function($scope, $routeParams, login, admin){
    if($routeParams.Type===undefined){
        $scope.page = "main";
    }else if($routeParams.Type=="courses"){
        $scope.page = "listCourses";
        var page = $routeParams.p ? $routeParams.p : 0 ;
        admin.getAllCourses(page, function(response){
            $scope.courses = response.data.courses;
            for(var i in $scope.courses){
                var tempDate = new Date($scope.courses[i].TimeCreated);
                $scope.courses[i].newDate = tempDate.toLocaleDateString();
            }
        })
    }else if($routeParams.Type="course"){
        if($routeParams.Id === "new"){
            $scope.page = "newCourse";
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
            $scope.step.handler.push(function(){
                $scope.step.style.form0 = {'loading': true };
                admin.createNewCourse(
                    $scope.step.data.name,
                    $scope.step.data.description,
                    function(response){
                        if(response.error == 0){
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
        }
    }
})
.controller("404Ctrl", function($scope, $http, login){
})
.factory('login', function($http){
    var user = {};
    user.loggedIn = false;
    user.name = "";
    user.email = "";
    user.admin = false;
    var globScope = null;
    var checkLogin = function(callback){
        $http.get(API_URL + '/user/loginInfo', {
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
            $http.post(API_URL + '/user/login', {
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
            $http.get(API_URL + '/user/logout', {
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
        }
    }
})
.factory('register', function($http){
    return function(email, username, name, password, callback){
        $http.post(API_URL + '/user/register', {
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
.factory('admin', function($http, login){
    return {
        getAllCourses: function(p, callback){
            $http.get(API_URL + "/admin/courses", {
                params: {p: p},
                withCredentials: true,
            }).then(function successCallback(response){
                if(response.data.error == 0) if(callback) callback(response.data);
            },function errorCallback(response){
                console.log("error");
                if(callback) callback(response.data);
            });
        },
        createNewCourse: function(name, description, callback){
            $http.post(API_URL + "/admin/course/new", {
                name: name,
                description: description,
            }, {
                withCredentials: true,
            }).then(function successCallback(response){
                callback(response.data);
            }, function errorCallback(response){
                console.log("error");
                callback(response.data);
            })
        },
        findUserByIdentifier: function(identifier, callback){
            $http.get(API_URL + "/admin/findUser", {
                params: {q: identifier},
                withCredentials: true,
            }).then(function successCallback(response){
                if(response.data.error == 0) if(callback) callback(response.data);
            }, function errorCallback(response){
                console.log("error");
                if(callback) callback(response.data);
            });
        }
    }
})
.directive("findUser", function(){
    function link(scope, elem, attrs, ngModelCtrl){
        var cooe = 40;
        var inputCooe = 8;
        scope.inputWidth = {
            width: "8px",
        }
        scope.chosen = [];
        function calcHeight(list, cooe){
            try{
                if(list.length<=1) return (2 * cooe).toString() + "px";
                if(list.length>=6) return (6 * cooe).toString() + "px";
                return (list.length * cooe).toString() + "px";
            } catch (err) {
                return (2 * cooe).toString() + "px";
            }
        }
        scope.dropdownMain = {};
        var open = false;
        function openSearching(){
            open = true;
            $('.dropdown-search-input').focus();
            scope.dropdownMain = {
                height: calcHeight(scope.options, cooe),
            };
        }
        function closeSearching(){
            open = false;
            scope.dropdownMain = {
                height: cooe.toString() + "px",
            };
        }
        scope.query = "";
        scope.openSearching = function(){openSearching();}
        scope.toggleSearching = function(){
            if(!open) openSearching();
            else closeSearching();
        }
        scope.add = function(user){
            new_user = {};
            new_user.Id = user.Id;
            new_user.Email = user.Email;
            new_user.Name = user.Name;
            scope.chosen.push(new_user);
            console.log(scope.chosen);
            closeSearching();
        }
        scope.remove = function(user){
            var i = scope.chosen.indexOf(user);
            if(i > -1) scope.chosen.splice(i, 1);
        }
        scope.queryChange = function(query){
            scope.inputWidth = {
                width: ((query.length + 1) * inputCooe).toString() + "px",
            }
            scope.query = query;
        }
        scope.$watch('chosen+query', function() {  
            ngModelCtrl.$setViewValue({chose: scope.chosen, query: scope.query});
        });
    }
    return {
        restrict: "E",
        require: "ngModel",
        templateUrl: 'public/template/findUser.html',
        scope: {
            options: "=",
        },
        link: link,
    }
});

$('#menu-button').click(function() {
    $('.ui.sidebar').sidebar('toggle');
});

$('.side-menu-button').click(function() {
    $('.ui.sidebar').sidebar('toggle');
});
