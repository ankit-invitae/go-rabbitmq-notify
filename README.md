# RabbitMQ Notify

> Simple macos tray application for getting RabbitMQ updates

## Install/Run

### macOS

Download [.app](https://github.com/ankit-invitae/go-rabbitmq-notify/releases) file and run as any other application.

## RabbitMQ Config file

Config file with name `.rabbitmq.config` must be present at the `homedir`. Below is the sample,

```json
{
  "username": "********",
  "password": "********************************",
  "queues": [
    {
      "key": "sauron-event-queue",
      "display": "Sauron Event Queue",
      "virtualHost": "lis",
      "endpoint": "sauron-event-queue",
      "interval": 10
    },
    {
      "key": "sauron-internal",
      "display": "Sauron Internal",
      "virtualHost": "lis",
      "endpoint": "sauron-internal",
      "interval": 5
    }
  ]
}
```
