package tinyjson_test

import (
	"reflect"
	"testing"

	"github.com/cdvelop/tinyjson"
)

type TestStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func EncodeShared(t *testing.T, j *tinyjson.TinyJSON) {
	t.Run("Encode String", func(t *testing.T) {
		input := "hello"
		expected := `"hello"`
		result, err := j.Encode(input)
		if err != nil {
			t.Fatalf("Encode failed: %v", err)
		}
		if string(result) != expected {
			t.Errorf("Expected %s, got %s", expected, string(result))
		}
	})

	t.Run("Encode Int", func(t *testing.T) {
		input := 123
		expected := "123"
		result, err := j.Encode(input)
		if err != nil {
			t.Fatalf("Encode failed: %v", err)
		}
		if string(result) != expected {
			t.Errorf("Expected %s, got %s", expected, string(result))
		}
	})

	t.Run("Encode Struct", func(t *testing.T) {
		input := TestStruct{Name: "Alice", Age: 30}
		// JSON key order is not guaranteed, so we might need to check fields or use a more robust comparison if strict string match fails often.
		// For simple structs in standard json, it's usually consistent, but let's see.
		// Actually, for this simple case, let's just check if it contains the keys and values.
		result, err := j.Encode(input)
		if err != nil {
			t.Fatalf("Encode failed: %v", err)
		}
		resStr := string(result)
		if resStr != `{"name":"Alice","age":30}` && resStr != `{"age":30,"name":"Alice"}` {
			t.Errorf("Expected JSON representation of struct, got %s", resStr)
		}
	})

	t.Run("Encode Slice of Structs", func(t *testing.T) {
		input := []TestStruct{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
		}
		result, err := j.Encode(input)
		if err != nil {
			t.Fatalf("Encode failed: %v", err)
		}
		resStr := string(result)
		t.Logf("Encoded slice of structs: %s", resStr)

		// Should be a JSON array, not an empty string
		if resStr == `""` || resStr == "" {
			t.Errorf("BUG: Slice of structs encoded as empty string instead of JSON array, got: %s", resStr)
		}

		// Verify it's a valid JSON array
		if len(resStr) < 2 || resStr[0] != '[' || resStr[len(resStr)-1] != ']' {
			t.Errorf("Expected JSON array format [...], got: %s", resStr)
		}
	})
}

func DecodeShared(t *testing.T, j *tinyjson.TinyJSON) {
	t.Run("Decode String", func(t *testing.T) {
		input := `"world"`
		var result string
		err := j.Decode([]byte(input), &result)
		if err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		if result != "world" {
			t.Errorf("Expected 'world', got '%s'", result)
		}
	})

	t.Run("Decode Int", func(t *testing.T) {
		input := "456"
		var result int
		err := j.Decode([]byte(input), &result)
		if err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		if result != 456 {
			t.Errorf("Expected 456, got %d", result)
		}
	})

	t.Run("Decode Struct", func(t *testing.T) {
		input := `{"name":"Bob","age":25}`
		var result TestStruct
		err := j.Decode([]byte(input), &result)
		if err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		expected := TestStruct{Name: "Bob", Age: 25}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("Decode Slice of Structs", func(t *testing.T) {
		input := `[{"name":"Alice","age":30},{"name":"Bob","age":25}]`
		var result []TestStruct
		err := j.Decode([]byte(input), &result)
		if err != nil {
			t.Fatalf("Decode failed: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("Expected 2 structs, got %d", len(result))
		}

		if len(result) > 0 && result[0].Name != "Alice" {
			t.Errorf("Expected first name 'Alice', got '%s'", result[0].Name)
		}

		if len(result) > 1 && result[1].Name != "Bob" {
			t.Errorf("Expected second name 'Bob', got '%s'", result[1].Name)
		}
	})

	// Test case that replicates crudp Packet structure with [][]byte
	t.Run("Decode Struct with byte and [][]byte fields", func(t *testing.T) {
		// This replicates crudp.Packet structure
		type Packet struct {
			Action    byte     `json:"action"`
			HandlerID uint8    `json:"handler_id"`
			ReqID     string   `json:"req_id"`
			Data      [][]byte `json:"data"`
		}

		// First encode a packet
		innerData := []byte(`{"name":"John"}`)
		packet := Packet{
			Action:    'c',
			HandlerID: 0,
			ReqID:     "test-1",
			Data:      [][]byte{innerData},
		}

		encoded, err := j.Encode(packet)
		if err != nil {
			t.Fatalf("Failed to encode packet: %v", err)
		}
		t.Logf("Encoded packet: %s", string(encoded))

		// Now decode it back
		var decoded Packet
		err = j.Decode(encoded, &decoded)
		if err != nil {
			t.Fatalf("Failed to decode packet: %v", err)
		}

		if decoded.Action != 'c' {
			t.Errorf("Expected Action 'c' (%d), got %d", 'c', decoded.Action)
		}
		if decoded.HandlerID != 0 {
			t.Errorf("Expected HandlerID 0, got %d", decoded.HandlerID)
		}
		if decoded.ReqID != "test-1" {
			t.Errorf("Expected ReqID 'test-1', got '%s'", decoded.ReqID)
		}
		if len(decoded.Data) != 1 {
			t.Fatalf("Expected 1 data item, got %d", len(decoded.Data))
		}
		if string(decoded.Data[0]) != string(innerData) {
			t.Errorf("Expected data '%s', got '%s'", string(innerData), string(decoded.Data[0]))
		}
	})

	// Test BatchRequest structure (nested structs with [][]byte)
	t.Run("Decode BatchRequest with nested Packets", func(t *testing.T) {
		type Packet struct {
			Action    byte     `json:"action"`
			HandlerID uint8    `json:"handler_id"`
			ReqID     string   `json:"req_id"`
			Data      [][]byte `json:"data"`
		}
		type BatchRequest struct {
			Packets []Packet `json:"packets"`
		}

		innerData := []byte(`{"name":"John","email":"john@example.com"}`)
		batch := BatchRequest{
			Packets: []Packet{
				{
					Action:    'c',
					HandlerID: 0,
					ReqID:     "test-create",
					Data:      [][]byte{innerData},
				},
			},
		}

		encoded, err := j.Encode(batch)
		if err != nil {
			t.Fatalf("Failed to encode batch: %v", err)
		}
		t.Logf("Encoded batch: %s", string(encoded))

		var decoded BatchRequest
		err = j.Decode(encoded, &decoded)
		if err != nil {
			t.Fatalf("Failed to decode batch: %v", err)
		}

		if len(decoded.Packets) != 1 {
			t.Fatalf("Expected 1 packet, got %d", len(decoded.Packets))
		}

		pkt := decoded.Packets[0]
		if pkt.Action != 'c' {
			t.Errorf("Expected Action 'c' (%d), got %d", 'c', pkt.Action)
		}
		if pkt.ReqID != "test-create" {
			t.Errorf("Expected ReqID 'test-create', got '%s'", pkt.ReqID)
		}
		if len(pkt.Data) != 1 {
			t.Fatalf("Expected 1 data item, got %d", len(pkt.Data))
		}
		if string(pkt.Data[0]) != string(innerData) {
			t.Errorf("Expected data '%s', got '%s'", string(innerData), string(pkt.Data[0]))
		}
	})
}
