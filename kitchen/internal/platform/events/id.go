package events

import (
	"crypto/rand"
	"encoding/binary"
	"sync"
	"time"
)

// ids makes every envelope id of this process.
var ids idSource

// idSource makes envelope ids: UUIDv7 (RFC 9562) that sort in creation order within the process, so the relay's
// id order is creation order. Layout: 48 bits of Unix milliseconds, the version, a 12-bit counter for ids in the
// same millisecond (RFC 9562 section 6.2, method 1), the variant, 62 random bits. When the counter runs out or the
// clock steps back, the next id takes the last one's millisecond and counter plus one.
type idSource struct {
	mu   sync.Mutex
	last uint64 // milliseconds<<12 | counter of the last id
}

func (s *idSource) next(now time.Time) [16]byte {
	ms := max(now.UnixMilli(), 0)
	s.mu.Lock()
	v := max(uint64(ms)<<12, s.last+1)
	s.last = v
	s.mu.Unlock()

	var id [16]byte
	// 48 bits of milliseconds, version 7, the 12-bit counter; then the variant (10) and random bits.
	binary.BigEndian.PutUint64(id[:8], v>>12<<16|0x7000|v&0xfff)
	_, _ = rand.Read(id[8:]) // never fails (crypto/rand.Read)
	id[8] = 0x80 | id[8]&0x3f
	return id
}
