package polygon

// BaseResponse is the base structure for responses sent from a peer.
type IndicesResponse struct {
	Ev        string  `json:"ev"`
	Value     float64 `json:"val"`
	Ticker    string  `json:"T"`
	TimeStamp int64   `json:"t"`
}
