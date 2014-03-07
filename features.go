package actnet

import (
	"fmt"
	"sort"

	"github.com/thinxer/go-word2vec"
)

type FeatureExtractor struct {
	FeatureMap map[string]int

	wordVectors *word2vec.Model
}

func (f *FeatureExtractor) featureId(name string) int {
	if id, ok := f.FeatureMap[name]; ok {
		return id
	} else {
		id = len(f.FeatureMap)
		f.FeatureMap[name] = id
		return id
	}
}

func (f *FeatureExtractor) ExtractSingle(a *Activity) (features pairs) {
	var vector word2vec.Vector
	if id, ok := f.wordVectors.Vocab[a.Verb]; ok {
		vector = f.wordVectors.Vector(id)
	} else {
		panic("Word not found!" + a.Verb)
	}
	if len(a.Object) > 0 {
		if id, ok := f.wordVectors.Vocab[a.Object]; ok {
			vector.Add(1.0, f.wordVectors.Vector(id))
		}
	}
	vector.Normalize()

	for i, v := range vector {
		features = append(features, pair{f.featureId(fmt.Sprintf("v%d", i)), v})
	}

	// s := m.Summary(a)
	// for cat, kv := range s {
	// 	if cat == "next" {
	// 		continue
	// 	}
	// 	sum := float32(0)
	// 	for _, v := range kv {
	// 		sum += float32(v)
	// 	}
	// 	for k, v := range kv {
	// 		features = append(features, pair{getId(fmt.Sprintf("%s_%s", cat, k)), float32(v) / sum})
	// 	}
	// }
	sort.Sort(features)

	return
}

type pair struct {
	id    int
	value float32
}
type pairs []pair

func (p pairs) Len() int           { return len(p) }
func (p pairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pairs) Less(i, j int) bool { return p[i].id < p[j].id }
