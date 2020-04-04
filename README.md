# Keybase Docker

A [Keybase](https://keybase.io/) integration to notify Docker Events via [webhookbot](https://keybase.io/webhookbot)


## Keybase webhook setup

+ Add **Webhook Bot** from list of Bots
+ Create a new webhook for sending messages into the current conversation. You must supply a name as well to identify the webhook. Example: `!webhook create alerts`
+ Get the new url to send webhooks

![Docker events on keybase](https://raw.githubusercontent.com/logocomune/webhookdocker/master/_img/keybase.png)


## Keybase docker

Run with docker
```shell
$ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock:ro \
logocomune/keybase-docker:latest --keybase-endpoint=https://bots.keybase.io/webhookbot/....
```


### Application options

### Application options

| flag | Environment |type | Default | |
| --- | --- | --- | --- | --- |
| --node-name | WD_NODE_NAME |String| | Node name. If empty use the hostname |
| --hide-node-name | WD_HIDE_NODE_NAME |Boolean| false | Node name is omitted |
| --docker-show-running | WD_DOCKER_SHOW_RUNNING | Boolean | false | Send running container to webhook |
| --docker-listen-container-events | WD_DOCKER_LISTEN_CONTAINER_EVENTS | Boolean | true | Listen for container events |
| --docker-listen-network-events | WD_DOCKER_LISTEN_NETWORK_EVENTS | Boolean | true | Listen for network events | 
| --docker-listen-volume-events |WD_DOCKER_LISTEN_VOLUME_EVENTS | Boolean | true | Listen for volume events | 
| --docker-listen-container-actions | WD_DOCKER_LISTEN_CONTAINER_ACTIONS| Strings separated by ; | attach;create;destroy;detach;die;kill;oom;pause;rename;restart;start;stop;unpause;update | Docker container events  |
| --docker-listen-network-actions | WD_DOCKER_LISTEN_NETWORK_ACTIONS | Strings separated by ; | create;connect;destroy;disconnect;remove | Docker network events |
| --docker-listen-volume-actions | WD_DOCKER_LISTEN_VOLUME_ACTIONS |  Strings separated by ; |  create;destroy;mount;unmount | Docker volume events |
| --docker-filter-container-name | WD_DOCKER_FILTER_CONTAINER_NAME | Regexp | |Filter events by container name (default all) | 
| --docker-filter-image-name | WD_DOCKER_FILTER_IMAGE_NAME | Regexp | |Filter events by image name (default all) | 
| --keybase-endpoint | WD_KEYBASE_ENDPOINT | String | |  Keybase endpoint for webhook | 


