package goshare

import (
	golhashmap "github.com/abhishekkr/gol/golhashmap"
	levigoNS "github.com/abhishekkr/levigoNS"
	abkleveldb "github.com/abhishekkr/levigoNS/leveldb"
	levigoTSDS "github.com/abhishekkr/levigoTSDS"
)

/* Get value of given key */
func GetVal(key string) string {
	return abkleveldb.GetVal(key, db)
}

/* Get value for all descendents of Namespace */
func GetValNS(key string) string {
	hashmap := levigoNS.ReadNSRecursive(key, db)
	return golhashmap.Hashmap_to_csv(hashmap)
}

/* Get value for the asked time-frame key, aah same NS */
func GetValTSDS(key string) string {
	return golhashmap.Hashmap_to_csv(levigoTSDS.ReadTSDS(key, db))
}

/* Get a value based on task-type */
func GetValTask(task_type string, key string) string {
	if task_type == "tsds" {
		return GetValTSDS(key)

	} else if task_type == "ns" {
		return GetValNS(key)

	}

	return GetVal(key)
}
