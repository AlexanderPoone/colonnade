var app = angular.module("Colonnade", ['ngRoute', 'ngCookies'])
var API_URL = './api';

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
.controller("GlobCtrl", function($scope, $rootScope, $http){

})
.controller("mainCtrl", function($scope, $rootScope, $http){

})
.controller("dashboardCtrl", function($scope, $rootScope, $http){

})
.controller("loginCtrl", function($scope, $rootScope, $location, login, register){
	$scope.login = function() {
		login($scope.loginInfo.email, $scope.loginInfo.password, function(data){
			if(data.error == 0){
				$location.url("/");
			}
		});
	};
	$scope.register = function() {
		register(
			$scope.registerInfo.email,
			$scope.registerInfo.username,
			$scope.registerInfo.name,
			$scope.registerInfo.password,
			function(data){
				if(data.error == 0){
					$scope.registerInfo.email = ""
					$scope.registerInfo.username = ""
					$scope.registerInfo.name = ""
					$scope.registerInfo.password = ""
				}
			}
		);
	};
})
.controller("404Ctrl", function($scope, $rootScope, $http){

})
.factory('login', function($http){
	return function(email, password, callback){
		$http.post(API_URL + '/user/login', {
			email: email,
			password: password
		}, {
			withCredentials: true,
		}).then(function successCallback(response) {
			callback(response.data);
		}, function errorCallback(response) {
			console.log("error");
			callback(response.data);
		});
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
			callback(response.data);
		}, function errorCallback(response) {
			console.log("error");
			callback(response.data);
		});
	}
});

$('#menu-button').click(function() {
	$('.ui.sidebar').sidebar('toggle');
});

$('.side-menu-button').click(function() {
	$('.ui.sidebar').sidebar('toggle');
});
