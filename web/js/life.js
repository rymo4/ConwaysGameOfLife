var app = angular.module('Conways',[]);
 
app.controller('LifeCtrl', ['$scope', function($scope) {
    $scope.double = function(value) { return value * 2; };
}]);