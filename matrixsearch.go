// Updated matrixsearch.go
package matrixsearch

import (
	"fmt"
	"math/rand"
	"os/exec"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type DataStore[T any] struct {
	mu             sync.RWMutex
	items          map[string]T
	compositeIndex map[string][]string
	getID          func(T) string
	indexer        func(T) []string
}

func NewDataStore[T any](getID func(T) string, indexer func(T) []string) *DataStore[T] {
	return &DataStore[T]{
		items:          make(map[string]T),
		compositeIndex: make(map[string][]string),
		getID:          getID,
		indexer:        indexer,
	}
}

func getCombinations(keys []string) []string {
	var combs []string
	n := len(keys)
	for i := 1; i < (1 << n); i++ {
		var subset []string
		for j := 0; j < n; j++ {
			if i&(1<<j) != 0 {
				subset = append(subset, keys[j])
			}
		}
		combs = append(combs, strings.Join(subset, ":"))
	}
	return combs
}

func (ds *DataStore[T]) Insert(item T) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	id := ds.getID(item)
	ds.items[id] = item
	keys := ds.indexer(item)
	comps := getCombinations(keys)
	for _, key := range comps {
		ds.compositeIndex[key] = append(ds.compositeIndex[key], id)
	}
}

func (ds *DataStore[T]) Delete(item T) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	id := ds.getID(item)
	delete(ds.items, id)
	keys := ds.indexer(item)
	comps := getCombinations(keys)
	for _, key := range comps {
		newIDs := []string{}
		for _, itemID := range ds.compositeIndex[key] {
			if itemID != id {
				newIDs = append(newIDs, itemID)
			}
		}
		ds.compositeIndex[key] = newIDs
	}
}

func (ds *DataStore[T]) Search(query string) []T {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if ids, ok := ds.compositeIndex[query]; ok {
		var results []T
		for _, id := range ids {
			results = append(results, ds.items[id])
		}
		return results
	}
	return nil
}

func (ds *DataStore[T]) SearchRandom(query string) (T, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if ids, ok := ds.compositeIndex[query]; ok && len(ids) > 0 {
		n := rand.Intn(len(ids))
		return ds.items[ids[n]], true
	}
	var zero T
	return zero, false
}

func (ds *DataStore[T]) Update(item T) {
	ds.Delete(item)
	ds.Insert(item)
}

func (ds *DataStore[T]) Count() int {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return len(ds.items)
}

func (ds *DataStore[T]) Clear() {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.items = make(map[string]T)
	ds.compositeIndex = make(map[string][]string)
}

func AutoIndexer[T any](item T) []string {
	var keys []string
	v := reflect.ValueOf(item)
	keys = extractKeys(v)
	return keys
}

func extractKeys(v reflect.Value) []string {
	var keys []string
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return keys
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get("text")
		if tag != "" {
			var strVal string
			switch value.Kind() {
			case reflect.String:
				strVal = value.String()
			case reflect.Bool:
				strVal = strconv.FormatBool(value.Bool())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				strVal = strconv.FormatInt(value.Int(), 10)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				strVal = strconv.FormatUint(value.Uint(), 10)
			case reflect.Float32, reflect.Float64:
				strVal = strconv.FormatFloat(value.Float(), 'f', -1, 64)
			}
			if strVal != "" {
				keys = append(keys, tag+":"+strVal)
			}
		}
		if value.Kind() == reflect.Struct {
			subKeys := extractKeys(value)
			keys = append(keys, subKeys...)
		}
	}
	return keys
}

// escapeDOT escapes double quotes in strings for DOT.
func escapeDOT(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}

/// Completely AI generated code below this line. No clue what it does but looks great!.

