package code

import "slices"

func minimumDeletions(word string, k int) int {
	freqMap := make(map[byte]int)
	for i := range word {
		freqMap[word[i]]++
	}

	freqs := make([]int, 0, len(freqMap))
	for _, v := range freqMap {
		freqs = append(freqs, v)
	}

	lowest, highest := slices.Min(freqs), slices.Max(freqs)
	minSteps := 0
	for i := lowest; i+k <= highest; i++ {
		// use window [low, high] to scan all freq
		low, high := i, i+k
		steps := 0
		for _, freq := range freqs {
			if freq < low {
				steps += freq
			} else if freq > high {
				steps += (freq - high)
			}
		}
		if minSteps == 0 {
			minSteps = steps
		} else {
			minSteps = min(minSteps, steps)
		}
	}

	return minSteps
}
