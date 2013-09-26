## Go Share

#### Go Share any data among the nodes. Over HTTP (0MQ to come).

Tryout:
```Shell
 go run gohttp.go -dbpath=/tmp/GOTSDB
```
By default it runs at port 9797, make it run on another port using
```Shell
 go run gohttp.go -dbpath=/tmp/GOTSDB -port=8080
```

Now visit the the link asked by it and get the help page.

```ASCII
                        __
   ____ _____     _____/ /_  ____ _________
  / __ `/ __ \   / ___/ __ \/ __ `/ ___/ _ \
 / /_/ / /_/ /  (__  ) / / / /_/ / /  /  __/
 \__, /\____/  /____/_/ /_/\__,_/_/   \___/
/____/

```

##### Dependency
* [go lang](http://golang.org/doc/install) (obviously, the heart and soul of the app)
* [leveldb](http://en.wikipedia.org/wiki/LevelDB) (we are using for datastore, it's awesome)
* [levigo](https://github.com/jmhodges/levigo/blob/master/README.md) (the go library utilized to access leveldb)
