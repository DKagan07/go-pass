# GoPass - Password Manager

## Project Overview

**GoPass** is a secure, local-first CLI password manager written in Go with dual CLI/TUI interfaces. It provides encrypted password storage using OS keyring integration, AES-256-GCM encryption, and PBKDF2 key derivation.

**Key Features**:
- Dual interface: CLI (Cobra) and TUI (tview)
- AES-256-GCM authenticated encryption
- OS keyring integration for secure key storage
- PBKDF2 key derivation (100,000 iterations)
- Encrypted backups and restore
- Password generation
- Search and filtering capabilities
- Local-only storage (no network)

## Architecture

### Directory Structure

```
go-pass/
├── main.go              # Entry point, validates SECRET_PASSWORD_KEY (32 bytes)
├── cmd/                 # Command-line interface
│   ├── root.go          # Base command, TUI switch (-o flag)
│   ├── init.go          # Initialize vault/config
│   ├── login.go         # Authenticate with master password
│   ├── clean.go         # Remove all data
│   ├── cptf.go          # Create plaintext file export (migration tool)
│   ├── upload_cptf.go   # Import plaintext file (migration tool)
│   ├── tui/             # Terminal UI implementation
│   │   ├── run.go       # TUI entry point and main screen
│   │   ├── vault.go     # Vault display logic
│   │   ├── vault_add.go    # Add entry modal
│   │   ├── vault_delete.go # Delete confirmation modal
│   │   └── vault_update.go # Update entry modal
│   ├── vault/           # Password vault commands
│   │   ├── add.go       # Add new entry
│   │   ├── backup.go    # Create encrypted backup
│   │   ├── delete.go    # Remove entry
│   │   ├── generate.go  # Generate random password (24 chars)
│   │   ├── get.go       # Retrieve entry
│   │   ├── list.go      # List all entries
│   │   ├── restore.go   # Restore from backup
│   │   ├── search.go    # Case-insensitive search
│   │   └── update.go    # Modify existing entry
│   └── config/          # Configuration commands
│       ├── change-masterpass.go  # Update master password
│       ├── update-timeout.go     # Adjust session timeout
│       └── view.go              # Display current config
├── model/               # Data models
│   ├── model.go         # VaultEntry, Config structs
│   └── keyring.go       # OS keyring integration (MasterAESKeyManager)
├── crypt/               # Encryption/decryption wrappers
│   ├── crypt.go         # Constants (NONCE_SIZE, KEY_SIZE, NUM_ITERATIONS)
│   ├── encrypt.go       # Encrypt passwords/vault/config
│   └── decrypt.go       # Decrypt passwords/vault/config (supports legacy format)
└── utils/               # Utilities
    ├── setup.go         # File I/O, vault/config creation
    └── input.go         # User input (text, password, confirmations)
```

### Technology Stack

**Dependencies**:
- **CLI**: `github.com/spf13/cobra` - Command-line argument parsing
- **TUI**: `github.com/rivo/tview` + `github.com/gdamore/tcell` - Terminal UI
- **Keyring**: `github.com/zalando/go-keyring` - OS-level secure storage
- **Encryption**: `golang.org/x/crypto` (AES-GCM, bcrypt, PBKDF2)
- **Terminal**: `golang.org/x/term` - Password input handling
- **Testing**: `github.com/stretchr/testify/assert`

**Go Version**: 1.23

## Data Models

### VaultEntry
```go
type VaultEntry struct {
    Name      string    // Service/site name
    Username  string
    Password  []byte    // AES-256-GCM encrypted
    Notes     string
    UpdatedAt int64     // Unix milliseconds
}
```

### Config
```go
type Config struct {
    MasterPassword []byte              // bcrypt hashed
    VaultName      string               // Default: "pass.json"
    LastVisited    int64               // Unix milliseconds
    Timeout        int64               // Default: 30 minutes (1800000 ms)
    Keychain       MasterAESKeyManager  // OS keyring manager
}
```

