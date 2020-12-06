package translation

import (
	avl "github.com/emirpasic/gods/trees/avltree"
	"fmt"
	"strings"
	"sync"
)

type HistoryCache struct {
	Store *avl.Tree
	sync.RWMutex
}

var (
	WordsHistory = &HistoryCache{
		Store: avl.NewWithStringComparator(),
	}

	SentencesHistory = &HistoryCache{
		Store: avl.NewWithStringComparator(),
	}
)

func (wh *HistoryCache) Cache(eng, gopher string) {
	wh.Lock()
	defer wh.Unlock()

	if _, found := wh.Store.Get(eng); found {
		return
	}

	wh.Store.Put(eng, gopher)
}

func ExportCachedItemsAsJson() string {
	getKeyValuePairsAsJsonObjects := func(wh *HistoryCache) string {
		if (wh.Store.Size() == 0) {
			return ""
		}

		wh.RLock()
		defer wh.RUnlock()
    
		b := strings.Builder{}

		// Make an inorder traversal of the tree
		tree := wh.Store
		lastNode := tree.Right()
		for n := tree.Left(); n != lastNode; n = n.Next() {
			b.WriteString(fmt.Sprintf(`{"%s":"%s"},`, n.Key, n.Value))
		}

		// The last element shouldn't have , at the end
		b.WriteString(fmt.Sprintf(`{"%s":"%s"}`, lastNode.Key, lastNode.Value))

		return b.String()
	}

	builder := strings.Builder{}
	builder.WriteString(`{"history":[`)

	var words = getKeyValuePairsAsJsonObjects(WordsHistory)
	if words != "" {
		builder.WriteString(words)
	}

	var sentences = getKeyValuePairsAsJsonObjects(SentencesHistory)
	if sentences != "" {
		if words != "" {
			builder.WriteString(",")
		}
		builder.WriteString(sentences)
	}

	builder.WriteString("]}")

	return builder.String()
}
