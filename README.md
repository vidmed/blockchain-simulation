# This is a small blockchain simulation.

### Install app and dependencies
1. Install go dep if it is not installed `go get -u github.com/golang/dep/cmd/dep`
2. Ensure dependencies `dep ensure`
3. Build and run app `go build`

Or just run `make` and everything wil be done automatically

Command `./blockchain-simulation` - runs app with default config `config.toml`


### Testing app
To test you can call ab (apache benchmark):

`ab -c 8 -n 10000 "http://localhost:8080/tx?key=sd&value=sdfsd"`

Or just open `http://localhost:8080/tx?key=sd&value=sdfsd` in your browser.
You should see HTTP 200 Ok as a response. That means that transaction accepted.

### Stopping app
To stop application press `CTRL+C`. This will gracefully stop the server and flush in-memory block on a disk.

### Config example
```
[main]
     ListenStr  = "0.0.0.0:8080"
     LogLevel = 5 # panic = 0, fatal = 1, error = 2, warning = 3, info = 4, debug = 5
     FlushPeriod = 10 # seconds to store `transactions` in memory
     MaxTransactions = 2 # Max amount of transactions in block. When the capacity will be reached - the block will be flushed on disk
     FlushFile = "blocks.json" # file name for storing `transactions` on a disk```
