package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestFieldStringer(t *testing.T) {
	cases := []struct {
		f Field
		s string
	}{
		{Field(0), "         "},
		{AllPossible, "123456789"},
		{Field(0).add(1).add(2), "12       "},
		{Field(0).add(1).add(9), "1       9"},
		{AllPossible.remove(1), " 23456789"},
		{AllPossible.remove(1).remove(5), " 234 6789"},
	}

	for _, c := range cases {
		if c.f.String() != c.s {
			t.Errorf("String(%b): got %s, expected %s", c.f, c.f.String(), c.s)
		}
	}
}

func TestPrint(t *testing.T) {
	s := NewSudoku()
	res := s.String()
	lines := strings.Split(res, "\n")
	if len(lines) != 13 {
		t.Errorf("got %d lines, expected %d", len(lines), 13)
	}
	if !strings.HasPrefix(lines[9], "123456789|") {
		t.Errorf("got %s as last line, expected 1234... on it", lines[9])
	}
}

func TestSetField(t *testing.T) {
	s := NewSudoku()
	if err := s.SetField(0, 0, 1); err != nil {
		t.Error(err)
	}
	if err := s.SetField(0, 8, 2); err != nil {
		t.Error(err)
	}
	if err := s.SetField(8, 0, 3); err != nil {
		t.Error(err)
	}

	cases := []struct {
		x, y int
		f    Field
	}{
		{0, 0, Field(0).add(1)},
		{2, 2, AllPossible.remove(1)},
		{0, 8, Field(0).add(2)},
		{0, 7, AllPossible.remove(1).remove(2)},
		{8, 0, Field(0).add(3)},
	}

	for _, c := range cases {
		f := s.GetField(c.x, c.y)
		if f != c.f {
			t.Errorf("(%d, %d): got %s, expected %s", c.x, c.y, f, c.f)
		}
	}

	err := s.SetField(0, 0, 2)
	if err == nil {
		t.Errorf("Expected error setting field")
	}
}

func TestRemoveChecked(t *testing.T) {
	if _, err := Field(0).add(1).removeChecked(2); err != nil {
		t.Error(err)
	}
}

func TestMembers(t *testing.T) {
	cases := []struct {
		f Field
		e []uint
	}{
		{AllPossible, []uint{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{Field(0), nil},
		{Field(0).add(4), []uint{4}},
	}

	for _, c := range cases {
		if !reflect.DeepEqual(c.f.members(), c.e) {
			t.Errorf("%s.members() got %s, expected %s", c.f, c.f.members(), c.e)
		}
	}
}

func TestParseSudoku(t *testing.T) {
	s, err := Parse(
		`
1_____4__
__3______
_________

_____6__9
_6_______
__5_4____

_________
___5____6
2________
`)
	if err != nil {
		t.Fatalf("Parse failed: %s", err)
	}
	if s.GetField(0, 0) != Field(0).add(1) {
		t.Error("Expected (0, 0) = 1")
	}
}

func TestSolve(t *testing.T) {
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
		t.Fatal(err)
	}

	res, err := s.Solve()
	if err != nil {
		t.Fatalf("err: %s\n%s", err, res)
	}

	s2, err := Parse(`
    128 957 436
    975 436 182
    463 218 795

    254 189 367
    796 342 518
    381 675 249

    532 861 974
    647 593 821
    819 724 653`)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, s2) {
		t.Errorf("Solve: got:\n%sExpected:\n%s", res.CompactString(), s2.CompactString())
	}
}
