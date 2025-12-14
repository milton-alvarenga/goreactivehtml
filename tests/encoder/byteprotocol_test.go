package encoder

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"testing"

	"github.com/milton-alvarenga/goreactivehtml/internal/server/encode/byteprotocol"
)

//
// --- Test Helpers ---
//

// decodeWithNode sends the binary payload to Node.js along with an optional initial array
func decodeWithNode_byte(payload []byte, initial []interface{}) ([]interface{}, error) {
	cmd := exec.Command("node", "../decoder/node_byte.js")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// Encode initial array JSON
	initialJSON, err := json.Marshal(initial)
	if err != nil {
		return nil, err
	}

	// Write 4-byte length prefix
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(initialJSON)))

	if _, err := stdin.Write(lenBuf); err != nil {
		return nil, err
	}
	if _, err := stdin.Write(initialJSON); err != nil {
		return nil, err
	}

	// Write raw payload
	if _, err := stdin.Write(payload); err != nil {
		return nil, err
	}
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		log.Println("Node.js stderr:", stderr.String())
		return nil, err
	}

	log.Println("Node.js stdout:", stdout.String())

	var result []interface{}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func must_byte[T any](t *testing.T, v T, err error) T {
	if err != nil {
		t.Fatal(err)
	}
	return v
}

//
// --- Tests ---
//

func TestInsertSingle_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`"A"`))
	t.Log(v)
	bin := must_byte(t, v, err)
	t.Log(bin)
	t.Logf("Encoded binary data: %v", bin)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	t.Log(val)
	out := must_byte(t, val, err)
	t.Log(out)

	if out[0] != "A" {
		t.Fatalf("expected A, got %v", out)
	}
}

func TestEncodeInsert_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`"hello"`))
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[0] != "hello" {
		t.Fatalf("expected hello, got %v", out)
	}
}

func TestUpdateSingle_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}
	v, err := enc.EncodeUpdate(5, []byte(`123`))
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[5] != float64(123) {
		t.Fatalf("expected 123, got %v", out)
	}
}

func TestEncodeUpdate_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}
	v, err := enc.EncodeUpdate(0, []byte(`123`))
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[0] != float64(123) {
		t.Fatalf("expected 123, got %v", out)
	}
}

func TestPartialUpdate_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}
	v, err := enc.EncodePartialUpdate(5, []byte(`"patched"`))
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[5] != "patched" {
		t.Fatalf("expected patched, got %v", out)
	}
}

func TestBulkPartialUpdate_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	patches := []byteprotocol.PartialPatch{
		{Pos: 10, Data: []byte(`"A"`)},
		{Pos: 12, Data: []byte(`"B"`)},
	}

	v, err := enc.EncodePartialUpdateRange(10, 20, patches)
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[10] != "A" || out[12] != "B" {
		t.Fatalf("expected sparse patches, got %v", out)
	}
}

func TestDeleteSingle_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	// Start with empty array
	initial := []interface{}{}

	// Insert first element
	v, err := enc.EncodeInsert(0, []byte(`1`))
	out, err := decodeWithNode_byte(must_byte(t, v, err), initial)
	out = must_byte(t, out, err)

	if len(out) != 1 || int(out[0].(float64)) != 1 {
		t.Fatalf("expected [1], got %v", out)
	}

	// Insert second element
	v, err = enc.EncodeInsert(1, []byte(`2`))
	out, err = decodeWithNode_byte(must_byte(t, v, err), out)
	out = must_byte(t, out, err)

	if len(out) != 2 || int(out[1].(float64)) != 2 {
		t.Fatalf("expected [1,2], got %v", out)
	}

	// Insert third element
	v, err = enc.EncodeInsert(2, []byte(`3`))
	out, err = decodeWithNode_byte(must_byte(t, v, err), out)
	out = must_byte(t, out, err)

	if len(out) != 3 || int(out[2].(float64)) != 3 {
		t.Fatalf("expected [1,2,3], got %v", out)
	}

	// Now delete the second element (index 1)
	v, err = enc.EncodeDelete(1)
	out, err = decodeWithNode_byte(must_byte(t, v, err), out)
	out = must_byte(t, out, err)

	if len(out) != 2 || int(out[0].(float64)) != 1 || int(out[1].(float64)) != 3 {
		t.Fatalf("delete failed, expected [1,3], got %v", out)
	}
}

func TestInsertRange_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	payloads := [][]byte{
		[]byte(`"A"`),
		[]byte(`"B"`),
		[]byte(`"C"`),
	}

	v, err := enc.EncodeInsertRange(10, 12, payloads)
	val, err := decodeWithNode_byte(must_byte(t, v, err), []interface{}{})
	out := must_byte(t, val, err)

	if out[10] != "A" || out[11] != "B" || out[12] != "C" {
		t.Fatalf("range insert failed: %v", out)
	}
}

