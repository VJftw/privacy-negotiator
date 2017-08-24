# Privacy Negotiator
by [VJ Patel](https://vjpatel.me)

## Requirements

 * Docker CE (https://docker.com)
 * Docker Compose (https://docs.docker.com/compose/)
 * Make (https://www.gnu.org/software/make/)

## Structure
A single repository has been used to keep set up simple.

### Backend
A Go application that exposes an API and runs Worker services.
### Infrastructure
The configuration for deployment to AWS.

### Web App
An Angular 4 web application to interact with the API via a User Interface.

## Running

To get started with a built application locally:

1. Run
```
make build
```

2. Add the following to your `/etc/hosts` file:
```
127.0.0.1	alpha.privacymanager.social
```

3. Run
```
docker-compose up
```

4. Open your web browser and navigate to http://alpha.privacymanager.social:4200
