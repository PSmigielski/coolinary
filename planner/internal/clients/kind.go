package clients

import "connectrpc.com/connect"

// Kind is what feature code switches on after a failed call (PLN-01-02 rule 4).
type Kind int

const (
	// OK: no error.
	OK Kind = iota
	// Unavailable: the service is down, unreachable or too slow; fall back (non-AI result, cache) or answer UNAVAILABLE.
	Unavailable
	// Rejected: the service answered no (bad input, rule, not found, quota, conflict). connect.CodeOf(err) keeps
	// its code.
	Rejected
	// Bug: anything else.
	Bug
)

func (k Kind) String() string {
	switch k {
	case OK:
		return "ok"
	case Unavailable:
		return "unavailable"
	case Rejected:
		return "rejected"
	default:
		return "bug"
	}
}

// KindOf sorts an error from one of the clients into a Kind. Connection errors reach it as UNAVAILABLE.
func KindOf(err error) Kind {
	if err == nil {
		return OK
	}
	switch connect.CodeOf(err) {
	case connect.CodeUnavailable, connect.CodeDeadlineExceeded:
		return Unavailable
	case connect.CodeInvalidArgument, connect.CodeFailedPrecondition, connect.CodeNotFound,
		connect.CodePermissionDenied, connect.CodeResourceExhausted,
		connect.CodeAborted, connect.CodeAlreadyExists: // Q237
		return Rejected
	default:
		return Bug
	}
}
