package encoder

import (
	"fmt"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"

	"github.com/milton-alvarenga/goreactivehtml/internal/server/encode/protocol"
)

func TestInsertUpdateDeleteProperties(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	properties := gopter.NewProperties(parameters)

	enc := protocol.Encoder{}

	properties.Property("Insert followed by Update produces updated value", prop.ForAll(
		func(pos uint32, value int, updated int) bool {
			bin1, _ := enc.EncodeInsert(pos, []byte(fmt.Sprintf("%d", value)))
			arr1, _ := decodeWithNode(bin1)

			bin2, _ := enc.EncodeUpdate(pos, []byte(fmt.Sprintf("%d", updated)))
			arr2, _ := decodeWithNode(bin2)

			return int(arr2[pos].(float64)) == updated
		},
		gen.UInt32(),
		gen.Int(),
		gen.Int(),
	))

	properties.Property("Delete removes element or has no effect", prop.ForAll(
		func(pos uint32, value int) bool {
			bin1, _ := enc.EncodeInsert(pos, []byte("1"))
			arr1, _ := decodeWithNode(bin1)

			bin2, _ := enc.EncodeDelete(pos)
			arr2, _ := decodeWithNode(bin2)

			return len(arr2) <= len(arr1)
		},
		gen.UInt32(),
		gen.Int(),
	))

	properties.TestingRun(t)
}
