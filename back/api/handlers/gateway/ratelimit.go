package gateway

import "time"

func rateLimiter() func() bool {
	const rate = 8   // per second average
	const min = 0.01 // kick threshold

	// Minimum time difference between messages
	// Network sometimes delivers two messages in quick succession
	const minDif = int64(time.Millisecond * 50)

	last := time.Now().UnixNano()
	var v float32 = 1.0
	return func() bool {
		now := time.Now().UnixNano()
		dif := now - last
		if dif < minDif {
			dif = minDif
		}
		v *= float32(rate*dif) / float32(time.Second)
		if v > 1.0 {
			v = 1.0
		}
		last = now
		return v > min
	}
}
