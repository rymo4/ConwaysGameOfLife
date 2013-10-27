var app = angular.module('Conways',[]);
 
app.controller('LifeCtrl', ['$scope', '$http', function($scope, $http) {
    $scope.map = Array(1);
    $scope.n_rows = 0;
    $scope.n_columns = 0;
    $scope.canonical = 0;
    var ALIVE = 1;
    var DEAD = 0;


    $scope.getCol = function(col) {
      return $scope.map[col];
    }
    $scope.getElem = function(col, row) {
      if ($scope.map[col][row] == ALIVE)
        return 'cell.gif';
      else
        return 'dead.jpg'; 
    }
    $scope.getNumber = function(num) {
      return new Array(num);   
    }
    $scope.double = function(value) { return value * 2; };
    
    $scope.toCanonical = function() {
      var canonical = $scope.n_columns + "," + $scope.n_rows + "|";
      for (var col = 0; col < $scope.n_columns; col++) {
        for (var row = 0; row < $scope.n_rows; row++) {
          if ($scope.map[col][[row]] == ALIVE) {
            canonical += col +',' + row + ',';
          }
        } 
      } 
      return canonical;
    }

    $scope.fromCanonical = function(canonical) {
      var t = canonical.split("|");
      var size = t[0];
      $scope.n_columns = size[0];
      $scope.n_rows = size[1];
      var elems = t[1];
      elems = elems.split(',');

      // INITIALIZE AN EMPTY GRIDs
      $scope.map = Array($scope.n_columns);
      for (var col = 0; col < $scope.n_columns; col++) {
        $scope.map[col] = new Array($scope.n_rows);
        for (var row = 0; row < $scope.n_rows; row++) {
            $scope.map[col][row] = DEAD;
        }
      } 

      //LOOP THROUGH THE COMPUTED ALIVE AND ADD THEM TO MAP
      for (var i=0; i<elements.length; i++) {
        $scope.map[i][i+1] = ALIVE;
      }
    }

    $scope.nextGen = function() {
      var url = "/next?state=" + $scope.canonical;
      $scope.canonical = $http.get(url)
      console.log($scope.canonical);
      $scope.fromCanonical($scope.canonical)
    }

    $scope.initMap = function(n_columns, n_rows) {
      var map = Array(n_columns);
      for (var col = 0; col < n_columns; col++) {
        map[col] = new Array(n_rows);
        for (var row = 0; row < n_rows; row++) {
          t = Math.random();
          if( t > 0.5) {
            map[col][row] = ALIVE;
          }
          else {
            map[col][row] = DEAD;
          }
        } 
      }
      $scope.n_rows = n_rows;
      $scope.n_columns = n_columns;
      $scope.map = map;
      $scope.canonical = $scope.toCanonical();
    }

}]);

