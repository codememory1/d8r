package postgres

// Options contains the settings required to connect to PostgreSQL.
type Options struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Database string
}
