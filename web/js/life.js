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

    $scope.visited_colors = ['#EEE', '#DDD', '#CCC', '#BBB', '#AAA', '#999', '#888', "#777" ,'#666', '#555','#444', '#333', '#222', '#111'];


    $scope.getCol = function(col) {
      return $scope.map[col];
    }

    $scope.getElem = function(col, row) {
      if ($scope.map[col][row] == ALIVE)
        return 'cell.gif';
      else
        return 'dead.jpg'; 
    }
    
    $scope.toCanonical = function() {
      var canonical = $scope.n_columns + "," + $scope.n_rows + "|";
      for (var col = 0; col < $scope.n_columns; col++) {
        for (var row = 0; row < $scope.n_rows; row++) {
          //console.log($scope.map[col][row]);
          if ($scope.map[col][row].alive == ALIVE) {
            canonical += col +',' + row + ',';
          }
        } 
      } 
      return canonical;
    }

    $scope.color = function(col, row) {
      var cell = $scope.map[col][row];
      if (cell.alive == 1) {
        return 'red';
      }
      else {
        if (cell.n_dead < $scope.visited_colors.length) {
          return $scope.visited_colors[cell.n_dead];
        }
        else {
          return $scope.visited_colors[$scope.visited_colors.length-1];
        }
      }
    }

    $scope.fromCanonical = function(canonical) {
      var t = canonical.split("|");
      var size = t[0].split(",");
      var n_columns = parseInt(size[0]);
      var n_rows = parseInt(size[1]);
      var elems = t[1];
      $scope.alive = elems.split(',');

      //var cells = d3.select('.cell');
      //console.log(cells); 
      // KILL OF ALIVE CELLS FROM PREV
      for (var i=0; i< $scope.prev.length-1; i=i+2) {
        var col = parseInt($scope.prev[i]);
        var row = parseInt($scope.prev[i+1]);
        $scope.map[col][row].alive = 0;
        $scope.map[col][row].n_dead += 1;
        //d3.select('#C'+ col + 'R' + row).attr("class", "cell dead").attr("n_dead", function(d) {
          //console.log(d.n_dead);
          //return d.n_dead;
        //}).style("fill", function(d) {return $scope.color(d.n_dead);});
      }

      //LOOP THROUGH THE COMPUTED ALIVE AND ADD THEM TO MAP
      for (var i=0; i<$scope.alive.length-1; i=i+2) {
        var col = parseInt($scope.alive[i]);
        var row = parseInt($scope.alive[i+1]);
        $scope.map[col][row].alive = 1;
        //d3.select('#C'+ col + 'R' + row).attr("class", "cell alive");
      }

      $scope.prev = $scope.alive;
    }

    $scope.loop = function() {
        window.setInterval($scope.nextGen, 100);
    }

    $scope.nextGen = function() {
      var url = "/next?state=" + $scope.canonical;
      console.log('next gen');
      $http.get(url).then(function(response) {
          console.log('setting canonical');
          $scope.canonical = response.data
          console.log('from canonical');
          $scope.fromCanonical($scope.canonical);
          $scope.updateGrid('#chart')
          //$scope.toBoard();
        }, function(response) {
          return $q.reject(response.data.error);
      });

    }

    $scope.selectMap = function(mapName) {
      var url = "/maps?mapName=" + mapName;
      console.log(mapName);
      console.log(url);
      $http.get(url).then(function(response) {
          $scope.canonical = response.data
          $scope.fromCanonical($scope.canonical);
          $scope.updateGrid('#chart')
          //$scope.toBoard();
        }, function(response) {
                return $q.reject(response.data.error);
      });
    }

    function CellGrid(rows, columns) {
      this.cells = new Array(rows);
      var n = 0;
      for(var i = -1; ++i < rows;){
          this.cells[i] = new Array(columns);
          for(var j = -1; ++j < columns; ){
              var cell = new Cell(false);
              cell.n = n++;
              cell.x = i;
              cell.y = j;
              this.cells[i][j] = cell;
          }
      }
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
    }

    $scope.updateGrid = function(id) {
      console.log(id);
      id = id + " svg";
      var grid = d3.select('#world');

      console.log(grid);
      var row = grid.selectAll(".row")
                    .data($scope.map)
                    .selectAll(".cell")
                   .data(function (d) { return d; })
                   .attr("class", "cell new")
                   .attr("x", function(d) { return d.x; })
                   .attr("y", function(d) { return d.y; })
                   .attr("width", function(d) { return d.width; })
                   .attr("height", function(d) { return d.height; })
                   .attr("id", function(d) {return d.index; })
                   .attr("n_dead", function(d) { return d.n_dead; })
                   .style("fill", function(d) {return $scope.color(d.col, d.row); })
    }


    $scope.generateGrid = function (id, width, height, n_columns, n_rows)
    {
        $scope.map = $scope.randomData(width, height, n_columns, n_rows);
        $scope.n_columns = n_columns; 
        $scope.n_rows = n_rows;
        $scope.width = width;
        $scope.height = height;
        var grid = d3.select(id).append("svg")
                      .attr("width", $scope.width)
                      .attr("height", $scope.width)
                      .attr("class", "chart")
                      .attr("id", "world");
        console.log(grid);

        var row = grid.selectAll(".row")
                    .data($scope.map)
                  .enter().append("svg:g")
                    .attr("class", "row");

        var col = row.selectAll(".cell")
                   .data(function (d) { return d; })
                  .enter().append("svg:rect")
                   .attr("class", "cell new")
                   .attr("x", function(d) { return d.x; })
                   .attr("y", function(d) { return d.y; })
                   .attr("width", function(d) { return d.width; })
                   .attr("height", function(d) { return d.height; })
                   .attr("id", function(d) {return d.index; })
                   .attr("n_dead", function(d) { return d.n_dead; })
                   .style("fill", function(d) {return $scope.color(d.col, d.row); })
        $scope.canonical = $scope.toCanonical();
    }

    $scope.randomData = function (gridWidth, gridHeight, n_columns, n_rows)
    {
        $scope.map = new Array();
        $scope.alive = new Array();
        var gridItemWidth = gridWidth / n_columns;
        var gridItemHeight = gridHeight / n_rows;
        var startX = gridItemWidth / 2;
        var startY = gridItemHeight / 2;
        var stepX = gridItemWidth;
        var stepY = gridItemHeight;
        var xpos = startX;
        var ypos = startY;
        var newValue = 0;
        var count = 0;
        var status = 0;

        for (var index_a = 0; index_a < n_columns; index_a++)
        {
            $scope.map.push(new Array());
            for (var index_b = 0; index_b < n_rows; index_b++)
            {
                newValue = Math.random();
                if (newValue > 0.7) {
                  status = ALIVE;
                  $scope.alive.push(index_a);
                  $scope.alive.push(index_b);
                }
                else {
                  status = DEAD;
                }
                 $scope.map[index_a].push({ 
                                    col: index_a,
                                    row: index_b,
                                    alive: status,
                                    width: gridItemWidth,
                                    height: gridItemHeight,
                                    x: xpos,
                                    y: ypos,
                                    count: count,
                                    index: 'C' + index_a + 'R' + index_b,
                                    n_dead: 0
                                });
                
                xpos += stepX;
                count += 1;
            }
            xpos = startX;
            ypos += stepY;
        }
        $scope.prev = $scope.alive;
        return  $scope.map;
    }

}]);

