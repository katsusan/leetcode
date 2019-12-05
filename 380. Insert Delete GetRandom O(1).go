package main

import "math/rand"

type RandomizedSet struct {
	setm map[int]int //map who stores set element in key
	seta []int       //array for getRandom in O(1)
}

/** Initialize your data structure here. */
func RandomizedSetConstructor() RandomizedSet {
	return RandomizedSet{
		setm: make(map[int]int),
		seta: make([]int, 0),
	}
}

/** Inserts a value to the set. Returns true if the set did not already contain the specified element. */
func (this *RandomizedSet) Insert(val int) bool {
	if _, ok := this.setm[val]; ok {
		//already exists
		return false
	}

	this.seta = append(this.seta, val)
	index := len(this.seta) - 1
	this.setm[val] = index
	return true
}

/** Removes a value from the set. Returns true if the set contained the specified element. */
func (this *RandomizedSet) Remove(val int) bool {
	index, ok := this.setm[val]
	if !ok {
		//val doesn't exist
		return false
	}

	lena := len(this.seta)
	this.setm[this.seta[lena-1]] = index
	this.seta[index] = this.seta[lena-1]
	this.seta = this.seta[:lena-1]
	delete(this.setm, val)
	return true
}

/** Get a random element from the set. */
func (this *RandomizedSet) GetRandom() int {
	return this.seta[rand.Intn(len(this.seta))]
}

/**
 * Your RandomizedSet object will be instantiated and called as such:
 * obj := Constructor();
 * param_1 := obj.Insert(val);
 * param_2 := obj.Remove(val);
 * param_3 := obj.GetRandom();
 */
