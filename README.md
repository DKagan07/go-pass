![GoPass Logo](https://github.com/user-attachments/assets/df96f0ac-a1d5-4f61-9153-e245c8a5777c)

# GoPass

**A secure, local-first password manager with dual CLI/TUI interfaces**

GoPass is a modern command-line password manager built with Go that prioritizes security and privacy. All data is encrypted and stored locally on your machine—nothing ever touches the internet.

## Features

- ** Triple-Layer Security**: OS keyring integration + master password + environment salt
- ** Dual Interface**: Full-featured CLI and interactive TUI (Terminal User Interface)
- ** AES-256-GCM Encryption**: Authenticated encryption for all password storage
- ** PBKDF2 Key Derivation**: 100,000 iterations with SHA-256
- ** Encrypted Backups**: Create and restore from encrypted backup files
- ** Fast Search**: Real-time filtering and case-insensitive search
- ** Password Generator**: Cryptographically secure password generation (default 24 characters)
- ** Session Management**: Configurable timeout for enhanced security
- ** Local-Only Storage**: Zero network dependencies, complete privacy

## Security Architecture

GoPass implements a robust three-layer security model:

### Layer 1: Master Password Authentication
- User-provided master password
- Hashed with bcrypt before storage
- Required for all vault operations

### Layer 2: OS Keyring Integration
- 32-byte random key stored in system keyring
- Hardware-backed on supported platforms (Keychain on macOS, Secret Service on Linux, Credential Manager on Windows)
- Service: `gopass`, Account: `encryption_key`

### Layer 3: Environment-Based Salt
- `SECRET_PASSWORD_KEY` environment variable as a salt (32 bytes)
- Additional salt layer for PBKDF2 key derivation
- Combined with keyring key and master password

**Encryption Process:**
```
Base Key (from keyring) + Master Password → PBKDF2 (100k iterations, SHA-256) → AES-256-GCM Key
```

## Installation

### Requirements
- Go 1.23 or higher
- Linux, macOS, or Windows
- System keyring support (typically built-in)

### Quick Install

1. **Clone the repository:**
   ```bash
   git clone https://github.com/DKagan07/go-pass.git
   cd go-pass
   ```

2. **Set up environment variable:**

   Add this to your shell RC file (`~/.bashrc`, `~/.zshrc`, etc.):
   ```bash
   export SECRET_PASSWORD_KEY="your-32-character-secret-key-here"
   ```

   Generate a secure 32-character key [here](https://passwords-generator.org/32-character).

   **Note:** Escape special characters like `'`, `"`, `` ` ``, and `\` if needed.

3. **Build and install:**
   ```bash
   make
   ```

   This builds the binary and installs it to `/usr/local/bin`.

   **Custom installation path:**
   ```bash
   make PREFIX=/custom/path
   ```
   Ensure the path contains a `bin/` directory and is in your `$PATH`.

4. **Initialize GoPass:**
   ```bash
   gopass init
   ```

## Usage

### TUI Mode (Interactive Interface)

Launch the full-featured terminal UI:
```bash
gopass -o
```

**TUI Keyboard Shortcuts:**
- `a` - Add new entry
- `d` - Delete selected entry
- `u` - Update selected entry
- `Enter` - View entry details
- `Tab` - Switch between search bar and vault list
- `Esc` - Exit modals/cancel operations
- `/` - Focus search bar

### CLI Mode

#### Getting Started
```bash
# Initialize vault (first-time setup)
gopass init

# Login (required after timeout)
gopass login
```

#### Vault Commands

```bash
# Add a new password entry
gopass vault add

# List all entries
gopass vault list

# Search for entries (case-insensitive)
gopass vault search <query>

# Get a specific entry
gopass vault get <name>

# Update an existing entry
gopass vault update

# Delete an entry
gopass vault delete

# Generate a secure password
gopass vault generate [--length 24]

# Create an encrypted backup
gopass vault backup

# Restore from backup
gopass vault restore

# List available backups
gopass vault list --backup
```

#### Configuration Commands

```bash
# View current configuration
gopass config view

# Change master password
gopass config change-masterpass

# Update session timeout (in milliseconds)
gopass config update-timeout
```

#### Maintenance Commands

```bash
# Remove all data (vault + config + keyring entry)
gopass clean

# Display help
gopass help
gopass vault help
```

## File Locations

- **Vault:** `~/.local/gopass/pass.json` (encrypted)
- **Config:** `~/.config/gopass/gopass-cfg.json` (encrypted)
- **Backups:** `~/.local/gopass-backup/backup__<timestamp>.json`
- **Keyring:** System keyring (location varies by OS)

## Development

### Project Structure

```
go-pass/
├── cmd/                    # Command implementations
│   ├── vault/             # Vault-related commands
│   ├── config/            # Configuration commands
│   └── test_tview_root.go # TUI implementation
├── model/                 # Data models and keyring manager
├── crypt/                 # Encryption/decryption logic
├── utils/                 # File I/O and utilities
└── testutils/             # Testing utilities
```

### Running Tests

```bash
# Run all tests
make test

# Run tests for specific package
go test ./cmd/...
go test ./crypt/...
```

### Building from Source

```bash
# Build binary
go build -o gopass

# Run without installing
go run main.go <command>
go run main.go -o  # TUI mode
```

### Cleanup

```bash
# Remove vault and config files
make remove

# Uninstall binary
make uninstall
```

## Dependencies

- **CLI Framework:** [spf13/cobra](https://github.com/spf13/cobra)
- **TUI Framework:** [rivo/tview](https://github.com/rivo/tview) + [gdamore/tcell](https://github.com/gdamore/tcell)
- **Keyring:** [zalando/go-keyring](https://github.com/zalando/go-keyring)
- **Crypto:** `golang.org/x/crypto` (bcrypt, PBKDF2)
- **Testing:** [stretchr/testify](https://github.com/stretchr/testify)

## Security Best Practices

1. **Never share your `SECRET_PASSWORD_KEY`** - Treat it like a password
2. **Use a strong master password** - It's your first line of defense
3. **Regular backups** - Use `gopass vault backup` periodically
4. **Session timeout** - Default is 30 minutes; adjust based on your security needs
5. **Secure your backup files** - They contain encrypted vault data

## Troubleshooting

**"Salt not appropriate length or not present"**
- Ensure `SECRET_PASSWORD_KEY` is exactly 32 characters
- Verify it's exported in your current shell session

**Keyring errors on Linux**
- Ensure `gnome-keyring` or `seahorse` is installed
- For headless systems, consider using `secret-tool` manually

**Permission denied during installation**
- The Makefile uses `sudo` to install to `/usr/local/bin`
- Alternatively, use `make PREFIX=~/.local` for user-local installation

## License

See [LICENSE](LICENSE) file for details.

## Acknowledgments

Built with security and privacy in mind. Special thanks to the Go community and the maintainers of the excellent libraries this project depends on.
