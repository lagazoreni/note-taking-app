package platform

import "time"

type Clock interface{ Now() time.Time }

type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now().UTC() }

type FixedClock struct{ value time.Time }

func NewFixedClock(value time.Time) *FixedClock { return &FixedClock{value: value.UTC()} }
func (c *FixedClock) Now() time.Time            { return c.value }
func (c *FixedClock) Set(value time.Time)       { c.value = value.UTC() }

func FormatTimestamp(value time.Time) string         { return value.UTC().Format(time.RFC3339Nano) }
func ParseTimestamp(value string) (time.Time, error) { return time.Parse(time.RFC3339Nano, value) }
