# Raindrop bookmark backup

`rdbak` is a minimalistic command line tool that downloads a local backup of your [Raindrop.io](https://raindrop.io) and stores them in a YAML file. Note, this is an enhanced version of [rdbak](https://github.com/gugray/rdbak). I crerated this version because the original version lacked some things I wanted, e.g., proper error handling. However, I needed to give credit where it's due.

## Differences
* Added a logging mechanism.
* Instead of using `panic(err)`, I properly return and process the errors.
* The config file and backup files are stored as YAML instead of JSON.
* I'm using the [cobra](https://github.com/spf13/cobra) framework.
* When a backup is saved, the existing file is copied to a timestamp'd version.
* Optional `--prune` flag to delete backup files older than 7 days.
* Other small improvements.

## Installation
* Clone this repository
* `cd <repo root>`
* `make install`
* `mkdir -p ~/.config/rdbak`
* `cp config.SAMPLE.yaml ~/.config/rdbak/config.yaml`
* Populate the `~/.config/rdbak/config.yaml` with your details
* `rdbak encrypt-password`
* `rdbak backup`

## More Coming Soon!
I will populate more of this soon
