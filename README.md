# Raindrop CLI

`raindrop` started as an enhanced version of [raindrop](https://github.com/gugray/raindrop) and evolved into something more. It now allows you to manage bookmarks, collections, and tags. I've also replaced the email/password auth with API Token auth.

## Installation

1. `git clone https://github.com/bananazon/raindrop.git`
2. `cd raindrop`
3. `make install`
4. `cp config.SAMPLE.yaml ~/.config/raindrop/config.yaml`
5. [Create a token](#create-a-token)
6. `raindrop encrypt-token`
7. Verify that there is now an `encryptedApiToken` field in `~/.config/raindrop/config.yaml`
8. Test with something like `raindrop collections list`

## Create a Token
1. Login to your [raindrop](https://raindrop.io) account
2. Go to the [integrations](https://app.raindrop.io/settings/integrations) page
3. Click `Create new app`
4. In the `Create new app` dialog, give your app an arbitrary name, check the accept checkbok, and then click `Create`
5. Click on your new app in the list
6. Click the `Create test token` link then click `OK` to confirm
7. The newly created token will be used for the `apiToken` field in `~/.config/raindrop/config.yaml`

## Features
There is a fair amount to cover so I'll only show top-level stuff. For sub-commands, do something like `raindrop bookmarks add --help`

### Bookmarks

```
% raindrop bookmarks
Manage bookmarks in your raindrop.io account

Usage:
  raindrop bookmarks [command]

Aliases:
  bookmarks, b

Available Commands:
  add         Add a new bookmark to your raindrop.io account
  backup      Back your raindrop.io bookmarks up to a YAML file
  list        List the bookmarks in your raindrop.io account
  remove      Remove an existing bookmark from your raindrop.io account
  update      Update an existing bookmark in your raindrop.io account

Flags:
  -h, --help   help for bookmarks

Use "raindrop bookmarks [command] --help" for more information about a command.
```

### Collections

```
% raindrop collections
Manage collections in your raindrop.io account

Usage:
  raindrop collections [command]

Aliases:
  collections, c

Available Commands:
  add         Add a new collection to your raindrop.io account
  backup      Back your raindrop.io collections up to a YAML file
  list        List the collections in your raindrop.io account
  remove      Remove an existing collection from your raindrop.io account
  sort        Sort your raindrop.io collections
  update      Update an existing collection in your raindrop.io account

Flags:
  -h, --help   help for collections

Use "raindrop collections [command] --help" for more information about a command.
```

### Tags

```
% raindrop tags
Manage tags in your raindrop.io account

Usage:
  raindrop tags [command]

Aliases:
  tags, t

Available Commands:
  list        List the tags in your raindrop.io account
  remove      Remove tags from your raindrop.io account
  rename      Rename a tag in your raindrop.io account

Flags:
  -h, --help   help for tags

Use "raindrop tags [command] --help" for more information about a command.
```
