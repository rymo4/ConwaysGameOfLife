package universe

import (
	"bufio"
	"bytes"
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

type grid map[string]byte

type Universe struct {
	Width        int
	Height       int
	Space        grid
	generation   int
	initialState string
}

func (u *Universe) IsLiving(y, x int) bool {
	return u.Space[toKey(y, x)] == living
}

func (u *Universe) SetLiving(y, x int, alive bool) {
	key := toKey(y, x)
	if alive {
		u.Space[key] = living
	} else {
		delete(u.Space, key)
	}
}

func toKey(y, x int) string {
	return fmt.Sprintf("%d-%d", y, x)
}

func CoordsFromKey(key string) (y, x int) {
	parts := strings.Split(key, "-")
	y_, _ := strconv.ParseInt(parts[0], 10, 0)
	x_, _ := strconv.ParseInt(parts[1], 10, 0)
	x, y = int(x_), int(y_)
	return
}

func (u *Universe) String() string {
	var buffer bytes.Buffer
	for y := 0; y < u.Height; y++ {
		for x := 0; x < u.Width; x++ {
			if u.IsLiving(y, x) {
				buffer.WriteString(string(living))
			} else {
				buffer.WriteString(".")
			}
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (u *Universe) NeighborsCount(y, x int) int {
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
		Space:  make(grid),
	}
}

func (self *Universe) Clone() Universe {
	return Universe{
		Width:  self.Width,
		Height: self.Height,
		Space:  make(grid),
	}
}

func (u *Universe) AtGeneration(gen int) {
	for u.generation < gen {
		u.Next()
	}
	//u.Show()
}

func (u *Universe) Next() {
	nxGen := u.Clone()
	for key := range u.Space {
		yCenter, xCenter := CoordsFromKey(key)
		for y_ := -1; y_ < 2; y_++ {
			for x_ := -1; x_ < 2; x_++ {
				y := y_ + yCenter
				x := x_ + xCenter
				if x < 0 || x >= u.Width || y < 0 || y >= u.Height {
					continue
				}
				live := u.IsLiving(y, x)
				switch n := u.NeighborsCount(y, x); {
				case live && n < 2:
					nxGen.SetLiving(y, x, false)
				case live && (n == 2 || n == 3):
					nxGen.SetLiving(y, x, true)
				case live && n > 3:
					nxGen.SetLiving(y, x, false)
				case !live && n == 3:
					nxGen.SetLiving(y, x, true)
				default:
					nxGen.SetLiving(y, x, false)
				}
			}
		}
	}
	u.Space, nxGen.Space = nxGen.Space, u.Space
	u.generation++
}

func LoadFromString(str string) (u *Universe) {
	str = strings.Trim(str, "\n ")
	body := strings.Split(str, "\n")

	rows := len(body)
	cols := len(body[0])

	u = &Universe{Width: cols, Height: rows, Space: make(grid)}
	for i := range body {
		for j, e := range body[i] {
			if e != blankMarker {
				u.SetLiving(i, j, true)
			}
		}
	}
	return u
}

func LoadFromFile(path string) (u *Universe, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines += scanner.Text()
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return LoadFromString(lines), nil
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

		uni.SetLiving(int(y), int(x), true)
	}
	return uni
}

func (u *Universe) CanonicalString() string {
	state := fmt.Sprintf("%d,%d|", u.Width, u.Height)
	for key := range u.Space {
		i, j := CoordsFromKey(key)
		if u.IsLiving(i, j) {
			state = fmt.Sprintf("%s%d,%d,", state, j, i)
		}
	}
	return state
}
