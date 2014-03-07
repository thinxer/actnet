package actnet

import (
	"encoding/json"
	"io"
)

type MergeResult struct {
	Total, Valid, Accepted int
}

func (m *Model) Merge(fin io.Reader) (MergeResult, error) {
	var r MergeResult

	decoder := json.NewDecoder(fin)
	var item ExtractedItem
	for {
		err := decoder.Decode(&item)
		if err == io.EOF {
			break
		} else if err != nil {
			return r, err
		}

		if len(item.Name) > 0 {
			r.Valid++
			if m.Insert(&item) {
				r.Accepted++
			}
		}
		r.Total++
	}
	return r, nil
}
