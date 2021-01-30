package vlog

import "testing"

type unsafeStruct struct {
	Name string
}

type safeStruct struct {
	Name   string
	Unsafe string
}

func (s safeStruct) SafeString() string {
	return s.Name
}

func TestBasicJoining(t *testing.T) {
	result := redactAndJoinInterfaces("hello", 23, 12.56, []byte("world"), unsafeStruct{Name: "Jeremy"})

	if result != "hello 23 12.56 [redacted []uint8] [redacted vlog.unsafeStruct]" {
		t.Error("unexpected output: ", result)
	}

	result2 := redactAndJoinInterfaces("hello", 23, 12.56, []byte("world"), safeStruct{Name: "Anthony", Unsafe: "bad data!!"})

	if result2 != "hello 23 12.56 [redacted []uint8] Anthony" {
		t.Error("unexpected output: ", result2)
	}
}
