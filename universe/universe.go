package universe

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	living      = 'O'
	dead        = ' '
	blankMarker = '.'
)

type Universe struct {
	Width        int
	Height       int
	Space        [][]byte
	generation   int
	initialState string
}

func (u *Universe) IsLiving(y, x int) bool {
  return u.Space[y][x] == living
}

func (u *Universe) SetLiving(y, x int, alive bool) {
  if alive {
    u.Space[y][x] = living
  } else {
    u.Space[y][x] = dead
  }
}

func (u *Universe) Show() {
	fmt.Print("\x0c")
	for _, r := range u.Space {
		for i, b := range r {
			if b == 0 {
				r[i] = dead
			}
		}
		fmt.Println(string(r[:]))
	}
}

func (u *Universe) NeighborsCount(x, y int) int {
	n := 0
	for dx := -1; dx < 2; dx++ {
		for dy := -1; dy < 2; dy++ {
			xx, yy := x+dx, y+dy
			if (dx != 0 || dy != 0) &&
				xx >= 0 && xx < u.Width &&
				yy >= 0 && yy < u.Height && u.IsLiving(yy, xx) {
				n++
			}
		}
	}
	return n
}

func New(width, height int) *Universe {
	return &Universe{
		Width:  width,
		Height: height,
		Space:  newSpaceArray(height, width),
	}
}

func (self *Universe) Clone() Universe {
	return Universe{
		Width:  self.Width,
		Height: self.Height,
		Space:  newSpaceArray(self.Height, self.Width),
	}
}

func (u *Universe) AtGeneration(gen int) {
	for u.generation < gen {
		u.Next()
	}
	u.Show()
}

func (u *Universe) Next() {
	nxGen := u.Clone()
	for y := range u.Space {
		for x := range u.Space[0] {
			live := u.Space[y][x] == living
			switch n := u.NeighborsCount(x, y); {
			case live && n < 2:
				nxGen.Space[y][x] = dead
			case live && (n == 2 || n == 3):
				nxGen.Space[y][x] = living
			case live && n > 3:
				nxGen.Space[y][x] = dead
			case !live && n == 3:
				nxGen.Space[y][x] = living
			default:
				nxGen.Space[y][x] = dead
			}
		}
	}
	u.Space, nxGen.Space = nxGen.Space, u.Space
	u.generation++
}

func LoadFromFile(path string) (u *Universe, err error) {
	body, err := readLines(path)
	if err != nil {
		return nil, err
	}

	rows := len(body)
	cols := len(body[0])

	u = &Universe{Width: cols, Height: rows, Space: newSpaceArray(rows, cols)}

	for i := range body {
		for j, e := range body[i] {
			if e != blankMarker {
				u.Space[i][j] = living
			}
		}
	}
	return u, err
}

func newSpaceArray(rows, cols int) [][]byte {
	ar := make([][]byte, rows)
	for i := range ar {
		ar[i] = make([]byte, cols)
	}
	return ar
}

func LoadFromCanonicalString(state string) *Universe {
	parts := strings.Split(state, "|")
	mapSize := strings.Split(parts[0], ",")
	livingCells := strings.Split(parts[1], ",")

	//TODO: handle errors
	width, _ := strconv.ParseInt(mapSize[0], 10, 0)
	height, _ := strconv.ParseInt(mapSize[1], 10, 0)

	uni := New(int(width), int(height))
	for i := 0; i < len(livingCells)/2; i++ {
		xIndex := i * 2
		yIndex := xIndex + 1

		x, _ := strconv.ParseInt(livingCells[xIndex], 10, 0)
		y, _ := strconv.ParseInt(livingCells[yIndex], 10, 0)

		uni.Space[y][x] = living
	}
	return uni
}

func (u *Universe) CanonicalString() string {
	state := fmt.Sprintf("%d,%d|", u.Width, u.Height)
	for i, r := range u.Space {
		for j, _ := range r {
			if u.IsLiving(i, j) {
				state = fmt.Sprintf("%s%d,%d,", state, j, i)
			}
		}
	}
	return state
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
