package encoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"testing"

	"github.com/milton-alvarenga/goreactivehtml/internal/server/encode/protocol"
)

//
// --- Test Helpers ---
//

func decodeWithNode(input []byte) ([]interface{}, error) {
	cmd := exec.Command("nodejs", "../decoder/node.js")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	// Capture stdout and stderr
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	if _, err := stdin.Write(input); err != nil {
		return nil, err
	}
	stdin.Close()

	// Wait for the process to finish
	if err := cmd.Wait(); err != nil {
		// Log the error message from stderr
		log.Println("Node.js stderr:", stderr.String())
		return nil, err
	}

	// Log the stdout for debugging purposes
	log.Println("Node.js stdout:", stdout.String())

	var result []interface{}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func must[T any](t *testing.T, v T, err error) T {
	if err != nil {
		t.Fatal(err)
	}
	return v
}

//
// --- Tests ---
//

func TestInsertSingle(t *testing.T) {
	enc := protocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`"A"`))
	t.Log(v)
	bin := must(t, v, err)
	t.Log(bin)
	t.Logf("Encoded binary data: %v", bin)
	val, err := decodeWithNode(bin)
	t.Log(val)
	out := must(t, val, err)
	t.Log(out)

	if out[0] != "A" {
		t.Fatalf("expected A, got %v", out)
	}
}

func TestEncodeInsert(t *testing.T) {
	enc := protocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`"hello"`))
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[0] != "hello" {
		t.Fatalf("expected hello, got %v", out)
	}
}

func TestUpdateSingle(t *testing.T) {
	enc := protocol.Encoder{}
	v, err := enc.EncodeUpdate(5, []byte(`123`))
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[5] != float64(123) {
		t.Fatalf("expected 123, got %v", out)
	}
}

func TestEncodeUpdate(t *testing.T) {
	enc := protocol.Encoder{}
	v, err := enc.EncodeUpdate(0, []byte(`123`))
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[0] != float64(123) {
		t.Fatalf("expected 123, got %v", out)
	}
}

func TestPartialUpdate(t *testing.T) {
	enc := protocol.Encoder{}
	v, err := enc.EncodePartialUpdate(5, []byte(`"patched"`))
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[5] != "patched" {
		t.Fatalf("expected patched, got %v", out)
	}
}

func TestBulkPartialUpdate(t *testing.T) {
	enc := protocol.Encoder{}

	patches := []protocol.PartialPatch{
		{Pos: 10, Data: []byte(`"A"`)},
		{Pos: 12, Data: []byte(`"B"`)},
	}

	v, err := enc.EncodePartialUpdateRange(10, 20, patches)
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[10] != "A" || out[12] != "B" {
		t.Fatalf("expected sparse patches, got %v", out)
	}
}

func TestDeleteSingle(t *testing.T) {
	enc := protocol.Encoder{}

	// Start with empty array
	fullArray := [][]byte{[]byte(`1`)}

	v, err := enc.EncodeUpdateRange(0, uint32(len(fullArray)-1), fullArray)
	val, err := decodeWithNode(must(t, v, err))
	out := must(t, val, err)
	if len(out) != 1 {
		t.Fatalf("delete failed one element: %v", out)
	}

	// Add second element
	fullArray = append(fullArray, []byte(`2`))
	v, err = enc.EncodeUpdateRange(0, uint32(len(fullArray)-1), fullArray)
	val, err = decodeWithNode(must(t, v, err))
	out = must(t, val, err)
	if len(out) != 2 {
		t.Fatalf("delete failed two element: %v", out)
	}

	// Add third element
	fullArray = append(fullArray, []byte(`3`))
	v, err = enc.EncodeUpdateRange(0, uint32(len(fullArray)-1), fullArray)
	val, err = decodeWithNode(must(t, v, err))
	out = must(t, val, err)
	if len(out) != 3 {
		t.Fatalf("delete failed three element: %v", out)
	}

	// Now delete second element
	fullArray = append(fullArray[:1], fullArray[2:]...) // remove index 1
	v, err = enc.EncodeUpdateRange(0, uint32(len(fullArray)-1), fullArray)
	val, err = decodeWithNode(must(t, v, err))
	out = must(t, val, err)

	if len(out) != 2 || int(out[1].(float64)) != 3 {
		t.Fatalf("delete failed: %v", out)
	}
}

