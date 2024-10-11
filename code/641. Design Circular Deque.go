package code

type MyCircularDeque struct {
	arr  []int
	size int
	head int
	tail int
	cnt  int
}

func ConstructorCirQueue(k int) MyCircularDeque {
	return MyCircularDeque{
		arr:  make([]int, k),
		size: k,
		head: -1,
		tail: -1,
		cnt:  0,
	}
}

func (this *MyCircularDeque) InsertFront(value int) bool {
	if this.IsFull() {
		return false
	}
	if this.IsEmpty() {
		this.head, this.tail = 0, 0
		this.arr[0] = value
		this.cnt++
		return true
	}
	hn := (this.head - 1 + this.size) % this.size
	this.arr[hn] = value
	this.head = hn
	this.cnt++
	return true
}

func (this *MyCircularDeque) InsertLast(value int) bool {
	if this.IsFull() {
		return false
	}
	if this.IsEmpty() {
		this.head, this.tail = 0, 0
		this.arr[0] = value
		this.cnt++
		return true
	}
	tn := (this.tail + 1) % this.size
	this.arr[tn] = value
	this.tail = tn
	this.cnt++
	return true
}

func (this *MyCircularDeque) DeleteFront() bool {
	if this.IsEmpty() {
		return false
	}
	this.head = (this.head + 1) % this.size
	this.cnt--
	if this.IsEmpty() {
		this.head, this.tail = -1, -1
	}
	return true
}

func (this *MyCircularDeque) DeleteLast() bool {
	if this.IsEmpty() {
		return false
	}
	this.tail = (this.tail - 1 + this.size) % this.size
	this.cnt--
	if this.IsEmpty() {
		this.head, this.tail = -1, -1
	}
	return true
}

func (this *MyCircularDeque) GetFront() int {
	if !this.IsEmpty() {
		return this.arr[this.head]
	}

	return -1
}

func (this *MyCircularDeque) GetRear() int {
	if !this.IsEmpty() {
		return this.arr[this.tail]
	}

	return -1
}

func (this *MyCircularDeque) IsEmpty() bool {
	return this.cnt == 0
}

func (this *MyCircularDeque) IsFull() bool {
	return this.cnt == this.size
}

/**
 * Your MyCircularDeque object will be instantiated and called as such:
 * obj := Constructor(k);
 * param_1 := obj.InsertFront(value);
 * param_2 := obj.InsertLast(value);
 * param_3 := obj.DeleteFront();
 * param_4 := obj.DeleteLast();
 * param_5 := obj.GetFront();
 * param_6 := obj.GetRear();
 * param_7 := obj.IsEmpty();
 * param_8 := obj.IsFull();
 */
