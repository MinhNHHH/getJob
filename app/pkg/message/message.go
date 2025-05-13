package message

type MType string

type Wrapper struct {
	Type MType
	Data interface{}
	From string
	To   string
}

const (
	// TODO refactor to make these message type as part of message
	// we probably only need RTC, Control, Connect types
	TRTCOffer     MType = "Offer"
	TRTCAnswer    MType = "Answer"
	TRTCCandidate MType = "Candidate"

	// Client can order the host to refresh the terminal
	// Used in case client resize and need to update the content to display correctly
	TTermRefresh MType = "Refresh"

	TWSPing MType = "Ping"

	// when connect, client will first send a connect messag
	TCConnect = "Connect"
	TCSend    = "send"
	TCYes     = "y"
	TCNo      = "No"
	// connection's response
	TCAuthenticated   = "Authenticated"
	TCUnauthenticated = "Unauthenticated"

	TCUnsupportedVersion = "UnsupportedVersion"
)
