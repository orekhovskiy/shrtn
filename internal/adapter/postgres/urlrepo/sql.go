package urlrepo

const (
	createTableIfNotExists = `
	CREATE TABLE IF NOT EXISTS url_records (
		id SERIAL PRIMARY KEY,
		uuid TEXT NOT NULL,
		short_url TEXT UNIQUE NOT NULL,
		original_url TEXT UNIQUE NOT NULL,
		user_id TEXT NOT NULL
	);`
	createIndexIfNotExists = "CREATE INDEX IF NOT EXISTS idx_user_id ON url_records (user_id);"
	isRecordExists         = "SELECT short_url FROM url_records WHERE original_url = $1"
	insertRecord           = "INSERT INTO url_records (uuid, short_url, original_url, user_id) VALUES ($1, $2, $3, $4)"
	getRecordByID          = "SELECT original_url FROM url_records WHERE short_url=$1"
	getRecordsByUserID     = "SELECT short_url, original_url FROM url_records WHERE user_id=$1"
)
