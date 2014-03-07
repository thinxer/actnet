package actnet

import (
	"time"
)

type Point struct {
	Lat float64
	Lng float64
}

type POI struct {
	// URI to identify this POI.
	Id string
	// Position.
	Point
	// Extras.
	Name, Type, City string
}

type ExtractedItem struct {
	// Name of the extracted activity.
	Name string
	// Simple representation of the item, from the language perspective.
	Verb, Object string
	// Location information.
	POI POI
	// Time when this happened.
	Timestamp time.Time

	// Id of this user, prefer to a URI.
	User string
	// Id of the original item, URI preferred, too.
	Entity string
	// Text of original information. For human evaluation.
	OriginalText string
}
