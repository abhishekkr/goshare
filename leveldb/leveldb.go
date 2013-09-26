package abkleveldb

import (
  "fmt"

  "github.com/jmhodges/levigo"
)

func boohoo(errstring string, rise bool){
  fmt.Println(errstring)
  if rise == true{ panic(errstring) }
}

func CreateDB(dbname string) (*levigo.DB) {
  opts := levigo.NewOptions()
  opts.SetCache(levigo.NewLRUCache(1<<10))
  opts.SetCreateIfMissing(true)
  db, err := levigo.Open(dbname, opts)
  if err != nil { boohoo("DB " + dbname + " Creation failed.", true) }
  return db
}

/* Push KeyVal */
func PushKeyVal(key string, val string, db *levigo.DB) bool{
  writer := levigo.NewWriteOptions()
  defer writer.Close()

  keyname := []byte(key)
  value := []byte(val)
  fmt.Printf("Writing for %s = %s\n", key, val)
  err := db.Put(writer, keyname, value)
  if err != nil {
    boohoo("Key " + key + " insertion failed. It's value was " + val, false)
    return false
  }
  return true
}


/* Get Key */
func GetValues(key string, db *levigo.DB) string {
  reader := levigo.NewReadOptions()
  defer reader.Close()

  data, err := db.Get(reader, []byte(key))
  if err != nil { boohoo("Key " + key + " query failed.", false) }
  return string(data)
}
