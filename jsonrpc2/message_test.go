package jsonrpc2

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestID_IsNull(t *testing.T) {
	tests := []struct {
		id       ID
		expected bool
	}{
		{id: ID{value: nil}, expected: true},
		{id: ID{value: "test"}, expected: false},
		{id: ID{value: 123}, expected: false},
	}

	for _, test := range tests {
		if result := test.id.IsNull(); result != test.expected {
			t.Errorf("ID(%v).IsNull() = %v; want %v", test.id, result, test.expected)
		}
	}
}

func TestID_String(t *testing.T) {
	tests := []struct {
		id       ID
		expected string
	}{
		{id: ID{value: "test"}, expected: "test"},
		{id: ID{value: 123}, expected: "123"},
		{id: ID{value: nil}, expected: ""},
	}

	for _, test := range tests {
		if result := test.id.String(); result != test.expected {
			t.Errorf("ID(%v).String() = %v; want %v", test.id, result, test.expected)
		}
	}
}

func TestID_String_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid ID type")
		}
	}()

	id := ID{value: true} // Invalid type
	_ = id.String()
}

func TestID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected ID
	}{
		{input: `"test"`, expected: ID{value: "test"}},
		{input: `123`, expected: ID{value: 123}},
		{input: `null`, expected: ID{value: nil}},
	}

	for _, test := range tests {
		var id ID
		if err := json.Unmarshal([]byte(test.input), &id); err != nil {
			t.Errorf("UnmarshalJSON(%v) error: %v", test.input, err)
		}
		if id != test.expected {
			t.Errorf("UnmarshalJSON(%v) = %v; want %v", test.input, id, test.expected)
		}
	}
}

