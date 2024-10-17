package urlrepo

func (r *Repository) Ping() error {
	return r.db.Ping()
}
