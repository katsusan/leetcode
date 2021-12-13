package stubs

import (
	"math/rand"
	"time"
)

/*
see: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle, also known as Knuth shuffle.

- Write down the numbers from 1 through N.
- Pick a random number k between one and the number of unstruck numbers remaining (inclusive).
- Counting from the low end, strike out the kth number not yet struck out, and write it down at the end of a separate list.
- Repeat from step 2 until all the numbers have been struck out.
- The sequence of numbers written down in step 3 is now a random permutation of the original numbers.

可以这样证明：
	洗牌的公平性要保证每个元素在每个位置上的概率均等。以[N]int为例，概率为1/N。
	第一次在[1,N]中随机选择一个数，选取对应索引处的数放在最后一个位置，概率为1/N，也就是每个元素在最后一个位置的概率均等。
	第二次在[1,N-1]中随机一个数，放在倒数第二个位置，此时最后一个元素在该位置的概率为1/N,
		其余元素在该位置的概率为(N-1)/N * 1/(N-1) = 1/N, 其中(N-1)/N为第一次没被选中的概率，1/(N-1)为第二次被选中的概率。
	其余元素同理。

*/

func init() {
	rand.Seed(time.Now().UnixNano())
}

func FisherYates(s []int) {
	for i := len(s) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}
