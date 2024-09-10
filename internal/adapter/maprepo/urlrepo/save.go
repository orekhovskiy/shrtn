package urlrepo

func (r Repository) Save(id string, url string) {
	r.urlMapping[id] = url
}
