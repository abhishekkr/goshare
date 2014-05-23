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
	return levigoTSDS.DeleteTSDS(key, db)
}

/* Delete a key on task-type */
func DelKeyTask(key_type string, key string) bool {
	if key_type == "tsds" {
		return DelKeyTSDS(key)

	} else if key_type == "ns" {
		return DelKeyNS(key)

	}

	return DelKey(key)
}
