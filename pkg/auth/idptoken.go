package auth

// IDPTokenPayload Spec: https://wiki.jarvis.trendmicro.com/display/SG/IdP+Auth+Service+External+API
type IDPTokenPayload struct {
	ConsumerProductID string `json:"cpid"`
	ProducerProductID string `json:"ppid"`
	CustomerID        string `json:"cid"` // same as company_id
	UserID            string `json:"uid"` // same as X-Sender-ID
	Payload           string `json:"pl"`
	IssueTime         int64  `json:"it"`
	ExpiredTime       int64  `json:"et"`
}
