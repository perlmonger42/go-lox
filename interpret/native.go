package interpret

import (
	"time"

	"github.com/perlmonger42/go-lox/token"
)

// ClockNative implements the builtin function `clock()`
type ClockNative struct {
}

var _ token.Object = ClockNative{}

func (f ClockNative) Show() string   { return f.String() }
func (f ClockNative) String() string { return `[native function "clock()"]` }

var _ Callable = &ClockNative{}

var EpochStart = time.Date(1970, time.January, 1, 00, 00, 00, 00, time.UTC)

func (f ClockNative) Arity() int { return 0 }

// Return the number of seconds since 00:00:00 on 1 January 1970 (UTC).
func (f ClockNative) Call(i T, arguments []token.Value) token.Value {
	var elapsedTimeSinceTheEpoch time.Duration = time.Since(EpochStart)
	// A time.Duration represents the elapsed time between two instants as an
	// int64 nanosecond count. There are 1e9 nanoseconds in a second.
	// Duration.Seconds() returns the duration as a floating point number of
	// seconds.
	return token.NumberValue{elapsedTimeSinceTheEpoch.Seconds()}
}

func (f ClockNative) EqualsObject(o token.Object) bool {
	return false
}

// StrNative implements the builtin function `str(whatever)`
type StrNative struct {
}

var _ token.Object = &StrNative{}

func (f StrNative) Show() string   { return f.String() }
func (f StrNative) String() string { return `[native function "str(x)"]` }

var _ Callable = StrNative{}

func (f StrNative) Arity() int { return 1 }

// Return the number of seconds since 00:00:00 on 1 January 1970 (UTC).
func (f StrNative) Call(i T, arguments []token.Value) token.Value {
	return token.StringValue{arguments[0].String()}
}

func (f StrNative) EqualsObject(o token.Object) bool {
	return false
}
