package networking

const (
	PENDING MessageStatus = iota
	PROCESSING
	DELIVERED
	FAILED
	COMPLETED
)