### MasterAESKeyManager
```go
type MasterAESKeyManager struct {
    Masterpassword string  // Master password for key derivation
}
```

**Key Methods**:
- `InitializeKeychain()` - Generates and stores 32-byte random key in OS keyring
- `GetEncryptionKey()` - Derives AES key using PBKDF2 (keyring_key + master_password)
- `Encrypt(plaintext []byte)` - AES-256-GCM encryption
- `Decrypt(ciphertext []byte)` - AES-256-GCM decryption
- `GenerateNonce()` - Creates 12-byte cryptographic nonce
- `GetSalt()` - Retrieves SECRET_PASSWORD_KEY from environment

The `crypt/` package provides wrapper functions that delegate to `MasterAESKeyManager` methods for backward compatibility.

## Security Model

### Two-Layer Security Architecture

**Layer 1: Authentication**
- Master password provided by user
- Hashed with bcrypt and stored in config
- Used to validate user access
- Checked on every vault operation

**Layer 2: Encryption**
1. **Base Key**: 32-byte random key stored in OS keyring
   - Service: "gopass"
   - Account: "encryption_key"
   - Generated once during initialization
   - Hardware-backed on supported platforms

2. **Key Derivation**: PBKDF2
   - Input: Base key (from keyring) + Master password
   - Salt: SECRET_PASSWORD_KEY environment variable (32 bytes)
   - Iterations: 100,000
   - Hash: SHA256
   - Output: 32-byte AES encryption key

3. **Encryption**: AES-256-GCM
   - Authenticated encryption (prevents tampering)
   - Random 12-byte nonce per operation
   - Format: `base64(nonce || ciphertext)`

**Security Benefits**:
- OS keyring provides hardware-backed security and handles authentication
- Master password required even with keyring access
- PBKDF2 slows brute-force attacks
- Environment variable adds additional secret layer
- GCM mode provides authentication (tamper detection)

### Session Management
- Infrastructure exists for timeout enforcement (default: 30 minutes)
- `LastVisited` timestamp tracked in config
- Timeout is **not actively enforced** in vault operations since OS keyring already requires authentication
- TUI enforces timeout before displaying vault on launch

## Migration Tools

The `cptf` (create plaintext file) and `upload` commands exist to support migration from legacy encryption formats to the current keyring-based implementation:

1. **Export**: `gopass cptf` - Decrypts vault using legacy format and exports to plaintext file named `out`
2. **Import**: `gopass upload` - Reads plaintext file `out` and creates new vault with current keyring-based encryption

These tools are intended for one-time migration and should be used with caution as they create temporary plaintext files.

## File Storage

### Locations

**Vault File**:
- Path: `~/.local/gopass/pass.json`
- Format: Base64-encoded AES-GCM ciphertext
- Plaintext: JSON array of VaultEntry objects
- Permissions: 0600 (owner read/write only)

**Config File**:
- Path: `~/.config/gopass/gopass-cfg.json`
- Format: Base64-encoded AES-GCM ciphertext
- Contains: Master password hash, vault name, timeout, last login
- Permissions: 0600 (owner read/write only)

**Backups**:
- Path: `~/.local/gopass-backup/`
- Format: Same as vault (encrypted)
- Naming: `backup__<timestamp>.json`
- Permissions: 0600 (owner read/write only)

**OS Keyring**:
- Service: "gopass"
- Account: "encryption_key"
- Content: Base64-encoded 32-byte random key

## TUI Interface

### Implementation (`cmd/tui/`)

**Framework**: tview with tcell backend

**Features**:
- Login screen with master password input
- Search bar (real-time vault filtering)
- Vault list (alphabetically sorted)
- Modal dialogs for all operations (add, update, delete, view)

**Keyboard Shortcuts**:
- `a` - Add new entry
- `d` - Delete selected entry
- `u` - Update selected entry
- `Enter` - View entry details
- `Tab` - Switch focus between search and list
- `Esc` - Exit modals/cancel operations
- `/` - Focus search bar

