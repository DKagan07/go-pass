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

type Config struct {
	// MasterPassword is a bcrypt-hashed password that the user will need to
	// input in to use the app.
	MasterPassword []byte `json:"master_password"`
	// VaultName is the name of the vault file in the FS. There is a default
	// value of "pass.json".
	VaultName string `json:"vault_name"`
	// LastVisited is the time in UnixMilli when the user last used the app. The
	// idea is that the user will have to re-input the password in after 30 mins
	// of 'inactivity', for security.
	LastVisited int64 `json:"last_visited"`
	// Timeout is the number of minutes that the user will have to re-input the
	// password after.
	Timeout int64 `json:"timeout"`
}

type UserInput struct {
	// Username is the username that is obtained from the user.
	Username string
	// Password is the encrypted password that we get from the user and then
	// encrypt.
	Password []byte
	// Notes is the single string of notes that we get from the user.
	Notes string
}