func (ds *DataStore[T]) Dump(filename string) error {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	var b strings.Builder
	b.WriteString("digraph DataStore {\n")
	b.WriteString("  // Basic settings\n")
	b.WriteString("  rankdir=TB;\n")         // Top to bottom layout
	b.WriteString("  splines=polyline;\n")   // Straight line segments
	b.WriteString("  ranksep=0.8;\n")        // Spacing between ranks
	b.WriteString("  nodesep=0.5;\n")        // Spacing between nodes
	b.WriteString("  fontname=\"Arial\";\n") // Default font
	b.WriteString("  node [fontname=\"Arial\", fontsize=11];\n")
	b.WriteString("  edge [fontname=\"Arial\", fontsize=9, arrowsize=0.7];\n")

	// Use HTML-like labels for better formatting
	// Create a structured section for stats
	totalKeys := len(ds.compositeIndex)

	// Count unique items
	itemSet := make(map[string]bool)
	for _, ids := range ds.compositeIndex {
		for _, id := range ids {
			itemSet[id] = true
		}
	}
	totalItems := len(itemSet)

	// Add a stats header node
	b.WriteString("  // Stats header node\n")
	b.WriteString(fmt.Sprintf("  \"stats\" [shape=plaintext, label=<<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\"><TR><TD BGCOLOR=\"#E6E6FA\"><B>DataStore Statistics</B></TD></TR><TR><TD ALIGN=\"left\">Total Keys: %d</TD></TR><TR><TD ALIGN=\"left\">Total Items: %d</TD></TR></TABLE>>, fontsize=12];\n",
		totalKeys, totalItems))

	// Create a key category node
	b.WriteString("  // Key category node\n")
	b.WriteString("  \"keyCategory\" [shape=plaintext, label=<<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\"><TR><TD BGCOLOR=\"#D0E0FF\"><B>Composite Keys</B></TD></TR></TABLE>>, fontsize=12];\n")

	// Connect stats to key category
	b.WriteString("  \"stats\" -> \"keyCategory\" [style=invis];\n")

	// Group keys by their length/complexity to create categories
	keysByLength := make(map[int][]string)

	for key := range ds.compositeIndex {
		// Count the number of components in each key
		components := strings.Count(key, ":") + 1
		keysByLength[components] = append(keysByLength[components], key)
	}

	// Sort the categories
	var lengths []int
	for length := range keysByLength {
		lengths = append(lengths, length)
	}
	sort.Ints(lengths)

	// Create key complexity level nodes
	var complexityNodes []string
	for _, length := range lengths {
		nodeName := fmt.Sprintf("keyLevel_%d", length)
		label := fmt.Sprintf("%d-Component Keys", length)
		if length == 1 {
			label = "Simple Keys"
		}
		b.WriteString(fmt.Sprintf("  \"%s\" [shape=plaintext, label=<<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\"><TR><TD BGCOLOR=\"#E0EFFF\">%s</TD></TR></TABLE>>, fontsize=11];\n",
			nodeName, label))
		complexityNodes = append(complexityNodes, nodeName)

		// Connect from key category
		b.WriteString(fmt.Sprintf("  \"keyCategory\" -> \"%s\";\n", nodeName))

		// Sort keys within this level for consistency
		sort.Strings(keysByLength[length])

		// Add individual key nodes
		for _, key := range keysByLength[length] {
			safeKey := escapeDOT(key)
			count := len(ds.compositeIndex[key])

			// Create a node for this key
			b.WriteString(fmt.Sprintf("  \"%s\" [shape=box, style=\"rounded,filled\", fillcolor=\"#F0F8FF\", label=\"%s\\n(%d items)\"];\n",
				safeKey, safeKey, count))

			// Connect from the complexity level
			b.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\";\n", nodeName, safeKey))

			// If this key has a reasonable number of items, show them directly
			MAX_DIRECT_ITEMS := 5 // Limit for direct connections
			ids := ds.compositeIndex[key]

			if len(ids) <= MAX_DIRECT_ITEMS {
				// Show all items directly
				for _, id := range ids {
					safeID := escapeDOT(id)
					b.WriteString(fmt.Sprintf("  \"%s_item_%s\" [shape=ellipse, style=\"filled\", fillcolor=\"#FFE6E6\", label=\"%s\"];\n",
						safeKey, safeID, safeID))
					b.WriteString(fmt.Sprintf("  \"%s\" -> \"%s_item_%s\";\n", safeKey, safeKey, safeID))
				}
			} else {
				// Create a collapsed node that can be expanded in interactive viewers
				b.WriteString(fmt.Sprintf("  \"%s_items\" [shape=folder, style=\"filled\", fillcolor=\"#FFEFEF\", label=\"%d items\"];\n",
					safeKey, len(ids)))
				b.WriteString(fmt.Sprintf("  \"%s\" -> \"%s_items\";\n", safeKey, safeKey))
			}
		}
	}

	// Make sure all complexityNodes are in the same rank
	if len(complexityNodes) > 1 {
		b.WriteString("  { rank=same; ")
		for _, node := range complexityNodes {
			b.WriteString(fmt.Sprintf("\"%s\"; ", node))
		}
		b.WriteString("}\n")
	}

	// Ensure stats and key category are at the top
	b.WriteString("  { rank=source; \"stats\"; \"keyCategory\"; }\n")

	// Closing the DOT representation
	b.WriteString("}\n")

	// Execute Graphviz command to generate the SVG output
	dotStr := b.String()
	cmd := exec.Command("dot", "-Tsvg", "-o", filename)
	cmd.Stdin = strings.NewReader(dotStr)
	return cmd.Run()
}
