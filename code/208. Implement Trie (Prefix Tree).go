package code

type TrieNode struct {
	val        byte
	wordflag   bool
	childNodes []*TrieNode
}

type Trie struct {
	*TrieNode //root Node
}

/** Initialize your data structure here. */
func TrieConstructor() Trie {
	rootNode := &TrieNode{
		childNodes: make([]*TrieNode, 0),
	}
	return Trie{rootNode}
}

/** Inserts a word into the trie. */
func (this *Trie) Insert(word string) {
	if word == "" {
		this.wordflag = true
		return
	}

	//first make sure if there has already word[0] existed in childNodes
	for _, child := range this.childNodes {
		if child.val == word[0] {
			(&Trie{child}).Insert(word[1:])
			return
		}
	}

	//need add childNode
	newchild := &TrieNode{
		val:        word[0],
		wordflag:   false,
		childNodes: make([]*TrieNode, 0),
	}
	this.childNodes = append(this.childNodes, newchild)
	(&Trie{newchild}).Insert(word[1:])
	return
}

/** Returns if the word is in the trie. */
func (this *Trie) Search(word string) bool {
	if word == "" {
		return this.wordflag
	}

	//search childNodes for word[0]
	for _, child := range this.childNodes {
		if child.val == word[0] {
			return (&Trie{child}).Search(word[1:])
		}
	}
	return false
}

/** Returns if there is any word in the trie that starts with the given prefix. */
func (this *Trie) StartsWith(prefix string) bool {
	if prefix == "" {
		return true
	}

	//see if any child is equal to prefix[0]
	for _, child := range this.childNodes {
		if child.val == prefix[0] {
			return (&Trie{child}).StartsWith(prefix[1:])
		}
	}
	return false
}

/**
 * Your Trie object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Insert(word);
 * param_2 := obj.Search(word);
 * param_3 := obj.StartsWith(prefix);
 */
