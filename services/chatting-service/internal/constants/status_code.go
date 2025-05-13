package constants

const (
	// Connection Lifecycle
	ServerConnected    = "902" // Successfully connected to WebSocket
	ServerDisconnected = "950" // Disconnected from server (client or server-side)
	ServerError        = "951" // Internal error from server

	// Matching States
	MatchRequestAccepted = "910" // Received tags, starting matching
	MatchSearching       = "911" // Searching for a match
	MatchNotFound        = "912" // No match found yet, waiting
	MatchFound           = "913" // Partner found

	// Messaging
	PartnerMessage       = "904" // Message from matched partner
	PartnerSystemMessage = "905" // System message about partner (e.g., joined, left)

	// Disconnection Events
	PartnerDisconnected = "920" // Partner disconnected
	UserDisconnected    = "921" // You disconnected (acknowledged)
	ReconnectAttempt    = "922" // Trying to re-match after disconnect

	// Error Handling
	InvalidRequestFormat = "960" // JSON parsing or format invalid
	UnauthorizedAccess   = "961" // Access denied or invalid user/session
	ModuleNotSupported   = "962" // Module type not recognized (e.g., "video" when not implemented)
)
