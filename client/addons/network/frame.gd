class_name Frame
extends RefCounted

const Pb = preload("res://addons/network/packets.gd")

var packets: Array[Pb.Packet] = []

static func build_frame(bytes: PackedByteArray) -> Frame:
	var frame = Frame.new()
	
	if bytes.is_empty():
		push_warning("Empty frame received")
		return frame
	
	var reader = StreamPeerBuffer.new()
	reader.data_array = bytes
	reader.big_endian = true
	
	while reader.get_position() < reader.get_size():
		if reader.get_available_bytes() < 2:
			push_error("invalid length prefix packet")
			break
			
		var length = reader.get_u16()
		if length == 0:
			push_error("invalid zero length packet")
			break
			
		if reader.get_available_bytes() < length:
			push_error("packet length is not compatible with packet data")
			break
			
		var raw_packet = reader.get_data(length)
		var packet = Pb.Packet.new()
		
		var err = packet.from_bytes(raw_packet)
		if err != Pb.PB_ERR.NO_ERRORS:
			push_error("failed to parse raw packet data")
			break
			
		frame.packets.append(packet)
	
	return frame

func append(packet: Pb.Packet) -> void:
	packets.append(packet)

func bytes() -> PackedByteArray:
	var writer = StreamPeerBuffer.new()
	writer.big_endian = true
	
	for packet in packets:
		var raw_packet = packet.to_bytes()
		if raw_packet.is_empty():
			push_warning("invalid raw empty packet")
			continue
			
		var length = raw_packet.size()
		
		writer.put_u16(length)
		writer.put_data(raw_packet)
	
	return writer.data_array

func size() -> int:
	return packets.size()
	
func empty() -> void:
	packets.resize(0)
	
