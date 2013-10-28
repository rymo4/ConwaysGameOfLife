var app = angular.module('Conways',[]);
 
app.controller('LifeCtrl', ['$scope', '$http', '$q', function($scope, $http, $q) {
    $scope.map = Array(1);
    $scope.map_size = Array(1);
    $scope.n_rows = 0;
    $scope.n_columns = 0;
    $scope.canonical = 0;
    $scope.board = "";
    $scope.alive = [];
    $scope.prev = [];
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
      var size = t[0].split(",");
      var n_columns = parseInt(size[0]);
      var n_rows = parseInt(size[1]);
      var elems = t[1];
      $scope.alive = elems.split(',');

      // KILL OF ALIVE CELLS FROM PREV
      //var map = Array($scope.n_columns);
      // for (var col = 0; col < $scope.n_columns; col++) {
      //   map[col] = new Array($scope.n_rows);
      //   for (var row = 0; row < $scope.n_rows; row++) {
      //       map[col][row] = DEAD;
      //       var elem = document.getElementById(col + '-' + row); 
      //       elem.setAttribute('class', 'cell dead');
      //   }
      // } 
      var map = document.getElementById('map');
      console.log(map);
      var insertFunction = removeToInsertLater(map);
      for (var i=0; i< $scope.prev.length-1; i=i+2) {
        var col = parseInt($scope.prev[i]);
        var row = parseInt($scope.prev[i+1]);
        //$scope.map[col][row] = DEAD;
        var elem = map.querySelector('li[id="' + col + '-' + row + '"]'); 
        elem.setAttribute('class', 'cell dead');
      }
      //LOOP THROUGH THE COMPUTED ALIVE AND ADD THEM TO MAP
      for (var i=0; i<$scope.alive.length-1; i=i+2) {
        var col = parseInt($scope.alive[i]);
        var row = parseInt($scope.alive[i+1]);
        //$scope.map[col][row] = ALIVE;
        console.log('li[id="' + col + '-' + row + '"]');
        var elem = map.querySelector('li[id="' + col + '-' + row + '"]');
        elem.setAttribute('class', 'cell alive');
      }
      console.log(map);
      insertFunction();
      $scope.prev = $scope.alive;
      //$scope.n_columns = n_columns;
      //$scope.n_rows = n_rows;
      //$scope.map = map;
    }

  
    /**
     * Remove an element and provide a function that inserts it into its original position
     * @param element {Element} The element to be temporarily removed
     * @return {Function} A function that inserts the element into its original position
     **/
    function removeToInsertLater(element) {
      var parentNode = element.parentNode;
      var nextSibling = element.nextSibling;
      parentNode.removeChild(element);
      return function() {
        if (nextSibling) {
          parentNode.insertBefore(element, nextSibling);
        } else {
          parentNode.appendChild(element);
        }
      };
    }

    $scope.loop = function() {
      //Reset the trail back to zero
      for (var col = 0; col < $scope.n_columns; col++) {
        for (var row = 0; row < $scope.n_rows; row++) {
          var elem = document.getElementById(col + '-' + row);
          elem.setAttribute('class', 'cell new') 
        }
      }
      setInterval($scope.nextGen, 50);
    }

    $scope.nextGen = function() {
      var url = "/next?state=" + $scope.canonical;
      $http.get(url).then(function(response) {
          $scope.canonical = response.data
          $scope.fromCanonical($scope.canonical);
          //$scope.toBoard();
        }, function(response) {
          return $q.reject(response.data.error);
      });

    }

    $scope.initMap = function(n_columns, n_rows) {
      var map = Array(n_columns);
      var temp = Array(n_columns);
      for (var col = 0; col < n_columns; col++) {
        map[col] = new Array(n_rows);
        temp[col] = new Array(n_rows);
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
      $scope.map_size = temp;
      $scope.canonical = $scope.toCanonical();
      //$scope.toBoard();
    }

    $scope.toBoard = function() {
      console.log('starting');
      var map = '<ul class="columns">'
      for (var col = 0; col < $scope.n_columns; col++) {
        map += '<li class="col"><ul>'
        for (var row = 0; row < $scope.n_rows; row++) {
          elem = $scope.getElem(col, row);
          map += '<li class="cell"><img src="/static/images/'+ elem + '"></img></li>';
        } 
        map += '</ul></li>';
      } 
      map += '</ul>';
      $scope.board = map;
    }

}]);

