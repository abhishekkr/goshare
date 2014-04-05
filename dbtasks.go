package goshare

import (
  "strings"

  abkleveldb "github.com/abhishekkr/levigoNS/leveldb"
  levigoNS "github.com/abhishekkr/levigoNS"
  levigoTSDS "github.com/abhishekkr/levigoTSDS"
  golhashmap "github.com/abhishekkr/gol/golhashmap"
  goltime "github.com/abhishekkr/gol/goltime"
)


/* Get value of given key */
func GetVal(key string) string{
  return abkleveldb.GetVal(key, db)
}


/* Push a given set of Key-Val */
func PushKeyVal(key string, val string) bool{
  return abkleveldb.PushKeyVal(key, val, db)
}


/* Empty Val for a given Key */
func DelKey(key string) bool{
  return abkleveldb.DelKey(key, db)
}


/* Get value for all descendents of Namespace */
func GetValNS(key string) string{
  hashmap := levigoNS.ReadNSRecursive(key, db)
  return golhashmap.Hashmap_to_csv(hashmap)
}


/* Push a given Namespace-Key and its value */
func PushKeyValNS(key string, val string) bool{
  return levigoNS.PushNS(key, val, db)
}


/* Delete a Namespace Key and all its value */
func DelKeyNS(key string) bool{
  return levigoNS.DeleteNSRecursive(key, db)
}


/* Get value for the asked time-frame key, aah same NS */
func GetValTSDS(key string) string{
  return golhashmap.Hashmap_to_csv(levigoTSDS.ReadTSDS(key, db))
}


/* Push a key namespace-d with goltime.Timestamp  */
func PushKeyValTSDS(key string, val string, timestamp goltime.Timestamp) bool{
  if levigoTSDS.PushTSDS(key, val, timestamp.Time(), db) == "" { return false }
  return true
}


/* Push a key namespace-d with current time */
func PushKeyValNowTSDS(key string, val string) bool{
  if levigoTSDS.PushNowTSDS(key, val, db) == "" { return false }
  return true
}


/* Push a key-val TSDS pair from CSV message */
func PushCSVTSDS(csv_string string, timestamp goltime.Timestamp) bool{
  status := true
  _time := timestamp.Time()
  hashmap_key_value := golhashmap.Csv_to_hashmap(csv_string)
  for _key, _val := range hashmap_key_value {
    _val = strings.Replace(_val, "\n", " ", -1)
    if levigoTSDS.PushTSDS(_key, _val, _time, db) == "" { status = false }
  }
  return status
}


/* Delete all keys under given namespace, same as NS */
func DelKeyTSDS(key string) bool{
  current_val := levigoTSDS.DeleteTSDS(key, db)
  if len(current_val) > 0 { return false }
  return true
}


/* Get a value based on task-type */
func GetValTask(task_type string, key string) string{
  if task_type == "tsds" {
    return GetValTSDS(key)

  } else if task_type == "ns" {
    return GetValNS(key)

  }

  return GetVal(key)
}


/* Push a key-val based on task-type; except on with goltime.Timestamp */
func PushKeyValTask(task_type string, key string, value string) bool{
  if task_type == "tsds-now" {
    return PushKeyValNowTSDS(key, value)

  } else if task_type == "ns" {
    return PushKeyValNS(key, value)

  }

  return PushKeyVal(key, value)
}


/* Push a key-val based on task-type; except on with goltime.Timestamp */
func PushKeyValByType(task_type string, message_array []string) bool {
  _key := message_array[2]

  if task_type == "tsds" {
    timestamp := goltime.CreateTimestamp(message_array[3:9])
    _value := strings.Join(message_array[9:], " ")
    return PushKeyValTSDS(_key, _value, timestamp)

  } else if task_type == "tsds-csv" {
    timestamp := goltime.CreateTimestamp(message_array[2:8])
    csv_value := strings.Join(message_array[8:], "\n")
    return PushCSVTSDS(csv_value, timestamp)

  }

  _value := strings.Join(message_array[3:], " ")
  return PushKeyValTask(task_type, _key, _value)
}


/* Delete a key on task-type */
func DelKeyTask(task_type string, key string) bool{
  if task_type == "tsds" {
    return DelKeyTSDS(key)

  } else if task_type == "ns" {
    return DelKeyNS(key)

  }

  return DelKey(key)
}