func TestInsertRangeText_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	payloads := [][]byte{
		[]byte(`"ABC"`),
		[]byte(`"BCD"`),
		[]byte(`"CDE"`),
	}

	v, err := enc.EncodeInsertRange(4, 6, payloads)
	val, err := decodeWithNode_byte(must_byte(t, v, err), []interface{}{})
	out := must_byte(t, val, err)

	if out[4] != "ABC" || out[5] != "BCD" || out[6] != "CDE" {
		t.Fatalf("range insert failed: %v", out)
	}
}

func TestUpdateRange_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	payloads := [][]byte{
		[]byte(`10`),
		[]byte(`20`),
		[]byte(`30`),
	}

	v, err := enc.EncodeUpdateRange(3, 5, payloads)
	val, err := decodeWithNode_byte(must_byte(t, v, err), []interface{}{})
	out := must_byte(t, val, err)

	if out[3] != float64(10) ||
		out[4] != float64(20) ||
		out[5] != float64(30) {
		t.Fatalf("range update failed: %v", out)
	}
}

func TestDeleteRange_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	// Start with 10 elements
	fullArray := make([][]byte, 10)
	for i := 0; i < 10; i++ {
		fullArray[i] = []byte(`1`)
	}

	// Bulk insert/update the full array in Node
	v, err := enc.EncodeUpdateRange(0, uint32(len(fullArray)-1), fullArray)
	out, err := decodeWithNode_byte(must_byte(t, v, err), []interface{}{})
	must_byte(t, out, err)
	if len(out) != 10 {
		t.Fatalf("expected 10 elements remain, got %v", out)
	}

	// Delete range 3..6 locally
	fullArray = append(fullArray[:3], fullArray[7:]...)

	// Send the new full array to Node again
	v, err = enc.EncodeUpdateRange(0, uint32(len(fullArray)-1), fullArray)
	out, err = decodeWithNode_byte(must_byte(t, v, err), []interface{}{})
	outArray := must_byte(t, out, err)

	if len(outArray) != 6 {
		t.Fatalf("expected 6 elements remain, got %v", outArray)
	}
}

func TestInsertAtEnd_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	state := []interface{}{} // initial empty array

	// Insert 0
	v, err := enc.EncodeInsert(0, []byte(`0`))
	state, err = decodeWithNode_byte(must_byte(t, v, err), state)
	state = must_byte(t, state, err)

	// Insert 1
	v, err = enc.EncodeInsert(1, []byte(`1`))
	state, err = decodeWithNode_byte(must_byte(t, v, err), state)
	state = must_byte(t, state, err)

	// Insert 2
	v, err = enc.EncodeInsert(2, []byte(`2`))
	state, err = decodeWithNode_byte(must_byte(t, v, err), state)
	state = must_byte(t, state, err)

	// Insert END at the end
	v, err = enc.EncodeInsert(3, []byte(`"END"`))
	state, err = decodeWithNode_byte(must_byte(t, v, err), state)
	state = must_byte(t, state, err)

	if state[3] != "END" {
		t.Fatalf("insert at end failed: %v", state)
	}
}

func TestInsertBeyondEnd_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}
	v, err := enc.EncodeInsert(10, []byte(`"X"`))
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[10] != "X" {
		t.Fatalf("expected X at pos 10, got %v", out)
	}

	for i := 0; i < 10; i++ {
		if out[i] != nil {
			t.Fatalf("expected null at %d got %v", i, out[i])
		}
	}
}

func TestUpdateBeyondEnd_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}
	v, err := enc.EncodeUpdate(15, []byte(`999`))
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[15] != float64(999) {
		t.Fatalf("expected 999, got %v", out)
	}
}

func TestZeroLengthPayload_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`""`))
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[0] != "" {
		t.Fatalf("expected empty string, got %v", out)
	}
}

func TestJSONString_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`"hello world"`))
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[0] != "hello world" {
		t.Fatalf("expected hello world, got %v", out)
	}
}

func TestJSONNumber_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`42`))
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[0] != float64(42) {
		t.Fatalf("expected 42, got %v", out)
	}
}

func TestJSONObject_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`{"a":1,"b":2}`))
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	obj := out[0].(map[string]interface{})
	if obj["a"] != float64(1) || obj["b"] != float64(2) {
		t.Fatalf("object mismatch: %v", obj)
	}
}

func TestJSONArray_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`[1,2,3]`))
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	arr := out[0].([]interface{})
	if len(arr) != 3 || arr[1] != float64(2) {
		t.Fatalf("array mismatch: %v", arr)
	}
}

