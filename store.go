package actnet

import (
	"encoding/gob"
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"
)

type Model struct {
	Items      []ExtractedItem
	Keys       map[string]int
	Activities map[string][]int
	Users      map[string][]int
	Days       map[int][]int

	Year    [3000]int
	Month   [12]int
	YearDay [366]int
	Hour    [24]int

	LastUpdate time.Time
}

func NewModel() *Model {
	return &Model{
		Keys:       make(map[string]int),
		Activities: make(map[string][]int),
		Users:      make(map[string][]int),
		Days:       make(map[int][]int),
	}
}

func (m *Model) String() string {
	return fmt.Sprintf("Activities: %d, Items: %d. Last Update: %v", len(m.Activities), len(m.Items), m.LastUpdate)
}

func (m *Model) Insert(e *ExtractedItem) bool {
	if _, ok := m.Keys[e.Entity]; ok {
		return false
	}

	m.Year[e.Timestamp.Year()]++
	m.Month[int(e.Timestamp.Month())-1]++
	m.YearDay[e.Timestamp.YearDay()]++
	m.Hour[e.Timestamp.Hour()]++

	i := len(m.Items)
	m.Items = append(m.Items, *e)
	m.Keys[e.Entity] = i
	m.Activities[e.Name] = append(m.Activities[e.Name], i)
	m.Users[e.User] = append(m.Users[e.User], i)
	daysSinceEpoch := int(e.Timestamp.Unix() / 86400)
	m.Days[daysSinceEpoch] = append(m.Days[daysSinceEpoch], i)
	return true
}

func (m *Model) Load(in io.Reader) error {
	return gob.NewDecoder(in).Decode(m)
}

func (m *Model) Save(out io.Writer) error {
	m.LastUpdate = time.Now()
	return gob.NewEncoder(out).Encode(m)
}

func (m *Model) PrintStats() {
	seen := make(map[string]bool)
	today := int(time.Now().Unix()/86400) + 1
	for i := 0; i <= today; i++ {
		todayList := []string{}

		ids := m.Days[i]
		todayNew := 0
		for _, id := range ids {
			name := m.Items[id].Name
			if !seen[name] {
				seen[name] = true
				todayNew++
				todayList = append(todayList, name)
			}
		}
		if todayNew > 0 {
			fmt.Println(time.Unix(int64(i*86400), 0), todayNew, len(ids), todayList)
		}
	}
}

func (m *Model) Summary(name string) *Activity {
	ids := m.Activities[name]
	if len(ids) < 5 {
		return nil
	}

	a := &Activity{}
	a.Verb = m.Items[ids[0]].Verb
	a.Object = m.Items[ids[0]].Object
	a.Summary = ActivitySummary{}

	for _, i := range ids {
		item := m.Items[i]
		a.Summary.Count("POIType", item.POI.Type, 1)
		a.Summary.Count("City", item.POI.City, 1)
		a.Summary.Count("Month", strconv.Itoa(int(item.Timestamp.Month())), 1)
		a.Summary.Count("Weekday", strconv.Itoa(int(item.Timestamp.Weekday())), 1)
		a.Summary.Count("Hour", strconv.Itoa(item.Timestamp.Hour()), 1)
	}

	return a
}

func (m *Model) Casuality() map[string][]string {
	following := make(map[string][]string)
	for _, ids := range m.Users {
		m.sortIdByTime(ids)
		for i, id := range ids {
			cut := m.Items[id].Timestamp.Add(time.Hour * 3)
			for j := i + 1; j < len(ids); j++ {
				id2 := ids[j]
				if cut.Before(m.Items[id2].Timestamp) {
					break
				}
				following[m.Items[id].Name] = append(following[m.Items[id].Name], m.Items[id2].Name)
			}
		}
	}
	return following
}

type xslice struct {
	sort.IntSlice
	items []ExtractedItem
}

func (x xslice) Less(i, j int) bool {
	return x.items[i].Timestamp.Before(x.items[j].Timestamp)
}

func (m *Model) sortIdByTime(ids []int) {
	sort.Sort(xslice{sort.IntSlice(ids), m.Items})
}

type Activity struct {
	Name string

	Verb, Object string
	Summary      ActivitySummary
}

type ActivitySummary map[string]map[string]float32

func (a ActivitySummary) Count(class, value string, count float32) {
	if a[class] == nil {
		a[class] = make(map[string]float32)
	}
	a[class][value] += count
}

func (a ActivitySummary) Flatten() (r map[string]float32) {
	r = make(map[string]float32)
	for c, m := range a {
		for k, v := range m {
			r[c+"_"+k] = v
		}
	}
	return
}
