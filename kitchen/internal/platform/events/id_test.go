package events

import (
	"encoding/binary"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

var uuidV7 = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// t0 is a fixed clock reading: no assertion here depends on the real time.
var t0 = time.Date(2026, 11, 2, 10, 30, 15, 0, time.UTC)

func idString(b [16]byte) string { return pgtype.UUID{Bytes: b, Valid: true}.String() }

func TestIDs_SortInCreationOrder(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		clock func(i int) time.Time // the clock reading for the i-th id
	}{
		{
			name: "five ids in every millisecond of one second", n: 5000,
			clock: func(i int) time.Time { return t0.Add(time.Duration(i/5) * time.Millisecond) },
		},
		{
			name: "more ids in one millisecond than the 12-bit counter holds", n: 5000,
			clock: func(int) time.Time { return t0 },
		},
		{
			name: "readings a microsecond apart", n: 3000,
			clock: func(i int) time.Time { return t0.Add(time.Duration(i) * time.Microsecond) },
		},
		{
			name: "clock steps back by a second halfway", n: 100,
			clock: func(i int) time.Time {
				if i < 50 {
					return t0
				}
				return t0.Add(-time.Second)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s idSource
			prev := ""
			for i := range tt.n {
				id := idString(s.next(tt.clock(i)))
				if !uuidV7.MatchString(id) {
					t.Fatalf("id %d = %s, not a UUIDv7", i, id)
				}
				if id <= prev {
					t.Fatalf("id %d = %s, not after id %d = %s", i, id, i-1, prev)
				}
				prev = id
			}
		})
	}
}

func TestIDs_StartWithTheirMillisecond(t *testing.T) {
	var s idSource
	now := t0.Add(123*time.Millisecond + 456*time.Microsecond)

	id := s.next(now)

	var ms [8]byte
	copy(ms[2:], id[:6])
	if got, want := int64(binary.BigEndian.Uint64(ms[:])), now.UnixMilli(); got != want {
		t.Errorf("id %s carries millisecond %d, want %d", idString(id), got, want)
	}
}

func TestIDs_DifferAcrossProcesses(t *testing.T) {
	var a, b idSource

	idA, idB := a.next(t0), b.next(t0)

	if idA == idB {
		t.Errorf("two sources made the same id %s in the same millisecond", idString(idA))
	}
}
