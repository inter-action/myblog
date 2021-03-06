package articles

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/ungerik/go-dry"

	"github.com/inter-action/myblog/server/utils"
)

var JSON_FILE_PATH = ".articles.tmp"

var mutex = &sync.Mutex{}

func LoadArticles() (Articles, os.FileInfo) {
	info, err := os.Stat(JSON_FILE_PATH)
	utils.NoError(err)
	bs, err := dry.FileGetBytes(JSON_FILE_PATH)
	utils.NoError(err)
	ars := Articles{}
	err = json.Unmarshal(bs, &ars)
	utils.NoError(err)
	return ars, info
}

func Persist(ars Articles) {
	mutex.Lock()
	defer mutex.Unlock()

	json, err := json.Marshal(ars)
	utils.NoError(err)
	err = dry.FileSetBytes(JSON_FILE_PATH, json)
	utils.NoError(err)
}

type Articles []*Article

func ParseArticles(path string) (Articles, error) {
	ars := make(Articles, 0, 30)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".md" {
			if ar, err := ParseArticle(path); err != nil {
				return nil
			} else if ar.Title == "" { // skip any articles without title field
				return nil
			} else {
				ars = append(ars, ar)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return ars, nil
}

func (self Articles) Sort(sortFunc func(i, j int) bool) {
	sort.SliceStable(self, sortFunc)
}

func (self Articles) SortByCreated() {
	self.Sort(func(i, j int) bool { return self[i].DateCreated.After(self[j].DateCreated) })
}

func (self Articles) SortByUpdated() {
	self.Sort(func(i, j int) bool { return self[i].DateUpdated.After(self[j].DateUpdated) })
}

func (self Articles) Paginate(page int32, pageSize int32) Articles {
	offset := (page - 1) * pageSize
	limit := page * pageSize

	size := int32(len(self))
	if limit > size {
		limit = size
	}

	return self[offset:limit]
}

func (self Articles) Update(path string) {
	ars, err := ParseArticles(path)
	utils.NoError(err)
	Persist(ars)
	self = ars
}
