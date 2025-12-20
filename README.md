![GoPass Logo](https://github.com/user-attachments/assets/df96f0ac-a1d5-4f61-9153-e245c8a5777c)

<div align="center">

# GoPass

**A secure, local-first password manager with CLI and TUI interfaces**

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)](https://github.com/DKagan07/go-pass)

</div>

---

## Overview

GoPass is a command-line password manager that stores your passwords encrypted on your local computer. It uses AES-256-GCM encryption, PBKDF2 key derivation, and OS keyring integration for secure password storage. All data stays local—no cloud services or internet required.

**Key Features:**
- Triple-layer security (OS keyring + master password + environment salt)
- Interactive TUI and traditional CLI interfaces
- Encrypted backups and restore
- Password generation and search
- Session timeout management
- Cross-platform (Linux, macOS, Windows)

---

## Quick Start

```bash
# Clone the repository
git clone https://github.com/DKagan07/go-pass.git
cd go-pass

# Set up environment variable (must be exactly 32 characters)
export SECRET_PASSWORD_KEY="your-32-character-secret-key-here"

# Build and install
make

# Initialize vault
gopass init

# Launch TUI
gopass -o
```

---

## Installation

### Requirements

- Go 1.23 or higher
- Linux, macOS, or Windows
- System keyring support (gnome-keyring, Keychain, or Credential Manager)

### Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/DKagan07/go-pass.git
   cd go-pass
   ```

2. **Configure environment variable:**
   
   Generate a 32-character secret key:
   ```bash
   openssl rand -base64 24 | head -c 32
   ```
   
   Add to your shell config (`~/.bashrc`, `~/.zshrc`, etc.):
   ```bash
   export SECRET_PASSWORD_KEY="your-32-character-secret-key-here"
   ```
   
   Reload your shell:
   ```bash
   source ~/.bashrc  # or ~/.zshrc
   ```

3. **Build and install:**
   ```bash
   make
   ```
   
   For custom installation path:
   ```bash
   make PREFIX=~/.local
   ```

4. **Initialize:**
   ```bash
   gopass init
   ```

---

## Usage

### TUI Mode

Launch the interactive terminal interface:
```bash
gopass -o
```

**Keyboard shortcuts:**
- `a` - Add entry
- `d` - Delete entry
- `u` - Update entry
- `Enter` - View entry
- `Tab` - Switch focus
- `Esc` - Close/cancel
- `Ctrl+C` - Exit

### CLI Commands

**Vault operations:**
```bash
gopass vault add                    # Add password
gopass vault list                   # List all passwords
gopass vault search <query>         # Search passwords
gopass vault get <name>             # Get specific password
gopass vault update                 # Update password
gopass vault delete                 # Delete password
gopass vault generate [--length N]  # Generate password
```

**Backup and restore:**
```bash
gopass vault backup                 # Create backup
gopass vault list --backup          # List backups
gopass vault restore                # Restore from backup
```

**Configuration:**
```bash
gopass config view                  # View settings
gopass config change-masterpass     # Change master password
gopass config update-timeout        # Update session timeout
```

**Maintenance:**
```bash
gopass login                        # Login after timeout
gopass clean                        # Remove all data
gopass help                         # Show help
```

---

## Security

GoPass uses a three-layer security model:

1. **Master Password** - Bcrypt hashed, required for all operations
2. **OS Keyring** - 32-byte random key stored in system keyring (hardware-backed where available)
3. **Environment Salt** - 32-character `SECRET_PASSWORD_KEY` for PBKDF2 derivation

**Encryption:** AES-256-GCM with PBKDF2 (100,000 iterations, SHA-256)

All three layers must be compromised to decrypt your vault. Data is authenticated to prevent tampering.

### File Locations

- Vault: `~/.local/gopass/pass.json` (encrypted)
- Config: `~/.config/gopass/gopass-cfg.json` (encrypted)
- Backups: `~/.local/gopass-backup/backup__<timestamp>.json` (encrypted)
- Keyring: System-dependent (OS-managed)

---

## Troubleshooting

**"Salt not appropriate length or not present"**
- Verify `SECRET_PASSWORD_KEY` is exactly 32 characters
- Check it's exported: `echo $SECRET_PASSWORD_KEY`
- Reload shell: `source ~/.bashrc`

**Keyring errors (Linux)**
- Install gnome-keyring: `sudo apt-get install gnome-keyring`
- Ensure daemon is running: `gnome-keyring-daemon --start`

**Permission denied during installation**
- Use `sudo` when prompted, or install to user directory: `make PREFIX=~/.local`

**Forgot master password**
- If you have a backup with the old password, you can restore after reinitializing
- Without a backup, passwords cannot be recovered (by design)

For more issues, see [GitHub Issues](https://github.com/DKagan07/go-pass/issues).

---

## Development

### Project Structure

```
go-pass/
├── cmd/          # CLI commands and TUI
├── model/        # Data models and keyring
├── crypt/        # Encryption/decryption
├── utils/        # File I/O and utilities
└── testutils/    # Testing helpers
```

### Building

```bash
go build -o gopass          # Build binary
go run main.go <command>    # Run without installing
make test                   # Run tests
```

### Contributing

Contributions welcome! Please:
1. Check existing issues first
2. Fork and create a feature branch
3. Run tests before submitting
4. Follow existing code style
5. Open a pull request

---

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

## Acknowledgments

Thanks to the maintainers of:
- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [rivo/tview](https://github.com/rivo/tview) - TUI library
- [zalando/go-keyring](https://github.com/zalando/go-keyring) - Keyring integration
- Go crypto libraries

---

<div align="center">

[Back to Top](#gopass)

</div>