func TestMax24BitPosition_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}
	pos := uint32(0xFFFFFF)

	v, err := enc.EncodeInsert(pos, []byte(`"MAXPOS"`))
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[pos] != "MAXPOS" {
		t.Fatalf("expected MAXPOS at %d", pos)
	}
}

func TestBulkInsertMixedSize_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	payloads := [][]byte{
		[]byte(`"A"`),
		[]byte(`"BBBBBB"`),
		[]byte(`"C"`),
		[]byte(`"DDDDDDDDDDDDDDDDDDD"`),
	}

	v, err := enc.EncodeInsertRange(5, 8, payloads)
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[5] != "A" ||
		out[6] != "BBBBBB" ||
		out[7] != "C" ||
		out[8] != "DDDDDDDDDDDDDDDDDDD" {
		t.Fatalf("bulk insert mismatch: %v", out)
	}
}

func TestRandomFixedCases_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	// Start with empty array state
	state := []interface{}{}

	for i := 0; i < 200; i++ {
		pos := uint32(i % 50)          // random but deterministic
		valStr := fmt.Sprintf("%d", i) // JSON number (not string)

		// Encode insert
		bin, err := enc.EncodeInsert(pos, []byte(valStr))
		if err != nil {
			t.Fatalf("encode error: %v", err)
		}

		// Feed previous state
		val, err := decodeWithNode_byte(bin, state)
		state = must_byte(t, val, err)

		// Ensure padding behavior
		if state[pos].(float64) != float64(i) {
			t.Fatalf("random test failed at iteration=%d pos=%d (value=%v full=%v)",
				i, pos, state[pos], state)
		}
	}
}

func TestSparsePartialUpdate_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	for i := 0; i < 5; i++ {
		v, err := enc.EncodeInsert(uint32(i), []byte(fmt.Sprintf(`%d`, i)))
		val, err := decodeWithNode_byte(must_byte(t, v, err), []interface{}{})
		must_byte(t, val, err)
	}

	patches := []byteprotocol.PartialPatch{
		{Pos: 1, Data: []byte(`"A"`)},
		{Pos: 3, Data: []byte(`"B"`)},
	}

	v, err := enc.EncodePartialUpdateRange(0, 4, patches)
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[1] != "A" || out[3] != "B" {
		t.Fatalf("sparse partial update failed: %v", out)
	}
}

func TestDensePartialUpdate_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	patches := make([]byteprotocol.PartialPatch, 5)
	for i := 0; i < 5; i++ {
		patches[i] = byteprotocol.PartialPatch{
			Pos:  uint32(i),
			Data: []byte(fmt.Sprintf(`%d`, i+10)),
		}
	}

	v, err := enc.EncodePartialUpdateRange(0, 4, patches)
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	for i := 0; i < 5; i++ {
		if out[i] != float64(i+10) {
			t.Fatalf("dense partial update failed: %v", out)
		}
	}
}

func TestBulkSparsePartialUpdate_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	patches := []byteprotocol.PartialPatch{
		{Pos: 100, Data: []byte(`"X"`)},
		{Pos: 105, Data: []byte(`"Y"`)},
		{Pos: 120, Data: []byte(`"Z"`)},
	}

	v, err := enc.EncodePartialUpdateRange(100, 200, patches)
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if out[100] != "X" || out[105] != "Y" || out[120] != "Z" {
		t.Fatalf("bulk sparse partial update failed: %v", out)
	}
}

func TestBulkDensePartialUpdate_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	patches := make([]byteprotocol.PartialPatch, 20)
	for i := 0; i < 20; i++ {
		patches[i] = byteprotocol.PartialPatch{
			Pos:  uint32(50 + i),
			Data: []byte(fmt.Sprintf(`"%d"`, i)),
		}
	}

	v, err := enc.EncodePartialUpdateRange(50, 69, patches)
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	for i := 0; i < 20; i++ {
		if out[50+i] != fmt.Sprintf("%d", i) {
			t.Fatalf("bulk dense partial update failed: %v", out)
		}
	}
}

func TestMax24BitDataLength_byte(t *testing.T) {
	enc := byteprotocol.Encoder{}

	size := 0xFFFFF
	big := make([]byte, size)
	for i := range big {
		big[i] = 'A'
	}

	jsonVal := append([]byte(`"`), append(big, '"')...)

	v, err := enc.EncodeInsert(0, jsonVal)
	bin := must_byte(t, v, err)
	val, err := decodeWithNode_byte(bin, []interface{}{})
	out := must_byte(t, val, err)

	if len(out[0].(string)) != size {
		t.Fatalf("expected len=%d got=%d", size, len(out[0].(string)))
	}
}
