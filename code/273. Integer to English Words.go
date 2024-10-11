package code

import (
	"fmt"
	"strings"
)

func numberToWords(num int) string {
	prmap := map[int]string{
		0:          "Zero",
		1:          "One",
		2:          "Two",
		3:          "Three",
		4:          "Four",
		5:          "Five",
		6:          "Six",
		7:          "Seven",
		8:          "Eight",
		9:          "Nine",
		10:         "Ten",
		11:         "Eleven",
		12:         "Twelve",
		13:         "Thirteen",
		14:         "Fourteen",
		15:         "Fifteen",
		16:         "Sixteen",
		17:         "Seventeen",
		18:         "Eighteen",
		19:         "Nineteen",
		20:         "Twenty",
		30:         "Thirty",
		40:         "Forty",
		50:         "Fifty",
		60:         "Sixty",
		70:         "Seventy",
		80:         "Eighty",
		90:         "Ninety",
		100:        "Hundred",
		1000:       "Thousand",
		1000000:    "Million",
		1000000000: "Billion",
	}

	var w strings.Builder
	doNumtoWords(&w, num, prmap)
	return strings.TrimSpace(w.String())
}

func doNumtoWords(w *strings.Builder, n int, m map[int]string) {
	// over 1 billion?
	nbillion := n / 1000000000
	if nbillion != 0 {
		fmt.Fprintf(w, "%s %s ", ntowForK(nbillion, m), m[1e9])
		if n%1e9 != 0 {
			doNumtoWords(w, n-nbillion*1e9, m)
		}
		return
	}

	// over 1 million?
	nmillion := n / 1000000
	if nmillion != 0 {
		fmt.Fprintf(w, "%s %s ", ntowForK(nmillion, m), m[1e6])
		if n%1e6 != 0 {
			doNumtoWords(w, n-nmillion*1e6, m)
		}
		return
	}

	// over 1k?
	nk := n / 1000
	if nk != 0 {
		fmt.Fprintf(w, "%s %s ", ntowForK(nk, m), m[1e3])
		if n%1000 != 0 {
			doNumtoWords(w, n-nk*1e3, m)
		}
		return
	}

	// less than 1k
	fmt.Fprintf(w, "%s", ntowForK(n, m))
}

// fast method for nums that < 1000
func ntowForK(n int, m map[int]string) string {
	var w strings.Builder
	// n = a * 100 + b * 10 + c
	a := n / 100
	if a != 0 {
		fmt.Fprintf(&w, "%s %s ", m[a], m[100])
	}

	b := (n % 100) / 10
	if b != 0 {
		if b == 1 {
			fmt.Fprintf(&w, "%s", m[n%100])
		} else {
			if n%10 == 0 {
				fmt.Fprintf(&w, "%s", m[b*10])
			} else {
				fmt.Fprintf(&w, "%s %s", m[b*10], m[n%10])
			}
		}
	} else {
		if a == 0 || n%10 != 0 {
			fmt.Fprintf(&w, "%s", m[n%10])
		}
	}

	return strings.TrimSpace(w.String())
}
