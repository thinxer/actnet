// The weibo data importer.
// Requirements:
//  sdfsdf
//  sdfsdf
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/thinxer/actnet"
)

var flagWeiboData = flag.String("data", "../data/weibo", "Path to weibo data.")
var flagMediumData = flag.String("output", "../output/weibo.med", "Path of output.")

func main() {
	var output io.Writer
	var err error

	flag.Parse()

	if *flagMediumData == "-" {
		output = os.Stdout
	} else {
		output, err = os.OpenFile(*flagMediumData, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
	}

	// Extract structs and save to some file.
	if err = ImportWeibo(*flagWeiboData, json.NewEncoder(output)); err != nil {
		panic(err)
	}
}

// ImportWeibo can read a weibo dataset.
func ImportWeibo(weibodir string, output *json.Encoder) (err error) {
	// Load POI
	pois, err := LoadPOIs(path.Join(weibodir, "weibo.poi"))
	log.Println(len(pois), err)

	// Process Weibo and Segmented Text
	n, err := iterateJsonsWithText(path.Join(weibodir, "data.jsons"), path.Join(weibodir, "pos.txt"), func(w *Weibo) {
		var item actnet.ExtractedItem
		item.POI = w.Location
		item.Entity = "weibo:" + strconv.FormatInt(w.Id, 10)
		item.User = "weibo:user:" + strconv.FormatInt(w.UserId, 10)
		item.OriginalText = w.Text
		item.Timestamp = w.Timestamp
		if a := extract(w.Segmented); a != nil {
			item.Name = a.Verb + a.Object
			item.Verb = a.Verb
			item.Object = a.Object
		}
		output.Encode(item)
	})

	log.Println("Total:", n)
	return
}

// iterate over jsons with segmented text
func iterateJsonsWithText(jsonsFile, segmentedFile string, fn func(*Weibo)) (int, error) {
	type (
		M map[string]interface{}
		V []interface{}
	)

	jsons, err := os.Open(jsonsFile)
	if err != nil {
		return 0, err
	}
	defer jsons.Close()
	segmented, err := os.Open(segmentedFile)
	if err != nil {
		return 0, err
	}
	defer segmented.Close()

	decoder := json.NewDecoder(jsons)
	reader := bufio.NewReader(segmented)
	count := 0
	for {
		w, err := nextWeibo(decoder)
		if err != nil {
			if err == io.EOF {
				break
			} else if errt, ok := err.(*time.ParseError); ok {
				log.Printf("Line %d parse error: %v, breaking", count, errt)
				break
			} else if errj, ok := err.(*json.UnmarshalTypeError); ok {
				log.Printf("Line %d json error: %v", count, errj)
				continue
			} else {
				return count, err
			}
		}
		w.Segmented, err = reader.ReadString('\n')
		w.Segmented = strings.Trim(w.Segmented, "\r\n")
		if err != nil {
			return count, err
		}
		fn(w)
		count++
		if count%10000 == 0 {
			log.Printf("Processed: %d\n", count)
		}
	}
	return count, nil
}
