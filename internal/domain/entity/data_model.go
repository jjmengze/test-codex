package entity

type PutDataInput struct {
	RawData     []byte
	ProductCode string
	TraceID     string
	CustomerID  string
	Encoding    string
	SubType     string
	SourceID    string
}
