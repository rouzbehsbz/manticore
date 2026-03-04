class_name Session
extends RefCounted

const Pb = preload("res://scripts/network/packets.gd")
const Frame = preload("res://scripts/network/frame.gd")

signal connected()
signal disconnected()

var dispatcher := Dispatcher.new()

var _socket := WebSocketPeer.new()
var _send_frame := Frame.new()
var _is_connect_attempt := false 
var _last_state := WebSocketPeer.STATE_CLOSED
	
func poll() -> void:
	_socket.poll()
	
	_check_state()
	_receive()

func connect_to(url: String) -> void:
	_is_connect_attempt = true
	
	if _socket.get_ready_state() != WebSocketPeer.STATE_CLOSED:
		push_warning("socket is already connected")
		return
		
	var err := _socket.connect_to_url(url)
	if err != OK:
		push_error("failed to connect to the server")
		return

func _check_state() -> void:
	var state = _socket.get_ready_state()
	if state == _last_state && !_is_connect_attempt:
		return
	
	match state:
		WebSocketPeer.STATE_OPEN:
			connected.emit()
		WebSocketPeer.STATE_CONNECTING:
			pass
		WebSocketPeer.STATE_CLOSED, WebSocketPeer.STATE_CLOSING:
			disconnected.emit()
	
	_last_state = state
	_is_connect_attempt = false

func _receive() -> void:
	while _socket.get_available_packet_count() > 0:
		var bytes = _socket.get_packet()
		if _socket.was_string_packet():
			push_warning("expected binary format packet")
			continue
		
		var frame = Frame.build_frame(bytes)
		for packet: Pb.Packet in frame.packets:
			dispatcher.dispatch(packet)
			
func write(packet: Pb.Packet) -> void:
	_send_frame.append(packet)

func flush() -> void:
	if _send_frame.size() == 0 or _last_state != WebSocketPeer.STATE_OPEN:
		return
	
	var bytes = _send_frame.to_bytes()
	
	var err = _socket.put_packet(bytes)
	if err != OK:
		push_error("failed to put packet in socket")
		return
	
	_send_frame.empty()