func TestID_UnmarshalJSON_InvalidType(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name:    "boolean",
			json:    `true`,
			wantErr: true,
		},
		{
			name:    "array",
			json:    `[1,2,3]`,
			wantErr: true,
		},
		{
			name:    "object",
			json:    `{"key":"value"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id ID
			err := json.Unmarshal([]byte(tt.json), &id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestID_MarshalJSON(t *testing.T) {
	tests := []struct {
		id       ID
		expected string
	}{
		{id: ID{value: "test"}, expected: `"test"`},
		{id: ID{value: 123}, expected: `123`},
		{id: ID{value: nil}, expected: `null`},
	}

	for _, test := range tests {
		result, err := json.Marshal(test.id)
		if err != nil {
			t.Errorf("MarshalJSON(%v) error: %v", test.id, err)
		}
		if string(result) != test.expected {
			t.Errorf("MarshalJSON(%v) = %v; want %v", test.id, string(result), test.expected)
		}
	}
}

func TestID_MarshalJSON_InvalidType(t *testing.T) {
	id := ID{value: true} // Invalid type
	_, err := json.Marshal(id)
	if err == nil {
		t.Error("Expected error for invalid ID type")
	}
}

func TestRequest_MarshalJSON(t *testing.T) {
	tests := []struct {
		req      *Request[any]
		expected string
	}{
		{
			req:      &Request[any]{ID: ID{value: "1"}, Method: "testMethod", Params: map[string]any{"param1": "value1"}},
			expected: `{"jsonrpc":"2.0","id":"1","method":"testMethod","params":{"param1":"value1"}}`,
		},
		{
			req:      &Request[any]{ID: ID{value: 1}, Method: "testMethod", Params: map[string]any{"param1": "value1"}},
			expected: `{"jsonrpc":"2.0","id":1,"method":"testMethod","params":{"param1":"value1"}}`,
		},
		{
			req:      &Request[any]{ID: ID{value: nil}, Method: "testMethod", Params: map[string]any{"param1": "value1"}},
			expected: `{"jsonrpc":"2.0","method":"testMethod","params":{"param1":"value1"}}`,
		},
	}

	for _, test := range tests {
		result, err := json.Marshal(test.req)
		if err != nil {
			t.Errorf("MarshalJSON(%v) error: %v", test.req, err)
		}
		if string(result) != test.expected {
			t.Errorf("MarshalJSON(%v) = %v; want %v", test.req, string(result), test.expected)
		}
	}
}

func TestRequest_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected Request[map[string]any]
	}{
		{
			input:    `{"jsonrpc":"2.0","id":"1","method":"testMethod","params":{"param1":"value1"}}`,
			expected: Request[map[string]any]{ID: ID{value: "1"}, Method: "testMethod", Params: map[string]any{"param1": "value1"}},
		},
		{
			input:    `{"jsonrpc":"2.0","id":1,"method":"testMethod","params":{"param1":"value1"}}`,
			expected: Request[map[string]any]{ID: ID{value: 1}, Method: "testMethod", Params: map[string]any{"param1": "value1"}},
		},
		{
			input:    `{"jsonrpc":"2.0","method":"testMethod","params":{"param1":"value1"}}`,
			expected: Request[map[string]any]{ID: ID{value: nil}, Method: "testMethod", Params: map[string]any{"param1": "value1"}},
		},
	}

	for _, test := range tests {
		var req Request[map[string]any]
		if err := json.Unmarshal([]byte(test.input), &req); err != nil {
			t.Errorf("UnmarshalJSON(%v) error: %v", test.input, err)
		}
		if req.ID != test.expected.ID || req.Method != test.expected.Method || fmt.Sprintf("%v", req.Params) != fmt.Sprintf("%v", test.expected.Params) {
			t.Errorf("UnmarshalJSON(%v) = %v; want %v", test.input, req, test.expected)
		}
	}
}

func TestResponse_MarshalJSON(t *testing.T) {
	tests := []struct {
		resp     *Response[any, any]
		expected string
	}{
		{
			resp:     &Response[any, any]{ID: ID{value: "1"}, Result: "result", Error: Error[any]{}},
			expected: `{"jsonrpc":"2.0","id":"1","result":"result"}`,
		},
		{
			resp:     &Response[any, any]{ID: ID{value: 1}, Result: 123, Error: Error[any]{}},
			expected: `{"jsonrpc":"2.0","id":1,"result":123}`,
		},
		{
			resp:     &Response[any, any]{ID: ID{value: nil}, Result: nil, Error: Error[any]{Code: -32000, Message: "error"}},
			expected: `{"jsonrpc":"2.0","error":{"code":-32000,"message":"error"}}`,
		},
	}

	for _, test := range tests {
		result, err := json.Marshal(test.resp)
		if err != nil {
			t.Errorf("MarshalJSON(%v) error: %v", test.resp, err)
		}
		if string(result) != test.expected {
			t.Errorf("MarshalJSON(%v) = %v; want %v", test.resp, string(result), test.expected)
		}
	}
}

func TestResponse_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected Response[any, any]
	}{
		{
			input:    `{"jsonrpc":"2.0","id":"1","result":"result"}`,
			expected: Response[any, any]{ID: ID{value: "1"}, Result: "result", Error: Error[any]{}},
		},
		{
			input:    `{"jsonrpc":"2.0","id":1,"result":123}`,
			expected: Response[any, any]{ID: ID{value: 1}, Result: 123, Error: Error[any]{}},
		},
		{
			input:    `{"jsonrpc":"2.0","error":{"code":-32000,"message":"error"}}`,
			expected: Response[any, any]{ID: ID{value: nil}, Result: nil, Error: Error[any]{Code: -32000, Message: "error"}},
		},
	}

	for _, test := range tests {
		var resp Response[any, any]
		if err := json.Unmarshal([]byte(test.input), &resp); err != nil {
			t.Errorf("UnmarshalJSON(%v) error: %v", test.input, err)
		}
		if resp.ID != test.expected.ID || fmt.Sprintf("%v", resp.Result) != fmt.Sprintf("%v", test.expected.Result) || resp.Error.Code != test.expected.Error.Code || resp.Error.Message != test.expected.Error.Message {
			t.Errorf("UnmarshalJSON(%v) = %v; want %v", test.input, resp, test.expected)
		}
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		err      Error[any]
		expected string
	}{
		{err: Error[any]{Code: -32000, Message: "error"}, expected: "error"},
		{err: Error[any]{Code: -32001, Message: "another error"}, expected: "another error"},
	}

	for _, test := range tests {
		if result := test.err.Error(); result != test.expected {
			t.Errorf("Error(%v).Error() = %v; want %v", test.err, result, test.expected)
		}
	}
}

func TestError_code(t *testing.T) {
	tests := []struct {
		err      Error[any]
		expected int
	}{
		{err: Error[any]{Code: -32000, Message: "error"}, expected: -32000},
		{err: Error[any]{Code: -32001, Message: "another error"}, expected: -32001},
	}

	for _, test := range tests {
		if result := test.err.code(); result != test.expected {
			t.Errorf("Error(%v).code() = %v; want %v", test.err, result, test.expected)
		}
	}
}

func TestError_message(t *testing.T) {
	tests := []struct {
		err      Error[any]
		expected string
	}{
		{err: Error[any]{Code: -32000, Message: "error"}, expected: "error"},
		{err: Error[any]{Code: -32001, Message: "another error"}, expected: "another error"},
	}

	for _, test := range tests {
		if result := test.err.message(); result != test.expected {
			t.Errorf("Error(%v).message() = %v; want %v", test.err, result, test.expected)
		}
	}
}

func TestError_data(t *testing.T) {
	tests := []struct {
		err      Error[string]
		expected string
	}{
		{err: Error[string]{Code: -32000, Message: "error", Data: "data"}, expected: "data"},
		{err: Error[string]{Code: -32001, Message: "another error", Data: "more data"}, expected: "more data"},
	}

	for _, test := range tests {
		if result := test.err.data(); result != test.expected {
			t.Errorf("Error(%v).data() = %v; want %v", test.err, result, test.expected)
		}
	}
}

func TestConvertError(t *testing.T) {
	tests := []struct {
		input    error
		expected Error[any]
	}{
		{
			input:    errors.New("standard error"),
			expected: Error[any]{Code: -32000, Message: "standard error", Data: errors.New("standard error")},
		},
		{
			input:    NewError(-32001, "custom error", "data"),
			expected: Error[any]{Code: -32001, Message: "custom error", Data: "data"},
		},
	}

	for _, test := range tests {
		result := convertError(test.input)
		if result.Code != test.expected.Code || result.Message != test.expected.Message || fmt.Sprintf("%v", result.Data) != fmt.Sprintf("%v", test.expected.Data) {
			t.Errorf("convertError(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

type customError struct {
	errCode    int
	errMessage string
	errData    any
}

func (e customError) Error() string {
	return e.errMessage
}

func (e customError) code() int {
	return e.errCode
}

func (e customError) message() string {
	return e.errMessage
}

func (e customError) data() any {
	return e.errData
}

func TestConvertError_CustomError(t *testing.T) {
	err := customError{
		errCode:    -32001,
		errMessage: "custom error",
		errData:    "error data",
	}

	converted := convertError(err)
	if converted.Code != -32001 {
		t.Errorf("Expected code -32001, got %d", converted.Code)
	}
	if converted.Message != "custom error" {
		t.Errorf("Expected message 'custom error', got %s", converted.Message)
	}
	if converted.Data != "error data" {
		t.Errorf("Expected data 'error data', got %v", converted.Data)
	}
}

func TestGetMessageType(t *testing.T) {
	tests := []struct {
		input    string
		expected messageType
		err      error
	}{
		{
			input:    `{"jsonrpc":"2.0","method":"testMethod","id":1}`,
			expected: messageRequest,
			err:      nil,
		},
		{
			input:    `{"jsonrpc":"2.0","method":"testMethod"}`,
			expected: messageNotification,
			err:      nil,
		},
		{
			input:    `{"jsonrpc":"2.0","id":1,"result":"testResult"}`,
			expected: messageResponse,
			err:      nil,
		},
		{
			input:    `{"jsonrpc":"2.0","error":{"code":-32000,"message":"error"}}`,
			expected: messageResponse,
			err:      nil,
		},
		{
			input:    `{"jsonrpc":"2.0","result":"testResult"}`,
			expected: 0,
			err:      errors.New("invalid message type"),
		},
		{
			input:    `{"jsonrpc":"1.0","method":"testMethod"}`,
			expected: 0,
			err:      errors.New("invalid JSON-RPC version"),
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			var msg json.RawMessage = []byte(test.input)
			result, err := getMessageType(msg)
			if result != test.expected || (err != nil && err.Error() != test.err.Error()) {
				t.Errorf("getMessageType(%v) = %v, %v; want %v, %v", test.input, result, err, test.expected, test.err)
			}
		})
	}
}
