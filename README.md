![Docker Automated build](https://img.shields.io/docker/automated/jrottenberg/ffmpeg.svg)  
![version](https://img.shields.io/badge/Version-0.2.0-brightgreen.svg)  

<img src="saltgopher.png" width="250">

# Salt Gopher
Simple chatbot written in Go to help you manage your SaltStack infrastructure from slack.

Current version: `0.2.0` available features:

* Basic help 
* Simple role management  
* Executing jobs on salt-minions

Currently in development or planned:
* Complex role management
* Full Access Control List  
  * Same principle as SaltStack itself have
* Aliases
* Support for runners and wheels
* Support for listing jobs and minions
* and more...


### Install SaltGopher to your slack workspace

<a href="https://slack.com/oauth/authorize?client_id=418342024391.489186556466&scope=bot,chat:write:bot,incoming-webhook,files:write:user,users.profile:read"><img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcset="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" /></a>

### Configure your SaltGopher

Edit config/config.json with your informations:
```json
{
    "salt": {
        "url": "your.saltmaster.url", // Hostname or IP where your salt-api are running
        "port": 8000,                 // salt-api publish port
        "user": "saltgopher",         // user used to authentificate against salt-api
        "password": "*******",        // password for user used to authentificate against salt-api
        "eauth": "pam",               // salt-api auth module. Recommanded to use pam
        "SSLSkipVerify": false        // set to True if you are not using valid SSL certificate
    },
    "slack": {
        "token": "xoxb-000000000-000000000-AAaaaa000aaa00AAA000"
    }
}
```
`SaltGopher is supporting only HTTPS protocol to communicate with salt-api!`  
You can use self signed SSL certs, just set SSLSkipVerify to true to avoid errors.

### Run your SaltGopher instance
Easiest way to run your SaltGopher is to use official docker image.

Simply run:
```bash
$ docker run -d --name saltgopher simianlabs/saltgopher
```

For easier configuration mount /app/config to your host filesystem.
```bash
$ docker run -d --name saltgopher -v /var/saltgopher:/app/config simianlabs/saltgopher
```

You can build your own instance if you have Go environment available adn use systemd or superviso to run your instance.
```bash
$ go build -o saltgopher
$ ./saltgopher
```

### SaltStack execution
SaltGopher listen for standard Salt command syntax:
```
salt “target” module.function argument {kwargs}
```
<img src="pics/salttest.png">

### Roles

Current version have only simple role management.
To be able execute salt commands user have to have `saltadmin` role.

Roles are stored in config/roles.json
```json
{
  "users": [
    {
      "id": "SLACKUID",
      "roles": [
        "gopheradmin",
        "saltadmin"
      ]
    }
  ]
}
```
#### Setting new role
There are two ways to add new role to user:
1. Edit config/roles.json and add role manually
2. Use bot command **set role**  
<img src="pics/setrole.png">

`Roles are persistent and don't require SaltGopher instance restart to load new roles`  
SaltGopher will add user if don't exist in file when adding role with **set role** command.


#### Pre-builded packages

Not available at the moment.

### Maintainer
Dev: Stefan Monko || smonko@simianlabs.io  
`Simian Labs` (https://github.com/simianlabs)  
http://simianlabs.io || sl@simianlabs.io
