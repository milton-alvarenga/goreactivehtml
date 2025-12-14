package encoder

import (
	"testing"

	"github.com/milton-alvarenga/goreactivehtml/internal/server/encode/byteprotocol"
	"github.com/milton-alvarenga/goreactivehtml/internal/server/encode/protocol"
)

// --- Bit Protocol Benchmarks ---

func BenchmarkBitProtocol_EncodeInsert(b *testing.B) {
	enc := protocol.Encoder{}
	data := []byte(`"hello world"`)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.EncodeInsert(uint32(i), data)
	}
}

func BenchmarkBitProtocol_EncodeUpdate(b *testing.B) {
	enc := protocol.Encoder{}
	data := []byte(`"hello world"`)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.EncodeUpdate(uint32(i), data)
	}
}

func BenchmarkBitProtocol_EncodeDelete(b *testing.B) {
	enc := protocol.Encoder{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.EncodeDelete(uint32(i))
	}
}

func BenchmarkBitProtocol_EncodeInsertRange(b *testing.B) {
	enc := protocol.Encoder{}
	payloads := [][]byte{
		[]byte(`"A"`),
		[]byte(`"B"`),
		[]byte(`"C"`),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.EncodeInsertRange(0, 2, payloads)
	}
}

func BenchmarkBitProtocol_EncodeUpdateRange(b *testing.B) {
	enc := protocol.Encoder{}
	payloads := [][]byte{
		[]byte(`"A"`),
		[]byte(`"B"`),
		[]byte(`"C"`),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.EncodeUpdateRange(0, 2, payloads)
	}
}

func BenchmarkBitProtocol_EncodeDeleteRange(b *testing.B) {
	enc := protocol.Encoder{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.EncodeDeleteRange(0, 100)
	}
}

// --- Byte Protocol Benchmarks ---

func BenchmarkByteProtocol_EncodeInsert(b *testing.B) {
	enc := byteprotocol.Encoder{}
	data := []byte(`"hello world"`)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.EncodeInsert(uint32(i), data)
	}
}

func BenchmarkByteProtocol_EncodeUpdate(b *testing.B) {
	enc := byteprotocol.Encoder{}
	data := []byte(`"hello world"`)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.EncodeUpdate(uint32(i), data)
	}
}

func BenchmarkByteProtocol_EncodeDelete(b *testing.B) {
	enc := byteprotocol.Encoder{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.EncodeDelete(uint32(i))
	}
}

func BenchmarkByteProtocol_EncodeInsertRange(b *testing.B) {
	enc := byteprotocol.Encoder{}
	payloads := [][]byte{
		[]byte(`"A"`),
		[]byte(`"B"`),
		[]byte(`"C"`),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.EncodeInsertRange(0, 2, payloads)
	}
}

func BenchmarkByteProtocol_EncodeUpdateRange(b *testing.B) {
	enc := byteprotocol.Encoder{}
	payloads := [][]byte{
		[]byte(`"A"`),
		[]byte(`"B"`),
		[]byte(`"C"`),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.EncodeUpdateRange(0, 2, payloads)
	}
}

func BenchmarkByteProtocol_EncodeDeleteRange(b *testing.B) {
	enc := byteprotocol.Encoder{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enc.EncodeDeleteRange(0, 100)
	}
}
