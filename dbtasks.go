package goshare

/*
[PATTERN]
action {read, push, delete}
type {default, ns, tsds, now}

## message_array here is devoided of axn and key_type
non-tsds {key&val, :type-data}
tsds(-*) {tdot&key&val, tdot&:type-data}
*/

/*
DBTasks can be provided standard Packet Data Array from any communication Protocol.
Communications handled on byte streams can use it by passing standard-ized packet-array
It prepares Packet and passes on to TasksOnPacket.
0MQ directly utilizes it.
*/
func DBTasks(packetArray []string) ([]byte, bool) {
	packet := CreatePacket(packetArray)
	return DBTasksOnPacket(packet)
}

/*
DBTasksOnPacket can utilize fromulated Packet.
Communication can directly create packet and pass it here.
HTTP directly utilizes it directly. 0MQ indirectly.
*/
func DBTasksOnPacket(packet Packet) ([]byte, bool) {
	response := ""
	axnStatus := false

	switch packet.DBAction {
	case "read":
		// returns axn error if key has empty value, if you gotta store then store, don't keep placeholders
		response = ReadFromPacket(packet)
		if response != "" {
			axnStatus = true
		}

	case "push":
		axnStatus = PushFromPacket(packet)

	case "delete":
		axnStatus = DeleteFromPacket(packet)
	}

	return []byte(response), axnStatus
}
