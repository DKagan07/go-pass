package model

type VaultEntry struct {
	// Name is the source name for the login info
	Name string `json:"name"`
	// Username is the username for the login
	Username string `json:"username"`
	// Password is an encrypted password, encrypted with AES-256-GCM
	Password []byte `json:"password"`
	// Notes is a section that can be empty that the user can add extra notes
	// about needing to login.
	Notes string `json:"notes,omitempty"`
	// UpdatedAt is the timestamp when the entry was created, in milliseconds
	UpdatedAt int64 `json:"updated_at"`
}
