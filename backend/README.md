## Privacy Negotiator API

### Local Development

1. `make install`
2. `docker-compose up`

This will auto-reload when changes are made to Go files.

### Building

1. `make build`


### Database Backup from AWS

```
pg_dump -h <hostname> -p <port> -U <user> <db_name> > <YYYYMMDD>.sql
```
