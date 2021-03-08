# Hukz

**Releasability:** [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=factly_hukz&metric=alert_status)](https://sonarcloud.io/dashboard?id=factly_hukz)  
**Reliability:** [![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=factly_hukz&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=factly_hukz) [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=factly_hukz&metric=bugs)](https://sonarcloud.io/dashboard?id=factly_hukz)  
**Security:** [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=factly_hukz&metric=security_rating)](https://sonarcloud.io/dashboard?id=factly_hukz) [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=factly_hukz&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=factly_hukz)  
**Maintainability:** [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=factly_hukz&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=factly_hukz) [![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=factly_hukz&metric=sqale_index)](https://sonarcloud.io/dashboard?id=factly_hukz) [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=factly_hukz&metric=code_smells)](https://sonarcloud.io/dashboard?id=factly_hukz)  
**Other:** [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=factly_hukz&metric=ncloc)](https://sonarcloud.io/dashboard?id=factly_hukz) [![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=factly_hukz&metric=duplicated_lines_density)](https://sonarcloud.io/dashboard?id=factly_hukz) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=factly_hukz&metric=coverage)](https://sonarcloud.io/dashboard?id=factly_hukz)  


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
* Environment variables can also be set for configuration parameters with `HUKZ_` prefix to above variables.