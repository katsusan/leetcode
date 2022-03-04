package code

import "strings"

// Note: I didn't do some sanity check in isValidBoard and stringReplace,
// as these are internal functions which should be only used in this file
// and I can be sure there are no problems in bound check.

func solveNQueens(n int) [][]string {
	var solutions = make([][]string, 0, n)
	initBoard := make([]string, n)
	for i := range initBoard {
		initBoard[i] = strings.Repeat(".", n)
	}
	backtrackNQueens(&solutions, initBoard, 0)
	return solutions
}

func backtrackNQueens(solu *[][]string, board []string, row int) {
	if row == len(board) {
		res := make([]string, len(board))
		copy(res, board)
		*solu = append(*solu, res)
		return
	}

	for col := 0; col < len(board); col++ {
		if !isValidBoard(board, row, col) {
			continue
		}

		board[row] = stringReplace(board[row], col, 'Q')
		backtrackNQueens(solu, board, row+1)
		board[row] = stringReplace(board[row], col, '.')
	}
}

// isValid reports whether position (x, y) can be a valid place for board.
func isValidBoard(board []string, x, y int) bool {
	// check row
	if strings.Contains(board[x], "Q") {
		return false
	}

	// check column and diagonal
	for i := 0; i < x; i++ {
		if board[i][y] == 'Q' {
			return false
		}

		for j := 0; j < len(board); j++ {
			if IntAbs(j-y) == x-i && board[i][j] == 'Q' {
				return false
			}
		}
	}
	return true
}

func stringReplace(s string, i int, b byte) string {
	tmp := []byte(s)
	tmp[i] = b
	return string(tmp)
}

func IntAbs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}
