package repository

const (
	StatusPending uint64 = 1 // Pending received values
	StatusClose   uint64 = 2 // Bid is closed
	StatusError   uint64 = 3 // Bid has error
)
