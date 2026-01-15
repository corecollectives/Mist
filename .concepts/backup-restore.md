# feature

we'll support backup and restore, on the same system and also on different systems, you can backup your
current state of mist into something like `bak.mist` or something, will be decided later, using the cli as well as ui (you should be able to download it through the ui), and on the same system or let's say a different system, assuming mist is already installed using the install script, `mist-cli` is also installed with it user should be able to do something like `mist-cli restore /path/to/bak.mist`, and this command should restore the backed up state, exactly as it was, with same database (almost), all the config files, logs, and all the apps up and running (atleast in the queue, and ready to be build and deployed).


## considerations

1. we can't just replace the database file `mist.db`:

let's say the backup was made at version `v1.0.2` and the restoration was done with the mist of version `v1.0.8` installed it will fuck up:
  - the versioning
  - new migrations (if any)

so we need to iterate through the old db from `bak.mist` and push the data into the new db present at `/var/lib/mist/mist.db`

2. deployments can't be replayed normally

in current implementation we deploy from the latest commit, but during the restoration we can't do that,
becuase let's say in between backup and restore i pushed 4 new commits to the project, then if we deploy the latest commit it won't truly restore the original state
