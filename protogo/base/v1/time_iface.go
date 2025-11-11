package basev1

import (
	"database/sql"
	"database/sql/driver"
	"log/slog"
	"time"
)

var (
	_ driver.Valuer = (*Time)(nil)
	_ sql.Scanner   = (*Time)(nil)
)

func Now() *Time             { return From(time.Now()) }
func From(t time.Time) *Time { return (&Time{}).Set(t) }

func (x *Time) Set(t time.Time) *Time { x.Seconds, x.Nanos = t.Unix(), int32(t.Nanosecond()); return x }
func (x *Time) AsTime() time.Time     { return time.Unix(x.GetSeconds(), int64(x.GetNanos())) }

func (x *Time) Parse(layout, value string) error {
	t, err := time.Parse(layout, value)
	if err != nil {
		return err
	}

	x.Set(t)

	return nil
}

func (x *Time) Value() (driver.Value, error) {
	return x.AsTime(), nil
}

func (x *Time) Scan(src any) error {
	switch src := src.(type) {
	case nil:
	case string:
		return x.Parse(time.RFC3339Nano, src)

	case []byte:
		t := &time.Time{}
		if err := t.UnmarshalText(src); err != nil {
			return err
		}

		x.Set(*t)

	case time.Time:
		x.Set(src)
	}

	return nil
}

func (x *Time) LogValue() slog.Value {
	if x == nil {
		return slog.StringValue("<time>")
	}

	return slog.StringValue(x.AsTime().Format(time.RFC3339))
}

func (x *Time) Equal(other any) bool {
	ot, ok := other.(*Time)
	if !ok || x == nil || ot == nil {
		return false
	}

	return x.Seconds == ot.Seconds && x.Nanos == ot.Nanos
}

func (x *Time) Formal() any {
	if x == nil {
		return nil
	}

	return x.AsTime()
}
