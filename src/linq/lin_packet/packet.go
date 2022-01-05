package lin_packet

type EN_PACKET_TYPE int

const (
	EN_PACKET_TYPE_none               = 0
	EN_PACKET_TYPE_client_sign_in     = 1
	EN_PACKET_TYPE_clinet_sign_in_res = 2
)

type PacketBase interface {
	PacketType() EN_PACKET_TYPE
}

type PacketClientSignin struct {
	clientName_ string
}
func (*PacketClientSignin) PacketType() EN_PACKET_TYPE {
	return EN_PACKET_TYPE_client_sign_in
}

type PacketClientSigninRes struct {
	clientId_ int64
}
func (*PacketClientSigninRes) PacketType() EN_PACKET_TYPE {
	return EN_PACKET_TYPE_clinet_sign_in_res
}
