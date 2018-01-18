package articles

import "testing"

func load() Articles {
	ars, err := ParseArticles("./__test__")
	if err != nil {
		panic(err)
	}
	return ars
}

func TestParseArticles(t *testing.T) {
	ars := load()
	if len(ars) != 2 {
		t.Error("test failed: articles count do not match", len(ars))
	}
}

// not easy to test since we sort by update time
// func TestSort(t *testing.T) {
// 	ars := load()
// 	if ars[0].Title != "you dont know js" {
// 		t.Error("failed")
// 	}
// 	ars.Sort()

// 	if ars[0].Title != "programming in go" {
// 		t.Error("failed", ars[0].DateUpdated, ars[1].DateUpdated)
// 	}
// }

func TestPaingate(t *testing.T) {
	ars := load()
	res := ars.Paginate(1, 1)
	if len(res) != 1 {
		t.Error("failed")
	}
}
