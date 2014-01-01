package goshare

import (
  "github.com/jmhodges/levigo"

  abkleveldb "github.com/abhishekkr/levigoNS/leveldb"
  levigoNS "github.com/abhishekkr/levigoNS"
  levigoTSDS "github.com/abhishekkr/levigoTSDS"
)


func GetVal(key string) string{
  return abkleveldb.GetVal(key, db)
}


func PushKeyVal(key string, val string) bool{
  return abkleveldb.PushKeyVal(key, val, db)
}


func DelKey(key string) bool{
  return abkleveldb.DelKey(key, db)
}
