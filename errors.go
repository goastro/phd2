package phd2

// Error is an error that can be returned by this package.
type Error string

func (err Error) Error() string {
	return string(err)
}

const (
	// ErrNotImplemented is returned if the function called has not been
	// implemented yet.
	ErrNotImplemented = Error("not implemented")
	// ErrNotConnected is returned if the client is not connected to the PHD2
	// server.
	ErrNotConnected = Error("not connected")
)
