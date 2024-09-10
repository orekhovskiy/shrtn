package api

type Service interface {
	GetById(id string) (url string, err error)
	Save(url string) (id string)
}

type Handler struct {
	urlService Service
}

func NewHandler(urlService Service) *Handler {
	return &Handler{urlService: urlService}
}
