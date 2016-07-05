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
.controller("GlobCtrl", function($scope, $http){

})
.controller("mainCtrl", function($scope, $http){

})
.controller("dashboardCtrl", function($scope, $http){

})
.controller("loginCtrl", function($scope, $http){

})
.controller("404Ctrl", function($scope, $http){

});

$('#menu-button').click(function() {
	$('.ui.sidebar').sidebar('toggle');
});

$('.side-menu-button').click(function() {
	$('.ui.sidebar').sidebar('toggle');
});
