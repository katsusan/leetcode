package main

/**
 * Your MyHashSet object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Add(key);
 * obj.Remove(key);
 * param_3 := obj.Contains(key);
 */

type MyHashSet struct {
	bits []byte
}

/** Initialize your data structure here. */
func HashSetConstructor() MyHashSet {
	return MyHashSet{
		bits: make([]byte, 0),
	}
}

func (this *MyHashSet) Add(key int) {
	if len(this.bits) <= key/8 {
		newbits := make([]byte, key/8+1)
		copy(newbits, this.bits)
		this.bits = newbits
	}

	this.bits[key/8] |= (1 << uint(key%8))
}

func (this *MyHashSet) Remove(key int) {
	if len(this.bits) <= key/8 {
		return
	}
	this.bits[key/8] &= ^(1 << uint(key%8))
}

/** Returns true if this set contains the specified element */
func (this *MyHashSet) Contains(key int) bool {
	if len(this.bits) <= key/8 {
		return false
	}

	cur := this.bits[key/8] >> uint(key%8) & 1
	if cur != 1 {
		return false
	}
	return true
}
