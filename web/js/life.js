var app = angular.module('Conways',[]);
 
app.controller('LifeCtrl', ['$scope', function($scope) {
    $scope.n_columns = 50;
    $scope.n_rows = 50;
    $scope.map = Array(1);

    $scope.initMap = function(n_columns, n_rows) {
      var map = Array(n_columns);
      for (var col = 0; col < n_columns; col++) {
        map[col] = new Array(n_rows);
        for (var row = 0; row < n_rows; row++) {
          map[col][row] = Math.random();
        } 
      }
      $scope.map = map;
    }
    $scope.getCol = function(col) {
      return $scope.map[col];
    }
    $scope.getElem = function(col, row) {
      if ($scope.map[col][row] > 0.5)
        return 'cell.gif';
      else
        return 'dead.jpg'; 
    }
    $scope.getNumber = function(num) {
      return new Array(num);   
    }
    $scope.double = function(value) { return value * 2; };
}]);