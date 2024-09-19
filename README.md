# WIP Book Manager
## Features
- Manage your ebook library
- Download books via IRC
- Send to Kindle

## Local Development
### MongoDB
If you run mongodb as a service you can ignore this. I chose not to and instead manually run it in the background when I am working on this project. To start it as a background process: 
```
mongod --config /usr/local/etc/mongod.conf --fork
```
### Start The Server
This is not a hot-reload server, whenever you make a code change you must restart the server yourself.
```
go run cmd/server.go
```