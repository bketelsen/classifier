package cmd

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/burntsushi/toml"
)

type Vendor struct {
	Name      string
	Directory string
	Keywords  []string
	Regex     string
}
type Vendors struct {
	Vendor []Vendor
}

type Match struct {
	Vendor Vendor
	File   os.FileInfo
	Year   string
	Month  string
}

func (v Vendor) KeywordMatch(s string) bool {
	var match bool
	for _, key := range v.Keywords {
		if strings.Contains(s, key) {
			match = true
			if !match {
				return false
			}
		} else {
			return false
		}
	}
	return match
}

func getMapping(config string) (Vendors, error) {
	var mpng Vendors
	bb, err := ioutil.ReadFile(config)
	if err != nil {
		return Vendors{}, err
	}
	if _, err := toml.Decode(string(bb), &mpng); err != nil {
		return Vendors{}, err
	}
	return mpng, nil
}
