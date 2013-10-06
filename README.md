## Go Share

```ASCII
                        __
   ____ _____     _____/ /_  ____ _________
  / __ `/ __ \   / ___/ __ \/ __ `/ ___/ _ \
 / /_/ / /_/ /  (__  ) / / / /_/ / /  /  __/
 \__, /\____/  /____/_/ /_/\__,_/_/   \___/
/____/

```

#### Go Share any data among the nodes. Over HTTP or ZeroMQ.
##### GOShare eases up communication over HTTP GET param based interaction.
##### OR ZeroMQ REQ/REP based synchronous communication model.

Tryout:

#### Over HTTP

```Shell
 go run gohttp.go -dbpath=/tmp/GOTSDB
```
By default it runs at port 9797, make it run on another port using
```Shell
 go run gohttp.go -dbpath=/tmp/GOTSDB -port=8080
```

```ASCII
  Dummy Client Using It

  * go run zxtra/gohttp_client.go

  for custom Port: 8080

  * go run zxtra/gohttp_client.go -port=8080
```

#### Over ZeroMQ

```Shell
 go run go0mq.go -dbpath=/tmp/GOTSDB
```
By default Binds at Ports: 9797, 9898, to opt for ports of choice 8000, 8080
```Shell
 go run go0mq.go -dbpath=/tmp/GOTSDB -req-port=8000 -rep-port=8080
```

```ASCII
  Dummy Client Using It

  * go run zxtra/go0mq_client.go

  for custom Port: 8080

  * go run zxtra/go0mq_client.go -req-port=8000 -rep-port=8080
```

Now visit the the link asked by it and get the help page.

##### Dependency
* [go lang](http://golang.org/doc/install) (obviously, the heart and soul of the app)
* [leveldb](http://en.wikipedia.org/wiki/LevelDB) (we are using for datastore, it's awesome)
* [levigo](https://github.com/jmhodges/levigo/blob/master/README.md) (the go library utilized to access leveldb)
* [zeroMQ](http://zeromq.org/) (the supercharged Sockets giving REQuest/REPly power)
* [gozmq](https://github.com/alecthomas/gozmq) GoLang ZeroMQ Bindings used here
