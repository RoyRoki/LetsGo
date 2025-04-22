package wsconstants

type WSMessageCode int32

const (
	WSMessageServerInfo WSMessageCode = 300
	WSMsgPartnerJoined  WSMessageCode = 301
	WSMsgPartnerLeft    WSMessageCode = 302
	WSMsgTyping         WSMessageCode = 303
	WSMsgPartnerMessage WSMessageCode = 400
)

var WSMessages = map[WSMessageCode]string{
	WSMessageServerInfo: "Server Message",
	WSMsgPartnerJoined:  "Your partner has joined",
	WSMsgPartnerLeft:    "Your partner has left",
	WSMsgTyping:         "Partner is typing...",
	WSMsgPartnerMessage: "Message from partner",
}
