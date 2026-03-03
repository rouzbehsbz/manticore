class_name PacketRegistry

const Pb = preload("res://scripts/network/packets.gd")

enum PacketId {
	LOGIN_REQUEST = 0,
	LOGIN_RESPONSE = 1,
	REGISTER_REQUEST = 2,
	REGISTER_RESPONSE = 3,
}

static func build_register_packet(username: String, password: String) -> Pb.Packet:
	var registerRequest = Pb.RegisterRequest.new()
	registerRequest.set_username(username)
	registerRequest.set_password(password)
	
	var packet = Pb.Packet.new()
	packet.set_id(PacketId.REGISTER_REQUEST)
	packet.__register_request.value = registerRequest
	
	return packet

static func build_login_packet(username: String, password: String) -> Pb.Packet:
	var loginRequest = Pb.LoginRequest.new()
	loginRequest.set_username(username)
	loginRequest.set_password(password)
	
	var packet = Pb.Packet.new()
	packet.set_id(PacketId.LOGIN_REQUEST)
	packet.__login_request.value = loginRequest

	return packet
