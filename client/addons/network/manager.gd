class_name NetworkManager
extends Node

const Session = preload("res://addons/network/session.gd")

@export var url: String = "ws://localhost:3000/ws"

var session: Session

func _ready() -> void:
	session = Session.new(url)
	session.connected.connect(_on_connect)
	session.disconnected.connect(_on_disconnect)

func _process(delta: float) -> void:
	if !session:
		return
	
	session.poll()
	session.flush()

func _on_connect() -> void:
	print("Connected")

func _on_disconnect() -> void:
	print("Disconnected")
