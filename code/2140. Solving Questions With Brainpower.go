package code

func mostPoints(questions [][]int) int64 {
	// dp[i] = max(point[i] + dp[i+N], dp[i+1], ..., dp[i+N-1])
	nq := len(questions)
	dp := make([]int, nq)
	dp[nq-1] = questions[nq-1][0]
	for i := nq - 2; i >= 0; i-- {
		// two choices: solve q[i] or skip q[i]
		sov := 0
		if i+questions[i][1]+1 > nq-1 {
			sov = questions[i][0]
		} else {
			sov = questions[i][0] + dp[i+questions[i][1]+1]
		}

		dp[i] = max(dp[i+1], sov)
	}

	return int64(dp[0])
}
