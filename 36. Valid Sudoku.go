package main

/*
Determine if a 9x9 Sudoku board is valid. Only the filled cells need to be validated according to the following rules:
	1. Each row must contain the digits 1-9 without repetition.
	2. Each column must contain the digits 1-9 without repetition.
	3. Each of the 9 3x3 sub-boxes of the grid must contain the digits 1-9 without repetition.
*/
func isValidSudoku(board [][]byte) bool {
	//first check the 9 rows
	for _, row := range board {
		if !isValidAll(row) {
			return false
		}
	}

	//check the 9 columns
	reverseBoard := make([][]byte, 9)
	for i := 0; i < 9; i++ {
		reverseBoard[i] = make([]byte, 9)
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			reverseBoard[i][j] = board[j][i]
		}
	}

	for _, row := range reverseBoard {
		if !isValidAll(row) {
			return false
		}
	}

	//check the 9 3x3 sub-boxes
	for _, i := range []int{1, 4, 7} {
		for _, j := range []int{1, 4, 7} {
			subbox := []byte{board[i-1][j-1], board[i-1][j], board[i-1][j+1],
				board[i][j-1], board[i][j], board[i][j+1],
				board[i+1][j-1], board[i+1][j], board[i+1][j+1]}
			if !isValidAll(subbox) {
				return false
			}
		}
	}

	return true
}

//isValidAll will check all elements in target to determine if repetitions exists
func isValidAll(target []byte) bool {
	s := NewSet()
	for _, elem := range target {
		if s.Exist(elem) {
			return false
		}
		if elem != '.' {
			s.Add(elem)
		}
	}
	return true
}

type Set struct {
	data map[byte]struct{}
}

func NewSet() *Set {
	return &Set{
		data: make(map[byte]struct{}),
	}
}

//check if n exists in set
func (s *Set) Exist(n byte) bool {
	if _, ok := s.data[n]; ok {
		return true
	}
	return false
}

func (s *Set) Add(n byte) {
	s.data[n] = struct{}{}
}
