# mongodb-changes-notifier

With this tool you'll be able to be notified of any data changes in MongoDB.
You can select what collections you want to watch, what change event types you interested in and generate notifications with the provided template.

## Features

- Watching for changes in your data
- Notification via Telegram Bot
- Filtering by [Change event type](https://docs.mongodb.com/manual/reference/change-events/)
- Watching for several collections at once
- Template engine with funcs provided by [sprig](https://github.com/Masterminds/sprig)

## Prerequirements

- [MongoDB replica set](https://docs.mongodb.com/manual/tutorial/deploy-replica-set/)
- [Webhook from CodeX notify bot](https://github.com/codex-bot/notify)

## Usage



## TODO
- [ ] Another notifications channels (Email, Slack, etc)
- [ ] Watching for collections by regexp
- [ ] Watching the entire deployment ([link](https://docs.mongodb.com/manual/changeStreams/#watch-collection-database-deployment))