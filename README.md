# GoPass

## GoPass is a CLI tool that help stores your passwords with security in mind.

This application encrypts, hashes, and stores passwords.
Everything is local to your computer! Nothing is stored on the internet

### Get Started

1. Ensure you have Golang installed on your system.

   a. You can install Golang [here](https://go.dev/doc/install)

2. Run `git clone https://github.com/DKagan07/go-pass.git` in your terminal to copy this repo onto your system.

3. Run `cd go-pass` to enter the root of the repo.

4. In the root of the repo, run `make`. This builds the binary and moves it to `/usr/local/bin`.

   a. If you want to move this binary to another custom path, you can run `make PREFIX=/custom/path`. Ensure this has a `bin` directory, and is in your `$PATH`.

5. You can now run `gopass <command>`!

### Commands (`gopass <command>`)

- **add**: Add a new password to the vault
- **clean**: Removes all storage of passwords
- **completion**: Generate the autocompletion script for the specified shell
- **delete**: Delete a specific item from your vault
- **get**: Get specific information from your vault by source name
- **help**: Help about any command
- **init**: Initialize all files and begins use of the app
- **list**: Lists all the sources of your login infos
- **login**: Login to the app
- **update**: Updates an entry in your vault with specific flags

### Workflow:

1. `gopass init` creates the files necessary and stores a master password.

   a. This counts as an initial login when creating files. For all future logins, you will need to `gopass login`.

2. Use the commands listed above to interact with the app and passwords.

3. `gopass login` takes the master password and unlocks use of the app.

   a. If the app is not in use for >30 minutes, you will need to login.
