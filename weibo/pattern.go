package main

import (
	"regexp"
	"strings"

	"github.com/thinxer/actnet"
)

// Pattern Matching

var (
	patternsRaw = []string{
		`(?P<verb>WORD/v) (?P<object>WORD/[nr])`,
		`(?P<verb>WORD/v)(?P<object>)`,
	}
	patterns []*regexp.Regexp
)

func init() {
	for _, str := range patternsRaw {
		re := strings.Replace(str, " ", `\s`, -1)
		re = strings.Replace(str, "WORD", `[^\s/]*`, -1)
		// log.Printf("Compiling pattern: [%s] [%s]", str, re)
		patterns = append(patterns, regexp.MustCompile(re))
	}
}

func zipdict(keys, values []string) map[string]string {
	ret := make(map[string]string)
	for i, k := range keys {
		ret[k] = values[i]
	}
	return ret
}

func stripsuffix(s string) string {
	return strings.SplitN(s, "/", 2)[0]
}

// extract activities from a tagged line, i.e. pattern matching.
// e.g.
// 好久/m 没/d 骑/v 自行车/n 了/y
// will yield Activity{骑 自行车}
func extract(line string) (vn *actnet.Text) {
	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(line)
		if len(matches) > 0 {
			dict := zipdict(pattern.SubexpNames(), matches)
			verb := stripsuffix(dict["verb"])
			object := stripsuffix(dict["object"])
			return &actnet.Text{Verb: verb, Object: object}
		}
	}
	return nil
}
