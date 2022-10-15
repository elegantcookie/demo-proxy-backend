package attempts

import (
	"time"
)

// DoWithAttempts allows to execute function with number of attempts and delay between tries
func DoWithAttempts(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}
	return
}
