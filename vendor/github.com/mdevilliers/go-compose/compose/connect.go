package compose

import "time"

const (
	defaultRetryCount     = 10
	defaultBaseRetryDelay = 100 * time.Millisecond
)

// Connect attempts to connect to a container using the given connector function.
// Use retryCount and retryDelay to configure the number of retries and the time waited between them (using exponential backoff).
func Connect(retryCount int, baseRetryDelay time.Duration, connectFunc func() error) error {
	var err error

	for i := 0; i < retryCount; i++ {
		err = connectFunc()
		if err == nil {
			return nil
		}
		time.Sleep(baseRetryDelay)
		baseRetryDelay *= 2
	}

	return err
}

// MustConnect is like Connect, but panics on error.
func MustConnect(retryCount int, baseRetryDelay time.Duration, connectFunc func() error) {
	if err := Connect(retryCount, baseRetryDelay, connectFunc); err != nil {
		panic(err)
	}
}

// ConnectWithDefaults is like Connect, with default values for retryCount and retryDelay.
func ConnectWithDefaults(connectFunc func() error) error {
	return Connect(defaultRetryCount, defaultBaseRetryDelay, connectFunc)
}

// MustConnectWithDefaults is like ConnectWithDefaults, but panics on error.
func MustConnectWithDefaults(connectFunc func() error) {
	if err := ConnectWithDefaults(connectFunc); err != nil {
		panic(err)
	}
}
