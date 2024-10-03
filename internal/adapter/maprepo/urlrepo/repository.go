package urlrepo

type Repository struct {
	urlMapping map[string]string
}

func NewRepository() Repository {
	urlMapping := make(map[string]string)
	return Repository{
		urlMapping: urlMapping,
	}
}
