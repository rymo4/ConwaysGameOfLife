package main

import (
  "fmt"
  "bufio"
  "time"
  "os"
)

const (
  living = 'O'
)

//w := 79
//h := 15

//type universe [h][w]byte

type universe struct {
  Width int
  Height int
  Space [][]byte
}

func loadFromFile(path string) (u* universe, err error) {
  body, err := readLines(path)
  if err != nil {
    // TODO: make new err
    return nil, err
  }

  // Find size
  rows := len(body)
  cols := len(body[0])

  fmt.Printf("rows: %d, cols: %d", rows, cols)

  u = &universe{Width: rows, Height: cols}
  u.Space = newSpaceArray(rows, cols)

  for i := range body {
    for j, e := range body[0] {
      if e == 'X' {
        u.Space[j][i] = living
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

func show(u *universe) {
  fmt.Print("\x0c")
  for _, r := range u.Space {
    for i, b := range r {
      if b == 0 {
        r[i] = ' '
      }
    }
    fmt.Println(string(r[:]))
  }
}

func neighbors(u *universe, x, y int) (n int) {
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
  ar := make([][]byte, cols)
  for i := range ar {
    ar[i] = make([]byte, rows)
  }
  return ar
}

func (self *universe) Clone() universe {
  return universe{Width: self.Width, Height: self.Height, Space: newSpaceArray(self.Width, self.Height)}
}

func main() {
  u, err := loadFromFile("maps/glider_gun.txt")
  nxGen := u.Clone()

  if err != nil {
    fmt.Printf("Please provide a valid file")
    return
  }

  for i := 0; i < 300; i++ {
    show(u)
    for y := range u.Space {
      for x := range u.Space[0] {
        live := u.Space[y][x] == living
        switch n := neighbors(u, x, y); {
        case live && n < 2:
          nxGen.Space[y][x] = 0
        case live && (n == 2 || n == 3):
          nxGen.Space[y][x] = living
        case live && n > 3:
          nxGen.Space[y][x] = 0
        case !live && n == 3:
          nxGen.Space[y][x] = living
        default:
          nxGen.Space[y][x] = 0
        }
      }
    }
    u.Space, nxGen.Space = nxGen.Space, u.Space
    time.Sleep(time.Second / 10)
  }
}
