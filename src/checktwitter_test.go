package main

import (
	_ "fmt"
	"testing"
)

func TestFindFirstWord(t *testing.T) {
	result := FindFirstWord("@golangphilbot blue")
	if result != "blue" {
		t.Log("Expected blue but got %v", result)
		t.FailNow()
	}
}

func TestFindFirstWordSentence(t *testing.T) {
	result := FindFirstWord("@golangphilbot blue now!")
	if result != "blue" {
		t.Log("Expected blue but got %v", result)
		t.FailNow()
	}
}

func TestFindFirstWordNone(t *testing.T) {
	result := FindFirstWord("")
	if result != "" {
		t.FailNow()
	}
}

func TestFindFirstWordNotMention(t *testing.T) {
	result := FindFirstWord("Hey what are you doing?")
	if result != "" {
		t.Log("Expected empty string but got %v", result)
		t.FailNow()
	}
}
