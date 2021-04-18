package code

import "fmt"

func solveSudoku(board [][]byte) {
	if solveSodukuEx(board, -1, -1) == false {
		fmt.Println("solution failed")
	}
}

//if given soduku can be solved, then returns true,
//otherwise returns false
func solveSodukuEx(board [][]byte, previ, prevj int) bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] == '.' {
				possdigits := getPossibleDigits(board, i, j)
				//means no right solution
				if len(possdigits) == 0 {
					if previ >= 0 && prevj >= 0 {
						board[previ][prevj] = '.'
					}
					return false
				}
				for _, pb := range possdigits {
					board[i][j] = pb
					if solveSodukuEx(board, i, j) {
						return true
					}
				}
				if previ >= 0 && prevj >= 0 {
					board[previ][prevj] = '.'
				}
				return false
			}
		}
	}
	return true
}

func getPossibleDigits(board [][]byte, i, j int) []byte {
	exclusionSet := make([]byte, 0)
	s := NewSet() //store the already occured digits
	for k := 0; k < 9; k++ {
		if board[i][k] >= '1' && board[i][k] <= '9' {
			s.Add(board[i][k]) //add digits at the same row
		}
		if board[k][j] >= '1' && board[k][j] <= '9' {
			s.Add(board[k][j]) //add digits at the same column
		}
	}

	//add the digits in the same 3x3 sub-box
	for _, m := range []int{0, 1, 2} {
		for _, n := range []int{0, 1, 2} {
			cur := board[i/3*3+m][j/3*3+n]
			if cur >= '1' && cur <= '9' {
				s.Add(cur)
			}
		}
	}

	for e := byte('1'); e <= '9'; e++ {
		if !s.Exist(e) {
			exclusionSet = append(exclusionSet, e)
		}
	}
	return exclusionSet
}
