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
