# go-pass

## GoPass is a CLI tool that help stores your passwords with security in mind.

This application encrypts, hashes, and stores passwords.
Everything is local to your computer! Nothing is stored on the internet

To enable it, `git clone` this repo, `cd` into it, and run `make` in the root of the repo!

#### TODO:

1. ~Add a notes section to model.VaultEntry for users to add misc info~
2. ~Look into starring the password ("\*") while being typed in, perhaps with delay?~
   ->~This seems to be complex, especially with cobra, so probably will not do.~
   ->~sudo`doesn't even have this, which makes me okay with this decision~
3. ~Implement Update functionality~
4. Implement 'init' functionality  
   a. ~Create config file~
   b. ~Ask for master password~
   c. ~Create vault here and other setup~
   d. ~Update other cmd functions to read vault name from cfg to properly read the correct vault~
   e. ~If init is already done, fail somehow?~
   f. Probably if config file doesn't exist, we remove all files and re-create them
5. Refactor for tests and write some tests for current implementation
   a. Need to implement the 'init' functionality before tests can be written for all the commands
6. Update all functionality to incorporate the master Password check

##### Current thoughts on workflow:

1. `gopass init` creates the file, user creates the master password (bcrypt)
2. `gopass login` takes the master password and unlocks use of the app
3. can then do all the CRUD stuff within the app now that the user is logged in
4. `gopass lock` will manually lock (read: logout) of the app, so the user can no longer use the app
5. Perhaps have a config file, also encrypted(?), to store config data,
   a. like last accessed, in MS, and the master password?
   b. Copmare the last logged in with 30 mins or something arbitrary to re-lock access
