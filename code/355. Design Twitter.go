package code

import (
	"sort"
	"time"
)

type Twitter struct {
	rmap map[int]map[int]int //relation map: user<->user
	emap map[int][]tweet     //twitter entry map: user <-> tweet
}

type tweet struct {
	tweetid  int
	posttime int64
}

const heapsize = 10

type heapTweet []tweet //define a minheap

func (h heapTweet) Len() int { return len(h) }

func (h heapTweet) Less(i, j int) bool { return !older(h[i], h[j]) }

func (h heapTweet) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func older(t1, t2 tweet) bool {
	if t1.posttime < t2.posttime {
		return true
	}
	return false
}

func (h *heapTweet) Push(t tweet) {
	if len(*h) < heapsize {
		//add to tail and swim
		*h = append(*h, t)
		k := len(*h) - 1
		for (k-1)/2 > 0 {
			if older((*h)[k], (*h)[(k-1)/2]) {
				//child < parent then swap
				(*h)[k], (*h)[(k-1)/2] = (*h)[(k-1)/2], (*h)[k]
			}
			k = (k - 1) / 2
		}
	} else {
		//compare t with top of heap and sink it if necessary
		if older((*h)[0], t) {
			(*h)[0] = t
			par := 0
			for 2*par+1 < len(*h) {
				chd := 2*par + 1
				if chd < len(*h)-1 && older((*h)[chd+1], (*h)[chd]) {
					chd++
				}
				if older((*h)[par], (*h)[chd]) {
					return
				}
				(*h)[par], (*h)[chd] = (*h)[chd], (*h)[par]
				par = chd
			}
		}
	}
}

/** Initialize your data structure here. */
func TwitterConstructor() Twitter {
	return Twitter{
		rmap: make(map[int]map[int]int, 0),
		emap: make(map[int][]tweet, 0),
	}
}

/** Compose a new tweet. */
func (this *Twitter) PostTweet(userId int, tweetId int) {

	if _, found := (*this).rmap[userId]; !found {
		(*this).rmap[userId] = make(map[int]int, 1)
		(*this).rmap[userId][userId] = 1
	}

	(*this).emap[userId] = append((*this).emap[userId], tweet{tweetId, time.Now().UnixNano()})
}

/** Retrieve the 10 most recent tweet ids in the user's news feed. Each item in the news feed must be posted by users who the user followed or by the user herself. Tweets must be ordered from most recent to least recent. */
func (this *Twitter) GetNewsFeed(userId int) []int {

	followmap, found := (*this).rmap[userId]
	if !found {
		return []int{}
	}

	//list of followed users and user herself
	userlist := make([]int, 0)

	for k := range followmap {
		userlist = append(userlist, k)
	}

	twheap := &heapTweet{}

	for _, user := range userlist {
		for _, tw := range (*this).emap[user] {
			twheap.Push(tw)
		}
	}
	sort.Sort(twheap)

	res := make([]int, 0)
	for _, t := range *twheap {
		res = append(res, t.tweetid)
	}
	return res
}

/** Follower follows a followee. If the operation is invalid, it should be a no-op. */
func (this *Twitter) Follow(followerId int, followeeId int) {
	_, foundFollower := (*this).rmap[followerId]

	if !foundFollower {
		(*this).rmap[followerId] = make(map[int]int, 0)
		(*this).rmap[followerId][followerId] = 1
	}

	(*this).rmap[followerId][followeeId] = 1
}

/** Follower unfollows a followee. If the operation is invalid, it should be a no-op. */
func (this *Twitter) Unfollow(followerId int, followeeId int) {
	if followerId == followeeId {
		return
	}

	if _, found := (*this).rmap[followerId]; found {
		delete((*this).rmap[followerId], followeeId)
	}
}

/**
 * Your Twitter object will be instantiated and called as such:
 * obj := Constructor();
 * obj.PostTweet(userId,tweetId);
 * param_2 := obj.GetNewsFeed(userId);
 * obj.Follow(followerId,followeeId);
 * obj.Unfollow(followerId,followeeId);
 */