// NOK
func TestDeleteSingleStateFul(t *testing.T) {
	enc := protocol.Encoder{}

	v, err := enc.EncodeInsert(0, []byte(`1`))
	val, err := decodeWithNode(must(t, v, err))
	must(t, val, err)

	v, err = enc.EncodeInsert(1, []byte(`2`))
	val, err = decodeWithNode(must(t, v, err))
	must(t, val, err)

	v, err = enc.EncodeInsert(2, []byte(`3`))
	val, err = decodeWithNode(must(t, v, err))
	must(t, val, err)

	v, err = enc.EncodeDelete(1)
	val, err = decodeWithNode(must(t, v, err))
	out := must(t, val, err)

	if len(out) != 2 || int(out[1].(float64)) != 3 {
		t.Fatalf("delete failed: %v", out)
	}
}

func TestInsertRange(t *testing.T) {
	enc := protocol.Encoder{}

	payloads := [][]byte{
		[]byte(`"A"`),
		[]byte(`"B"`),
		[]byte(`"C"`),
	}

	v, err := enc.EncodeInsertRange(10, 12, payloads)
	val, err := decodeWithNode(must(t, v, err))
	out := must(t, val, err)

	if out[10] != "A" || out[11] != "B" || out[12] != "C" {
		t.Fatalf("range insert failed: %v", out)
	}
}

func TestUpdateRange(t *testing.T) {
	enc := protocol.Encoder{}

	payloads := [][]byte{
		[]byte(`10`),
		[]byte(`20`),
		[]byte(`30`),
	}

	v, err := enc.EncodeUpdateRange(3, 5, payloads)
	val, err := decodeWithNode(must(t, v, err))
	out := must(t, val, err)

	if out[3] != float64(10) ||
		out[4] != float64(20) ||
		out[5] != float64(30) {
		t.Fatalf("range update failed: %v", out)
	}
}

// NOK
func TestDeleteRange(t *testing.T) {
	enc := protocol.Encoder{}

	for i := 0; i < 10; i++ {
		v, err := enc.EncodeInsert(uint32(i), []byte(`1`))
		val, err := decodeWithNode(must(t, v, err))
		must(t, val, err)
	}

	v, err := enc.EncodeDeleteRange(3, 6)
	val, err := decodeWithNode(must(t, v, err))
	out := must(t, val, err)

	if len(out) != 6 {
		t.Fatalf("expected 6 elements remain, got %v", out)
	}
}

// NOK
func TestInsertAtEnd(t *testing.T) {
	enc := protocol.Encoder{}

	v, err := enc.EncodeInsert(0, []byte(`0`))
	val, err := decodeWithNode(must(t, v, err))
	must(t, val, err)

	v, err = enc.EncodeInsert(1, []byte(`1`))
	val, err = decodeWithNode(must(t, v, err))
	must(t, val, err)

	v, err = enc.EncodeInsert(2, []byte(`2`))
	val, err = decodeWithNode(must(t, v, err))
	must(t, val, err)

	v, err = enc.EncodeInsert(3, []byte(`"END"`))
	val, err = decodeWithNode(must(t, v, err))
	out := must(t, val, err)

	if out[3] != "END" {
		t.Fatalf("insert at end failed: %v", out)
	}
}

// NOK
func TestInsertBeyondEnd(t *testing.T) {
	enc := protocol.Encoder{}
	v, err := enc.EncodeInsert(10, []byte(`"X"`))
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[10] != "X" {
		t.Fatalf("expected X at pos 10, got %v", out)
	}

	for i := 0; i < 10; i++ {
		if out[i] != nil {
			t.Fatalf("expected null at %d got %v", i, out[i])
		}
	}
}

func TestUpdateBeyondEnd(t *testing.T) {
	enc := protocol.Encoder{}
	v, err := enc.EncodeUpdate(15, []byte(`999`))
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[15] != float64(999) {
		t.Fatalf("expected 999, got %v", out)
	}
}

func TestZeroLengthPayload(t *testing.T) {
	enc := protocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`""`))
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[0] != "" {
		t.Fatalf("expected empty string, got %v", out)
	}
}

func TestJSONString(t *testing.T) {
	enc := protocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`"hello world"`))
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[0] != "hello world" {
		t.Fatalf("expected hello world, got %v", out)
	}
}

func TestJSONNumber(t *testing.T) {
	enc := protocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`42`))
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[0] != float64(42) {
		t.Fatalf("expected 42, got %v", out)
	}
}

func TestJSONObject(t *testing.T) {
	enc := protocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`{"a":1,"b":2}`))
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	obj := out[0].(map[string]interface{})
	if obj["a"] != float64(1) || obj["b"] != float64(2) {
		t.Fatalf("object mismatch: %v", obj)
	}
}

func TestJSONArray(t *testing.T) {
	enc := protocol.Encoder{}
	v, err := enc.EncodeInsert(0, []byte(`[1,2,3]`))
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	arr := out[0].([]interface{})
	if len(arr) != 3 || arr[1] != float64(2) {
		t.Fatalf("array mismatch: %v", arr)
	}
}

