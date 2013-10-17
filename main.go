package main

import (
  "fmt"
  "bufio"
  "time"
  "os"
)

const (
  living    = 'O'
  dead      = ' '
  blankMarker = '.'
  framerate = 10
)

type universe struct {
  Width int
  Height int
  Space [][]byte
}

func loadFromFile(path string) (u* universe, err error) {
  body, err := readLines(path)
  if err != nil {
    return nil, err
  }

  rows := len(body)
  cols := len(body[0])

  u = &universe{Width: cols, Height: rows, Space: newSpaceArray(rows, cols)}

  for i := range body {
    for j, e := range body[i] {
      if e != blankMarker {
        u.Space[i][j] = living
      }
    }
  }
  return u, err
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

func (u *universe) Show() {
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

func (u *universe) neighbors(x, y int) (n int) {
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
  return
}

func newSpaceArray(rows, cols int) [][]byte {
  ar := make([][]byte, rows)
  for i := range ar {
    ar[i] = make([]byte, cols)
  }
  return ar
}

func (self *universe) Clone() universe {
  return universe{Width: self.Width, Height: self.Height, Space: newSpaceArray(self.Height, self.Width)}
}

func main() {
  u, err := loadFromFile("maps/glider_gun.txt")
  nxGen := u.Clone()

  if err != nil {
    fmt.Printf("Please provide a valid file")
    return
  }

  for i := 0; i < 300; i++ {
    u.Show()
    for y := range u.Space {
      for x := range u.Space[0] {
        live := u.Space[y][x] == living
        switch n := u.neighbors(x, y); {
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
    time.Sleep(time.Second / framerate)
  }
}
