package urlrepo

const (
	createTableIfNotExists = `
	CREATE TABLE IF NOT EXISTS url_records (
		id SERIAL PRIMARY KEY,
		uuid TEXT NOT NULL,
		short_url TEXT UNIQUE NOT NULL,
		original_url TEXT NOT NULL
	);`
	isRecordExists = "SELECT EXISTS(SELECT 1 FROM url_records WHERE short_url=$1);"
	insertRecord   = "INSERT INTO url_records (uuid, short_url, original_url) VALUES ($1, $2, $3)"
	getRecordByID  = "SELECT original_url FROM url_records WHERE short_url=$1"
)