// NOK
func TestMax24BitPosition(t *testing.T) {
	enc := protocol.Encoder{}
	pos := uint32(0xFFFFFF)

	v, err := enc.EncodeInsert(pos, []byte(`"MAXPOS"`))
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[pos] != "MAXPOS" {
		t.Fatalf("expected MAXPOS at %d", pos)
	}
}

func TestMax24BitDataLength(t *testing.T) {
	enc := protocol.Encoder{}

	size := 0xFFFFF
	big := make([]byte, size)
	for i := range big {
		big[i] = 'A'
	}

	jsonVal := append([]byte(`"`), append(big, '"')...)

	v, err := enc.EncodeInsert(0, jsonVal)
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if len(out[0].(string)) != size {
		t.Fatalf("expected len=%d got=%d", size, len(out[0].(string)))
	}
}

func TestBulkInsertMixedSize(t *testing.T) {
	enc := protocol.Encoder{}

	payloads := [][]byte{
		[]byte(`"A"`),
		[]byte(`"BBBBBB"`),
		[]byte(`"C"`),
		[]byte(`"DDDDDDDDDDDDDDDDDDD"`),
	}

	v, err := enc.EncodeInsertRange(5, 8, payloads)
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[5] != "A" ||
		out[6] != "BBBBBB" ||
		out[7] != "C" ||
		out[8] != "DDDDDDDDDDDDDDDDDDD" {
		t.Fatalf("bulk insert mismatch: %v", out)
	}
}

// NOK
func TestRandomFixedCases(t *testing.T) {
	enc := protocol.Encoder{}

	for i := 0; i < 200; i++ {
		pos := uint32(i % 50)
		valStr := fmt.Sprintf("%d", i)

		v, err := enc.EncodeInsert(pos, []byte(valStr))
		bin := must(t, v, err)
		val, err := decodeWithNode(bin)
		out := must(t, val, err)

		if out[pos] != valStr {
			t.Fatalf("random test failed at pos=%d: %v", pos, out)
		}
	}
}

func TestSparsePartialUpdate(t *testing.T) {
	enc := protocol.Encoder{}

	for i := 0; i < 5; i++ {
		v, err := enc.EncodeInsert(uint32(i), []byte(fmt.Sprintf(`%d`, i)))
		val, err := decodeWithNode(must(t, v, err))
		must(t, val, err)
	}

	patches := []protocol.PartialPatch{
		{Pos: 1, Data: []byte(`"A"`)},
		{Pos: 3, Data: []byte(`"B"`)},
	}

	v, err := enc.EncodePartialUpdateRange(0, 4, patches)
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[1] != "A" || out[3] != "B" {
		t.Fatalf("sparse partial update failed: %v", out)
	}
}

func TestDensePartialUpdate(t *testing.T) {
	enc := protocol.Encoder{}

	patches := make([]protocol.PartialPatch, 5)
	for i := 0; i < 5; i++ {
		patches[i] = protocol.PartialPatch{
			Pos:  uint32(i),
			Data: []byte(fmt.Sprintf(`%d`, i+10)),
		}
	}

	v, err := enc.EncodePartialUpdateRange(0, 4, patches)
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	for i := 0; i < 5; i++ {
		if out[i] != float64(i+10) {
			t.Fatalf("dense partial update failed: %v", out)
		}
	}
}

func TestBulkSparsePartialUpdate(t *testing.T) {
	enc := protocol.Encoder{}

	patches := []protocol.PartialPatch{
		{Pos: 100, Data: []byte(`"X"`)},
		{Pos: 105, Data: []byte(`"Y"`)},
		{Pos: 120, Data: []byte(`"Z"`)},
	}

	v, err := enc.EncodePartialUpdateRange(100, 200, patches)
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	if out[100] != "X" || out[105] != "Y" || out[120] != "Z" {
		t.Fatalf("bulk sparse partial update failed: %v", out)
	}
}

func TestBulkDensePartialUpdate(t *testing.T) {
	enc := protocol.Encoder{}

	patches := make([]protocol.PartialPatch, 20)
	for i := 0; i < 20; i++ {
		patches[i] = protocol.PartialPatch{
			Pos:  uint32(50 + i),
			Data: []byte(fmt.Sprintf(`"%d"`, i)),
		}
	}

	v, err := enc.EncodePartialUpdateRange(50, 69, patches)
	bin := must(t, v, err)
	val, err := decodeWithNode(bin)
	out := must(t, val, err)

	for i := 0; i < 20; i++ {
		if out[50+i] != fmt.Sprintf("%d", i) {
			t.Fatalf("bulk dense partial update failed: %v", out)
		}
	}
}
