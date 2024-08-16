package netstat

import (
	"reflect"
	"testing"
)

func TestConvertIp(t *testing.T) {
	ip := "c0a80001" // 192.168.0.1 in hexadecimal
	expected := "1.0.168.192"
	result := convertIp(ip)
	if result != expected {
		t.Errorf("convertIp(%s) = %s; want %s", ip, result, expected)
	}
}
func TestHexToDec(t *testing.T) {
	hex := "1F"
	expected := int64(31)
	result := hexToDec(hex)
	if result != expected {
		t.Errorf("hexToDec(%s) = %d; want %d", hex, result, expected)
	}
}
func TestGetProcessName(t *testing.T) {
	exe := "/usr/bin/example"
	expected := "Example"
	result := getProcessName(exe)
	if result != expected {
		t.Errorf("getProcessName(%s) = %s; want %s", exe, result, expected)
	}
}

func TestRemoveEmpty(t *testing.T) {
	input := []string{"", "hello", "", "world", ""}
	expected := []string{"hello", "world"}

	result := removeEmpty(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("removeEmpty(%v) = %v; want %v", input, result, expected)
	}
}

func TestProcessNetstatLine(t *testing.T) {
	line := "   0: 0100007F:B1B2 0100007F:9C40 01 00000000:00000000 02:000003F4 00000000     0        0 335688 4 00000000810d32b5 20 4 27 10 -1                   "
	expected := Process{
		User:        "root",
		Name:        "",
		Pid:         "",
		Exe:         "",
		State:       "ESTABLISHED",
		Ip:          "127.0.0.1",
		Port:        45490,
		ForeignIp:   "127.0.0.1",
		ForeignPort: 40000,
	}

	output := make(chan Process, 1)
	go processNetstatLine(line, &[]iNode{}, output)
	result := <-output

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("processNetstatLine(%s) = %v; want %v", line, result, expected)
	}
}
func TestContains(t *testing.T) {
	s := []string{"apple", "banana", "orange"}

	// Test case: element exists in the slice
	if !contains(s, "banana") {
		t.Errorf("contains(%v, %s) = false; want true", s, "banana")
	}

	// Test case: element does not exist in the slice
	if contains(s, "grape") {
		t.Errorf("contains(%v, %s) = true; want false", s, "grape")
	}

	// Test case: empty slice
	if contains([]string{}, "apple") {
		t.Errorf("contains(%v, %s) = true; want false", []string{}, "apple")
	}
}
