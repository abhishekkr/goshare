package goshare

import (
	levigoNS "github.com/abhishekkr/levigoNS"
	abkleveldb "github.com/abhishekkr/levigoNS/leveldb"
	levigoTSDS "github.com/abhishekkr/levigoTSDS"
)

/* Empty Val for a given Key */
func DelKey(key string) bool {
	return abkleveldb.DelKey(key, db)
}

/* Delete a Namespace Key and all its value */
func DelKeyNS(key string) bool {
	return levigoNS.DeleteNSRecursive(key, db)
}

/* Delete all keys under given namespace, same as NS */
func DelKeyTSDS(key string) bool {
	current_val := levigoTSDS.DeleteTSDS(key, db)
	if len(current_val) > 0 {
		return false
	}
	return true
}

/* Delete a key on task-type */
func DelKeyTask(task_type string, key string) bool {
	if task_type == "tsds" {
		return DelKeyTSDS(key)

	} else if task_type == "ns" {
		return DelKeyNS(key)

	}

	return DelKey(key)
}
