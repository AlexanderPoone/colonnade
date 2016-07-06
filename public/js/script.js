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
.controller("404Ctrl", function($scope, $http, login){

})
.factory('login', function($http){
	var user = {};
	user.loggedIn = false;
	user.name = "";
	user.email = "";
	globScope = null;
	return {
		init: function(scope){
			globScope = scope;
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
		checkLogin: function(){
			return null;
		},
		logout: function(callback){
			$http.get(API_URL + '/user/logout', {
				withCredentials: true,
			}).then(function successCallback(response) {
				if(response.data.error == 0){
					user.loggedIn = false;
					user.name = "";
					user.email = "";
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
});

$('#menu-button').click(function() {
	$('.ui.sidebar').sidebar('toggle');
});

$('.side-menu-button').click(function() {
	$('.ui.sidebar').sidebar('toggle');
});
