package main

import (
	"encoding/json"
	"time"

	"github.com/thinxer/actnet"
)

type rawWeibo struct {
	Id                  int64
	Idstr               string
	UserId              int64 `json:"user_id"`
	Text                string
	Mid                 string
	BmiddlePic          string `json:"bmiddle_pic"`
	Source              string
	InReplyToUserId     string `json:"in_reply_to_user_id"`
	InReplyToStatusId   string `json:"in_reply_to_status_id"`
	InReplyToScreenName string `json:"in_reply_to_screen_name"`
	RepostsCount        int64  `json:"reposts_count"`
	Truncated           bool
	ThumbnailPic        string `json:"thumbnail_pic"`
	Favorited           bool
	CreatedAt           string `json:"created_at"`
	Mlevel              int64
	AttitudesCount      int64  `json:"attitudes_count"`
	OriginalPic         string `json:"original_pic"`
	CommentsCount       int64  `json:"comments_count"`

	Visible struct {
		Type   int64
		ListId int64 `json:"list_id"`
	}
	Geo struct {
		Type        string
		Coordinates []float64
	}
	Annotations []struct {
		Place struct {
			// The following two fields have string and float types, strange
			//Lat   float64
			//Lon   float64
			Type  string
			Poiid string
			Title string
		}
	}
}

type Weibo struct {
	Id        int64
	UserId    int64
	Text      string
	Segmented string
	Timestamp time.Time
	Location  actnet.POI
}

func nextWeibo(decoder *json.Decoder) (*Weibo, error) {
	var w Weibo
	var v rawWeibo

	err := decoder.Decode(&v)
	if err != nil {
		return nil, err
	}

	w.Id = v.Id
	w.Text = v.Text
	w.UserId = v.UserId
	w.Timestamp, err = time.Parse("Mon Jan 02 15:04:05 -0700 2006", v.CreatedAt)
	if err != nil {
		return nil, err
	}
	w.Location.Lat = v.Geo.Coordinates[0]
	w.Location.Lng = v.Geo.Coordinates[1]
	for _, ann := range v.Annotations {
		w.Location.Id = ann.Place.Poiid
		w.Location.Name = ann.Place.Title
	}
	return &w, nil
}
