package articles

type Articles []Article

func ParseArticles() *Article {
	return nil
}

func LoadArticles() *Article {
	return nil
}

func (ar *Articles) Persist() {
	//todo
}

func (*Articles) Sort() {
	//
}

func (*Articles) Paginate(page int32, pageSize int32) *Articles {
	return nil
}

func (*Articles) Update() {
	//
}
