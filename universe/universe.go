package universe

import (
	"bufio"
	"fmt"
	"os"
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
				yy >= 0 && yy < u.Height && u.Space[yy][xx] == living {
				n++
			}
		}
	}
	return n
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
	body, err := ReadLines(path)
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

func ReadLines(path string) ([]string, error) {
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
