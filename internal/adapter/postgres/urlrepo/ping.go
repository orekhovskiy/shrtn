package urlrepo

func (r *PostgresURLRepository) Ping() error {
	return r.db.Ping()
}
