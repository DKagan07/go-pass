package model

type VaultEntry struct {
	// Name is the source name for the login info
	Name string
	// Username is the username for the login
	Username string
	// Password is a hashed password that we want stored, hashed with bcrypt
	Password []byte
	// CreatedAt is the timestamp when the entry was created, in milliseconds
	CreatedAt int64
}
