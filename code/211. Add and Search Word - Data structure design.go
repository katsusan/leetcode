package code

type WordDictionary struct {
	wordtree *TrieNode //root node of Trie
}

/*
type TrieNode struct {
	val      byte
	wordflag bool
	childNodes   []*TrieNode
}
*/
func (node *TrieNode) Add(word string) {
	if len(word) == 0 {
		node.wordflag = true
		return
	}

	for _, child := range node.childNodes {
		if child.val == word[0] {
			//already exists
			child.Add(word[1:])
			return
		}
	}

	newchild := &TrieNode{
		val:        word[0],
		wordflag:   false,
		childNodes: []*TrieNode{},
	}
	node.childNodes = append(node.childNodes, newchild)
	newchild.Add(word[1:])
}

func (node *TrieNode) Search(word string) bool {
	if len(word) == 0 {
		return node.wordflag
	}

	switch word[0] {
	case '.':
		for _, child := range node.childNodes {
			if child.Search(word[1:]) {
				return true
			}
		}
	default:
		for _, child := range node.childNodes {
			if child.val == word[0] {
				return child.Search(word[1:])
			}
		}
	}

	return false
}

/** Initialize your data structure here. */
func trieConstructor() WordDictionary {
	return WordDictionary{
		wordtree: &TrieNode{childNodes: []*TrieNode{}},
	}
}

/** Adds a word into the data structure. */
func (this *WordDictionary) AddWord(word string) {
	this.wordtree.Add(word)
}

/** Returns if the word is in the data structure. A word could contain the dot character '.' to represent any one letter. */
func (this *WordDictionary) Search(word string) bool {
	return this.wordtree.Search(word)
}

/**
 * Your WordDictionary object will be instantiated and called as such:
 * obj := Constructor();
 * obj.AddWord(word);
 * param_2 := obj.Search(word);
 */
