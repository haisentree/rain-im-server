package basev1

import (
	"rain-im-server/protogo/base/v1/internal/b2u"

	"github.com/google/uuid"
)

func (x *UUID) FromUUID(u uuid.UUID) {
	x.Hi, x.Lo = b2u.B264X2(u[:])
}

func (x *UUID) FromString(s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	x.FromUUID(id)

	return nil
}

func (x *UUID) ToUUID() uuid.UUID {
	v := uuid.UUID{}
	copy(v[:], x.ToBytes())

	return v
}

func (x *UUID) UUID() string { return x.ToUUID().String() }

func (x *UUID) ToBytes() []byte {
	if x == nil {
		return nil
	}

	b := b2u.U264X2(x.Hi, x.Lo)

	return b[:]
}

func NewUUID() *UUID {
	x := new(UUID)
	x.FromUUID(uuid.Must(uuid.NewV7()))

	return x
}

func (x *UUID) Equal(other any) bool {
	if x == nil || other == nil {
		return false // nil != nil
	}

	u, ok := other.(*UUID)
	if !ok {
		return false
	}

	return x.Hi == u.Hi && x.Lo == u.Lo
}

func (x *UUID) Formal() any {
	return x.UUID()
}

func UUID2String(arr []*UUID) []string {
	v := make([]string, len(arr))
	for i := range arr {
		v[i] = arr[i].UUID()
	}

	return v
}

func UUID2Any(arr []*UUID) []any {
	v := make([]any, len(arr))
	for i := range arr {
		v[i] = arr[i].UUID()
	}

	return v
}
