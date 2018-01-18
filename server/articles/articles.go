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

func LoadArticles() Articles {
	bs, err := dry.FileGetBytes(JSON_FILE_PATH)
	utils.NoError(err)
	ars := Articles{}
	err = json.Unmarshal(bs, &ars)
	utils.NoError(err)
	return ars
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

func (self Articles) Sort() {
	sort.SliceStable(self, func(i, j int) bool { return self[i].DateUpdated.After(self[j].DateUpdated) })
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
