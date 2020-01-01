# chatr
This is a simple project built using a go backend and web (cljs) frontend.

## Development mode
All commands must be run in the project directory. Make sure you have go and lein + clojure installed.

### Build application:
Build the server.
```
make
```

Install webapp.
```
make app
```

### Run application:
Start the server.
```
bin/chatr
```
Start the webapp.
```
make start-app
```

## Run tests:
Server tests.
```
make test
```
Webapp tests.
Install karma and headless chrome
```
npm install -g karma-cli
```
Run the tests
```
lein karma
```