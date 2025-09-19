![Image](https://github.com/user-attachments/assets/df96f0ac-a1d5-4f61-9153-e245c8a5777c)

# GoPass

## GoPass is a CLI tool that help stores your passwords with security in mind.

This application encrypts, hashes, and stores passwords.
Everything is local to your computer! Nothing is stored on the internet

### Features

- Secure storage of passwords locally
- Encryption of passwords with AES-256-GCM algorithm
- Hashing of passwords with bcrypt
- Encrypted Backups
- Restores form backups
- Searching through your vault
- Secure password generation

### Get Started

1. Ensure you have Golang installed on your system.

   a. You can install Golang [here](https://go.dev/doc/install)

2. You want to have an environment variable, called `SECRET_PASSWORD_KEY` in your rc file (ex. `~/.zshrc` or `~/.bashrc`). It needs to be 32 bytes long.

   a. You can create one [here](https://passwords-generator.org/32-character), for example. NOTE: You might need to escape certain characters, like single quote, double quote, backticks, and backslashes

3. Run `git clone https://github.com/DKagan07/go-pass.git` in your terminal to copy this repo onto your system.

4. Run `cd go-pass` to enter the root of the repo.

5. In the root of the repo, run `make`. This builds the binary and moves it to `/usr/local/bin`.

   a. If you want to move this binary to another custom path, you can run `make PREFIX=/custom/path`. Ensure this has a `bin` directory, and is in your `$PATH`.

6. You can now run `gopass <command> <sub_command>`!

### Commands

Commands are divided into different categories. The main category are 'vault' related commands

#### Vault Commands

- **add**: Add a new password to the vault
- **backup**: Backups up your vault
- **delete**: Delete a specific item from your vault
- **generate**: Generates a secure password with any length (default is 24 characters)
- **get**: Get specific information from your vault by source name
- **list**: Lists all the sources of your login infos
- **restore**: Restores your vault from a chosen backup
- **search**: Insensitive-case search for sources in your vault. Will return a list of all matches
- **update**: Updates an entry in your vault with specific flags

#### Config Commands

- **change_timeout**: Change the timeout of your session -> the longer the time, the longer time to pass before you need to log in
- **change-masterpass**: Change the master password, used for logging in
- **view**: Views the current state of your config without showing the Master Password

#### General Commands

- **clean**: Removes all storage of passwords
- **completion**: Generate the autocompletion script for the specified shell
- **init**: Initialize all files and begins use of the app
- **login**: Login to the app
- **help**: Help for all commands and sub_commands

### Workflow:

1. `gopass init` creates the files necessary and stores a master password.

   a. This counts as an initial login when creating files. For all future logins, you will need to `gopass login`.

2. Use the commands listed above to interact with the app and passwords.

3. `gopass login` takes the master password and unlocks use of the app.

   a. If the app is not in use for your timeout time (default is 30 mins), you will need to login.
