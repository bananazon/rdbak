# Raindrop bookmark backup

`rdbak` is a minimalistic command line tool that downloads a local backup of your [Raindrop.io](https://raindrop.io) and stores them in a YAML file. Note, this is an enhanced version of [rdbak](https://github.com/gugray/rdbak). I crerated this version because the original version lacked some things I wanted, e.g., proper error handling. However, I needed to give credit where it's due.

## Differences
* I'm using the [cobra](https://github.com/spf13/cobra) framework.
* I added a logging mechanism that logs to STDOUT/STDERR when run from the terminal and to `~/.config/rdbak/rdbak.log` when run via cron.
* Instead of using `panic(err)`, I properly return and process the errors.
* I chose YAML for the config and bookmarks files because stings in YAML don't have to be quoted and a double quote can mess with JSON in some cases.
* I also detect deleted bookmarks and remove them from the `bookmarks.yaml` file. This wasn't being done in the original.
* I only create a backup file if there are changes.
* When a backup is saved, the existing file is copied to `bookmarks-{timestamp}.yaml`.
* The `backup` command has a `--prune` option to delete `bookmarks-{timestamp}.yaml` files older than seven days. The retention period will soon be made into a flag."
* At runtime, rdbak verifies `~/.config/rdbak` exists, is a directory, and is writable. If it doesn't exist, rdbak attempts to create it with 0600 as the mode.

## Installation
* Clone this repository
* `cd <repo root>`
* `make install`
* `mkdir -p ~/.config/rdbak`
* `cp config.SAMPLE.yaml ~/.config/rdbak/config.yaml`
* Populate the `~/.config/rdbak/config.yaml` with your details
* `rdbak encrypt-password`
* `rdbak backup`

## crontab
I run it every five minutes
```
% crontab -l
#rdbak
*/5 * * * * /home/dummy/go/bin/rdbak backup --prune
```

## More Coming Soon!
I will populate more of this soon
