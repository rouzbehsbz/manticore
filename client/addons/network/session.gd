class_name Session
extends RefCounted

const Pb = preload("res://addons/network/packets.gd")
const Frame = preload("res://addons/network/frame.gd")

signal connected()
signal disconnected()
signal packet_received(packet: Pb.Packet)

var socket := WebSocketPeer.new()
var frame := Frame.new()

var _is_running := false
var _last_state := WebSocketPeer.STATE_CLOSED

func _init(url: String) -> void:
	_connect(url)
	
func poll() -> void:
	if !_is_running:
		return
		
	socket.poll()
	
	_check_state()
	_receive()

func _connect(url: String) -> void:
	if socket.get_ready_state() != WebSocketPeer.STATE_CLOSED:
		push_warning("socket is already connected")
		return
		
	var err := socket.connect_to_url(url)
	if err != OK:
		push_error("failed to connect to the server")
		return
		
	_is_running = true

func _check_state() -> void:
	var state = socket.get_ready_state()
	if state == _last_state:
		return
	
	match state:
		WebSocketPeer.STATE_OPEN:
			connected.emit()
		WebSocketPeer.STATE_CONNECTING:
			pass
		WebSocketPeer.STATE_CLOSED, WebSocketPeer.STATE_CLOSING:
			if _is_running:
				_is_running = false
				disconnected.emit()
	
	_last_state = state

func _receive() -> void:
	while socket.get_available_packet_count() > 0:
		var bytes = socket.get_packet()
		if socket.was_string_packet():
			push_warning("expected binary format packet")
			continue
		
		var frame = Frame.build_frame(bytes)
		for packet in frame.packets:
			packet_received.emit(packet)
			
func write(packet: Pb.Packet) -> void:
	frame.append(packet)

func flush() -> void:
	if frame.size() == 0:
		return
	
	var bytes = frame.bytes()
	
	var err = socket.put_packet(bytes)
	if err != OK:
		push_error("failed to put packet in socket")
		return
	
	frame.empty()
