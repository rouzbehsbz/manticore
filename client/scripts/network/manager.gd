class_name Manager
extends Node

const Session = preload("res://scripts/network/session.gd")

@export var url: String = "ws://localhost:3000/ws"
@export var reconnect_delay := 5.0

var session: Session
var reconnect_timer: Timer

func _ready() -> void:
	session = Session.new()
	
	session.connected.connect(_on_connect)
	session.disconnected.connect(_on_disconnect)
	
	session.connect_to(url)
	
	reconnect_timer = Timer.new()
	reconnect_timer.connect("timeout", _attempt_reconnect)
	add_child(reconnect_timer)

func _process(_delta: float) -> void:
	if !session:
		return
	
	session.poll()
	session.flush()

func _on_connect() -> void:
	print("Connected")

func _on_disconnect() -> void:
	print("Disconnected")
	print("Reconnecting...")
	reconnect_timer.start(reconnect_delay)

func _attempt_reconnect() -> void:
	session.connect_to(url)
	reconnect_timer.stop()
