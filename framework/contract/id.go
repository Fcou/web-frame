package contract

const IDKey = "fcou:id"

type IDService interface {
	NewID() string
}