**UI Components**:
- InputField - Text/password input
- List - Vault entries display
- Modal - Dialogs for add/update/delete/view
- Pages - Screen management (login, main, modals)

**Session Handling**:
- Enforces timeout check on TUI launch
- Returns to login screen on timeout
- Clears sensitive data from memory

## Command Reference

### Core Commands
```bash
gopass init                    # Initialize vault and config
gopass login                   # Authenticate with master password
gopass clean                   # Remove all data (vault + config)
gopass -o                      # Launch TUI mode
gopass cptf                    # Export vault to plaintext file (migration tool)
gopass upload                  # Import plaintext file to new vault (migration tool)
```

### Vault Commands
```bash
gopass vault add               # Add new password entry
gopass vault backup            # Create encrypted backup
gopass vault delete            # Remove entry from vault
gopass vault generate          # Generate secure password (24 chars)
gopass vault get <name>        # Retrieve specific entry
gopass vault list [--backup]   # List all entries or backups
gopass vault restore           # Restore from backup
gopass vault search <query>    # Case-insensitive search
gopass vault update            # Modify existing entry
```

### Config Commands
```bash
gopass config change-masterpass  # Update master password
gopass config update-timeout     # Adjust session timeout
gopass config view              # Display current config
```

## Development Guidelines

### Code Patterns

**Command Handlers**:
- Follow naming: `<Name>CmdHandler()`
- Separate business logic from command setup
- Return errors, don't panic (except critical crypto failures)
- Use `log.Fatalf()` for unrecoverable errors

**Error Handling**:
- Wrap errors with context: `fmt.Errorf("context: %w", err)`
- Check all errors immediately
- User-facing errors should be descriptive
- Never log or display passwords

**Testing**:
- File naming: `*_test.go`
- Use testify assertions: `assert.NoError(t, err)`
- Separate test data from production (test vault/config names)
- Test utilities in `utils/test_utils.go`

**Security Best Practices**:
- Validate input lengths and formats
- Use `crypto/rand` for random generation
- Clear sensitive data from memory when done
- Use 0600 file permissions for sensitive data
- Use confirmation prompts for destructive operations

### Environment Setup

**Required Environment Variable**:
```bash
export SECRET_PASSWORD_KEY="<32-byte-hex-string>"
```

**Validation**: main.go validates this is exactly 32 bytes on startup

**Running**:
```bash
go run main.go <command>     # CLI mode
go run main.go -o            # TUI mode
```

**Building**:
```bash
go build -o gopass
```

## Project Strengths

1. **Well-Architected**: Clear separation of concerns across packages
2. **Comprehensive**: Full-featured password manager (CRUD, backup, search, generate)
3. **Secure Design**: Multiple layers of encryption, authenticated encryption
4. **Dual Interface**: Both CLI and TUI for different use cases
5. **Well-Tested**: Good test coverage for critical paths
6. **Local-First**: Privacy-focused, no network dependencies
7. **Good Documentation**: Inline comments and help text

<!-- ## Future Considerations -->
<!---->
<!-- **Security Enhancements**: -->
<!-- - Implement file locking for concurrent access -->
<!-- - Reduce file permissions to 0600 -->
<!-- - Add memory zeroing for sensitive data -->
<!-- - Consider removing environment variable requirement -->
<!---->
<!-- **Feature Additions**: -->
<!-- - Password strength analysis -->
<!-- - Password history/audit log -->
<!-- - Import/export formats (CSV, JSON) -->
<!-- - Password expiration/rotation reminders -->
<!-- - Multi-factor authentication -->
<!---->
<!-- **Code Quality**: -->
<!-- - Complete the keyring migration -->
<!-- - Remove commented-out code -->
<!-- - Add integration tests -->
<!-- - Implement atomic file writes with rollback -->
<!-- - Add benchmark tests for crypto operations -->
<!---->
<!-- --- -->
<!---->
<!-- **Last Updated**: 2025-11-02 -->
<!-- **Version**: In Development (pre-1.0) -->
<!-- **License**: Check LICENSE file in repository -->
