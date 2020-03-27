# Keybase Docker

A [Keybase](https://keybase.io/) integration to notify Docker Events via [webhookbot](https://keybase.io/webhookbot)


## Keybase webhook setup

+ Add **Webhook Bot** from list of Bots
+ Create a new webhook for sending messages into the current conversation. You must supply a name as well to identify the webhook. Example: `!webhook create alerts`
+ Get the new url to send webhooks

## Keybase docker

Run with docker
```shell
$ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock:ro \
logocomune/keybase-docker:1.0.0 --keybase-endpoint=https://bots.keybase.io/webhookbot/....
```


### Application options

| flag | Environment |type | Default | |
| --- | --- | --- | --- | --- |
| --node-name | KD_NODE_NAME |String| | Node name. If empty use the hostname |
| --docker-show-running | KD_DOCKER_SHOW_RUNNING | Boolean | false | Send running container to webhook |
| --docker-listen-container-events | KD_DOCKER_LISTEN_CONTAINER_EVENTS | Boolean | true | Listen for container events |
| --docker-listen-network-events | KD_DOCKER_LISTEN_NETWORK_EVENTS | Boolean | true | Listen for network events | 
| --docker-listen-volume-events | KD_DOCKER_LISTEN_VOLUME_EVENTS | Boolean | true | Listen for volume events | 
| --keybase-endpoint/ | KD_KEYBASE_ENDPOINT | String |  | Keybase endpoint for webhook | 

