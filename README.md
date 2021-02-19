# Hukz

A simple & lightweight service implemented in go to add webhooks in your application. This service is intended to fire webhooks on different events that are registered.

### Concepts
* **Event**: A simple event string (eg. post.created)
* **Webhook**: Webhook registration object. It contains information about events, url (to which webhook will be fired), tags (key value list to filter webhooks)

### Usage
* Need to install **[nats.io](https://nats.io/)** for queuing service.
* Needs to install postgres database.
* On publishing event to nats.io server, huks service will subscribe to all events added in db and fire webhooks accordingly.
* One can add multiple events, the naming convention of events is based on your application usage. 
* Every request takes "X-User" header and to keep track of user who created resources (webhook, event etc.).
* Whenever a webhook is fired its response, tags, etc are logged in webhook_logs.

### Sample Config
```
DATABASE_HOST=postgres 
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=hukz 
DATABASE_PORT=5432 
DATABASE_SSL_MODE=disable

NATS_URL=nats://nats:4222

MODE=development
```

* Config file should be stored in project root folder with name config (ext can be yml, json, env)
* Environment variables can also be set for configuration parameters.