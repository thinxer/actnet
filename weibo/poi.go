package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/thinxer/actnet"
)

type rawPOI struct {
	PoiId        string
	Title        string
	Lat          string
	Lon          string
	Address      string
	City         string
	Province     string
	Categorys    string
	CategoryName string `json:"category_name"`

	Country          interface{}
	PoiStreetSummary string `json:"poi_street_summary"`
	PoiStreetAddress string `json:"poi_street_address"`
	Postcode         string
	Category         string
	Enterprise       int64
	TodoNum          int64 `json:"todo_num"`
	// This line can be both int and string, WTF.
	//CheckinUserNum   int64 `json:"checkin_user_num,string"`
	PoiPic     string `json:"poi_pic"`
	CheckinNum int64  `json:"checkin_num"`
}

func LoadPOIs(jsons string) (map[string]actnet.POI, error) {
	log.Println("Loading POIs")
	file, err := os.Open(jsons)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pois := make(map[string]actnet.POI)
	decoder := json.NewDecoder(file)
	count := 0
	for {
		count++
		var v rawPOI
		err = decoder.Decode(&v)

		if err == io.EOF {
			break
		} else if errj, ok := err.(*json.UnmarshalTypeError); ok {
			log.Printf("Line %d json error: %v", count, errj)
			continue
		} else if err != nil {
			return pois, err
		}
		if len(v.Title) <= 1 {
			continue
		}
		poi := actnet.POI{}
		poi.Id = v.PoiId
		poi.Lat, _ = strconv.ParseFloat(v.Lat, 64)
		poi.Lng, _ = strconv.ParseFloat(v.Lon, 64)
		poi.Name = v.Title
		poi.Type = v.CategoryName
		poi.City = v.City
		pois[v.PoiId] = poi
	}
	log.Printf("POIs loaded, %d in total, %d valid.\n", count-1, len(pois))
	return pois, nil
}
