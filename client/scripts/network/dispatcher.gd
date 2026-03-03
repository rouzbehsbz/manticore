class_name Dispatcher
extends RefCounted

const Pb = preload("res://scripts/network/packets.gd")

var handlers := {}

func register_handler(packet_id: int, handler: Callable) -> void:
	handlers[packet_id] = handler
	
func dispatch(packet: Pb.Packet) -> void:
	var packet_id = packet.get_id()
	
	if !handlers.has(packet_id):
		push_error("no handler registered for this packet ID: " + str(packet_id))
		return
	
	var handler: Callable = handlers.get(packet_id)
	handler.call(packet)
