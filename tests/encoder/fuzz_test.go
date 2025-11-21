package encoder

import (
	"math/rand"
	"testing"

	"github.com/milton-alvarenga/goreactivehtml/internal/server/encode/protocol"
)

func FuzzBinaryProtocol(f *testing.F) {
	enc := protocol.Encoder{}

	f.Add(uint32(0), []byte(`"hello"`))
	f.Add(uint32(10), []byte(`123`))

	f.Fuzz(func(t *testing.T, pos uint32, payload []byte) {
		op := rand.Intn(3)
		var bin []byte
		var err error

		switch op {
		case 0:
			bin, err = enc.EncodeInsert(pos, payload)
		case 1:
			bin, err = enc.EncodeUpdate(pos, payload)
		case 2:
			bin, err = enc.EncodeDelete(pos)
		}

		if err != nil {
			return // encoder correctly rejected invalid input
		}

		_, err = decodeWithNode(bin)
		if err != nil {
			t.Fatalf("decoder crashed on fuzz input: %v", err)
		}
	})
}
