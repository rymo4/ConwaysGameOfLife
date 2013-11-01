var app = angular.module('Conways',[]);

var Universe = {
  empty: function(height, width) {
    var map = new Array(height);
    for (var y = 0; y < height; y++){
      map[y] = new Array(width);
      for (var x = 0; x < width; x++){
        map[y][x] = {
          x: x,
          y: y,
          count: 0,
          alive: false,
          n_dead: 0
        };
      }
    }
    return { map: map, alive: [] };
  },
  fromCanonical: function(canonical) {
    var parts  = canonical.split("|");
    var info   = parts[0].split(",");
    var width  = parseInt(info[0]);
    var height = parseInt(info[1]);

    universe = Universe.empty(height, width);
    universe = Universe._addFromCanonical(universe, canonical);
    return universe;
  },
  update: function(universe, canonical) {
    for (var i = 0; i < universe.alive.length - 1; i = i + 2) {
      var col = parseInt(universe.alive[i]);
      var row = parseInt(universe.alive[i+1]);
      universe.map[row][col].alive = false;
      universe.map[row][col].n_dead += 1;
    }
    universe.alive = [];
    universe = Universe._addFromCanonical(universe, canonical);
    Universe.render(universe);
    return universe;
  },
  _addFromCanonical: function(universe, canonical){
    var parts  = canonical.split("|");
    var alive = parts[1].split(',');

    for (var i=0; i < alive.length-1; i=i+2) {
      var col = parseInt(alive[i]);
      var row = parseInt(alive[i+1]);
      universe.alive.push(col);
      universe.alive.push(row);
      universe.map[row][col].alive = true;
    }
    universe.canonical = canonical;
    return universe;
  },
  initRender: function(universe) {
    universe.cells = new Array(universe.map.length * universe.map[0].length);
    for (var i=0; i < universe.cells.length; i++) {
      var y = i % universe.map.length;
      var x = i % universe.map[0].length;
      universe.cells[i] = universe.map[y][x];
    }

    var svg = d3.select("body").append("svg").attr("width", 500).attr("height", 500);
    Universe.rects = svg.selectAll("rect")
      .data(universe.cells)
      .enter().append("rect");
    Universe.render(universe);
  },
  render: function(universe) {
    Universe.rects
      .attr('x', function(d,i){ return i % universe.map[0].length * 10;})
      .attr('y', function(d,i){ return i % universe.map.length * 10;})
      .attr('width', 10)
      .attr('height', 10)
      .attr('fill', function(d){
        if (d.alive) return 'red';
        return 'white';
      });
  }
};

app.controller('LifeCtrl', ['$scope', '$http', '$q', function($scope, $http, $q) {
    $http.get("/mapslist").then(function(response) {
      $scope.maps = response.data;
      $scope.selectMap($scope.maps[0]);
    });

    $scope.visited_colors = [
      '#EEE', '#DDD', '#CCC', '#BBB',
      '#AAA', '#999', '#888', "#777",
      '#666', '#555', '#444', '#333',
      '#222', '#111'
    ];

    $scope.loop = function() {
        window.setInterval($scope.nextGen, 50);
    }

    $scope.nextGen = function() {
      var url = "/next?state=" + $scope.universe.canonical;
      $http.get(url).then(function(response) {
          var canonical = response.data
          Universe.update($scope.universe, canonical);
        }, function(response) {
          return $q.reject(response.data.error);
      });
    }

    $scope.selectMap = function(mapName) {
      var url = "/maps?mapName=" + mapName;
      $http.get(url).then(function(response) {
        $scope.universe = Universe.fromCanonical(response.data);
        Universe.initRender($scope.universe);
      }, function(response) {
        return $q.reject(response.data.error);
      });
    }
}]);
