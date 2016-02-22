package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

type Field uint

func bit(num uint) uint {
	num = num - 1
	if num < 0 || num > 8 {
		log.Fatalf("setting bit above 8")
	}
	return 1 << num
}

func (f Field) has(num uint) bool {
	return (uint(f) & bit(num)) != 0
}

func (f Field) add(num uint) Field {
	return f.set(num, true)
}

func (f Field) remove(num uint) Field {
	return f.set(num, false)
}

func (f Field) members() []uint {
	var res []uint
	for i := uint(1); i < 10; i++ {
		if f.has(i) {
			res = append(res, i)
		}
	}
	return res
}

func (f Field) removeChecked(num uint) (Field, error) {
	f2 := f.set(num, false)
	if f2 == 0 {
		return f, fmt.Errorf("%s.removeChecked(%d) == 0", f, num)
	}
	return f2, nil
}

func (f Field) set(num uint, val bool) Field {
	if val {
		return Field(uint(f) | bit(num))
	} else {
		return Field(uint(f) & ^bit(num))
	}
}

func (f Field) String() string {
	var buf [9]byte
	for i := uint(0); i < 9; i++ {
		if f.has(i + 1) {
			buf[i] = byte('1' + i)
		} else {
			buf[i] = ' '
		}
	}
	return string(buf[:])
}

var AllPossible = Field(0).
	add(1).
	add(2).
	add(3).
	add(4).
	add(5).
	add(6).
	add(7).
	add(8).
	add(9)

type Sudoku []Field

func NewSudoku() Sudoku {
	s := make(Sudoku, 9*9)
	for i := range s {
		s[i] = AllPossible
	}
	return s
}

func Parse(in string) (Sudoku, error) {
	s := make(Sudoku, 9*9)
	for i := range s {
		s[i] = AllPossible
	}

	y := 0
	for _, line := range strings.Split(in, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		x := 0
		for _, field := range line {
			if field == '_' {
				x++
				continue
			}
			if field < '0' || field > '9' {
				continue
			}
			num := uint(field - '0')
			if err := s.SetField(x, y, num); err != nil {
				return nil, err
			}
			x++
		}
		y++
	}

	return s, nil
}

func (s Sudoku) coord(x, y int) int {
	return y*9 + x
}

func (s Sudoku) GetField(x, y int) Field {
	return s[s.coord(x, y)]
}

func (s Sudoku) remove(x, y int, num uint) error {
	f := s[s.coord(x, y)]
	f2, err := f.removeChecked(num)
	if err != nil {
		return fmt.Errorf("(%d, %d): %s", x, y, err)
	}
	s[s.coord(x, y)] = f2
	if f2 != f && len(f2.members()) == 1 {
		s.SetField(x, y, f2.members()[0])
	}
	return nil
}

func (s Sudoku) SetField(x, y int, num uint) error {
	f := s[s.coord(x, y)]
	if !f.has(num) {
		return fmt.Errorf("(%d, %d): %s.set(%d)", x, y, f, num)
	}
	s[s.coord(x, y)] = Field(0).add(num)
	for i := 0; i < 9; i++ {
		if i != y {
			if err := s.remove(x, i, num); err != nil {
				return err
			}
		}
		if i != x {
			if err := s.remove(i, y, num); err != nil {
				return err
			}
		}
	}
	startX := (x / 3) * 3
	startY := (y / 3) * 3
	for xx := startX; xx < startX+3; xx++ {
		for yy := startY; yy < startY+3; yy++ {
			if xx != x && yy != y {
				if err := s.remove(xx, yy, num); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s Sudoku) Solve() (Sudoku, error) {
	return s.solve(0)
}

func (s Sudoku) solve(at int) (Sudoku, error) {
	for i, v := range s[at:] {
		i = i + at
		x := i % 9
		y := i / 9
		if choices := v.members(); len(choices) > 1 {
			var lastErr error
			for _, num := range choices {
				s2 := append(Sudoku(nil), s...)
				// fmt.Printf("Solve choice (%d, %d) = %d -> %d\n", x, y, choices, num)
				if err := s2.SetField(x, y, num); err != nil {
					continue
				}
				// fmt.Printf("Solve step:\n%s", s2)
				// os.Stdin.Read([]byte{0})
				var res Sudoku
				res, lastErr = s2.solve(i + 1)
				if lastErr == nil {
					return res, nil
				}
			}
			return s, fmt.Errorf("No choice at (%d, %d): %s", x, y, lastErr)
		}
	}
	return s, nil
}

func (s Sudoku) String() string {
	var buf bytes.Buffer
	for i, v := range s {
		if i%27 == 0 {
			for i := 0; i < 90; i++ {
				buf.WriteRune('-')
			}
			buf.WriteRune('\n')
		}
		buf.WriteString(v.String())
		buf.WriteRune('|')
		if i%9 == 8 {
			buf.WriteRune('\n')
		}
	}
	return buf.String()
}

func (s Sudoku) CompactString() string {
	var buf bytes.Buffer
	for i, v := range s {
		if i%9 != 0 && i%3 == 0 {
			buf.WriteRune(' ')
		}
		m := v.members()
		if len(m) == 1 {
			buf.WriteRune(rune('0' + m[0]))
		} else {
			buf.WriteRune('_')
		}
		if i%9 == 8 {
			buf.WriteRune('\n')
		}
		if i%27 == 26 {
			buf.WriteRune('\n')
		}
	}
	return buf.String()
}

func main() {
	s, err := Parse(`
    _2_9__4__
    9_54__1__
    _63__8___

    ___1___67
    ____4____
    38___5___

    ___8__97_
    __7__38_1
    __9__4_5_`)
	if err != nil {
		log.Fatalf("incorrect sudoku: %s", err)
	}
	s2, err := s.Solve()
	if err != nil {
		log.Fatalf("incorrect sudoku: %s", err)
	}
	fmt.Printf("Solved:\n%s", s2.CompactString())
}
