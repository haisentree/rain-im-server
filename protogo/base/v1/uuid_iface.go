package basev1

import (
	"database/sql"
	"database/sql/driver"
	"log/slog"

	"github.com/google/uuid"
	"github.com/uptrace/bun/schema"
)

var (
	_ driver.Valuer        = (*UUID)(nil)
	_ sql.Scanner          = (*UUID)(nil)
	_ schema.QueryAppender = (*UUID)(nil)
	_ slog.LogValuer       = (*UUID)(nil)
)

func (x *UUID) AppendQuery(_ schema.Formatter, b []byte) ([]byte, error) {
	b = append(b, '\'')
	b = append(b, x.UUID()...)

	return append(b, '\''), nil
}

func (x *UUID) Value() (driver.Value, error) {
	return x.ToUUID().Value()
}

func (x *UUID) Scan(src any) error {
	var u uuid.UUID
	if err := u.Scan(src); err != nil {
		return err
	}

	x.FromUUID(u)

	return nil
}

func (x *UUID) LogValue() slog.Value {
	if x == nil {
		return slog.StringValue("<uuid>")
	}

	return slog.StringValue(x.UUID())
}
