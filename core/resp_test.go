package core

import (
	"fmt"
	"testing"
)

func TestDecodeInt64(t *testing.T) {
	cases := map[string]int64{
		":0\r\n": 0,
		":1000\r\n": 1000,
	}
	for k, v := range cases {
		elem, _ := decode([]byte(k))
		if elem != v {
			t.Fail()
		}
	}
}

func TestDecodeSimpleString(t *testing.T) {
	cases := map[string]string{
		"+PING\r\n": "PING",
	}
	for k, v := range cases {
		elem, _ := decode([]byte(k))
		if elem != v {
			t.Fail()
		}
	}
}

func TestDecodeBulkString(t *testing.T) {
	cases := map[string]string{
		"$4\r\nPING\r\n": "PING",
		"$0\r\n\r\n": "",
	}
	for k, v := range cases {
		elem, _ := decode([]byte(k))
		if elem != v {
			t.Fail()
		}
	}
}

func TestDecodeArray(t *testing.T) {
	cases := map[string][]interface{}{
		"*0\r\n": {},
		"*3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n:123\r\n": {"SET", "KEY", int64(123)},
		"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n": {[]interface{}{int64(1), int64(2), int64(3)}, []interface{}{"Hello", "World"}},
	}
	for k, v := range cases {
		elem, _ := decode([]byte(k))
		elems := elem.([]interface{})
		if len(elems) != len(v) {
			t.Fail()
		}
		for i := range elems {
			if fmt.Sprintf("%v", elems[i]) != fmt.Sprintf("%v", v[i]) {
				t.Fail()
			}
		}
	}
}

func TestDecodeError(t *testing.T) {
	cases := map[string]string{
		"-no data\r\n": "no data",
	}
	for k, v := range cases {
		elem, _ := decode([]byte(k))
		if elem != v {
			t.Fail()
		}
	}
}
