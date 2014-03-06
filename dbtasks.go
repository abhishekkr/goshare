package goshare

import (
  "strings"
  "strconv"
  "time"

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


func hmap_to_csv(hmap levigoNS.HashMap) string{
  csv := ""
  for  key, value := range hmap {
    csv += key + "," + value + "\n"
  }
  return csv
}


func GetValNS(key string) string{
  return hmap_to_csv(levigoNS.ReadNSRecursive(key, db))
}


func PushKeyValNS(key string, val string) bool{
  return levigoNS.PushNS(key, val, db)
}


func DelKeyNS(key string) bool{
  levigoNS.DeleteNSRecursive(key, db)
  return true
}


func GetValTSDS(key string) string{
  return hmap_to_csv(levigoTSDS.ReadTSDS(key, db))
}


func PushKeyValTSDS(key string, val string,
                    year int, month int, day int,
                    hour int, min int, sec int) bool{
  key_time := time.Date(year, time.Month(month),
                        day, hour, min, sec, 0, time.UTC)
  levigoTSDS.PushTSDS(key, val, key_time, db)
  return true
}


func PushKeyValNowTSDS(key string, val string) bool{
  levigoTSDS.PushNowTSDS(key, val, db)
  return true
}


func DelKeyTSDS(key string) bool{
  levigoTSDS.DeleteTSDS(key, db)
  return true
}


func PushKeyMsgArrayTSDS(key string, msg_arr []string) bool{
  year, _ := strconv.Atoi(msg_arr[3])
  month, _ := strconv.Atoi(msg_arr[4])
  day, _ := strconv.Atoi(msg_arr[5])
  hour, _ := strconv.Atoi(msg_arr[6])
  min, _ := strconv.Atoi(msg_arr[7])
  sec, _ := strconv.Atoi(msg_arr[8])
  _value := strings.Join(msg_arr[9:], " ")
  return PushKeyValTSDS(key, _value, year, month, day, hour, min, sec)
}


func GetValTask(task_type string, key string) string{
  if task_type == "tsds" {
    return GetValTSDS(key)
  } else if task_type == "ns" {
    return GetValNS(key)
  }
  return GetVal(key)
}


func PushKeyValTask(task_type string, key string, value string) bool{
  if task_type == "tsds-now" {
    return PushKeyValNowTSDS(key, value)
  } else if task_type == "ns" {
    return PushKeyValNS(key, value)
  }
  return PushKeyVal(key, value)
}


func DelKeyTask(task_type string, key string) bool{
  if task_type == "tsds" {
    return DelKeyTSDS(key)
  } else if task_type == "ns" {
    return DelKeyNS(key)
  }
  return DelKey(key)
}
