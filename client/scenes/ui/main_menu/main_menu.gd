extends Control

const Pb = preload("res://scripts/network/packets.gd")

@onready var username_input: LineEdit = $%UsernameInput
@onready var password_input: LineEdit = $%PasswordInput
@onready var login_button: Button = $%LoginButton
@onready var register_button: Button = $%RegisterButton

var sm := SessionManager

func _ready() -> void:
	sm.session.dispatcher.register_handler(PacketRegistry.PacketId.REGISTER_RESPONSE, _on_register_handler)
	sm.session.dispatcher.register_handler(PacketRegistry.PacketId.LOGIN_RESPONSE, _on_login_handler)

func _on_login_button_pressed() -> void:
	var packet = PacketRegistry.build_login_packet(username_input.text, password_input.text)
	sm.session.write(packet)

func _on_register_button_pressed() -> void:
	var packet = PacketRegistry.build_register_packet(username_input.text, password_input.text)
	sm.session.write(packet)

func _on_login_handler(packet: Pb.Packet) -> void:
	var res = packet.get_login_response()
	print(res)
	
func _on_register_handler(packet: Pb.Packet) -> void:
	var res = packet.get_register_response()
	print(res)
