package model

type VaultEntry struct {
	// Name is the source name for the login info
	Name string `json:"name"`
	// Username is the username for the login
	Username string `json:"username"`
	// Password is an encrypted password, encrypted with AES-256-GCM
	Password []byte `json:"password"`
	// CreatedAt is the timestamp when the entry was created, in milliseconds
	CreatedAt int64 `json:"created_at"`
}
