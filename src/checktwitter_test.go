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

func TestBadWord(t *testing.T) {
	tweet := "@golangphilbot 150"

	rgb := reRGB.FindStringSubmatch(tweet)
	if rgb != nil {
		t.FailNow()
	}

	result := FindFirstWord(tweet)
	if result == "" {
		t.Log("Expected 150 but got %v", result)
		t.FailNow()
	}

	value := colorMap[result]
	if value != "" {
		t.FailNow()
	}
}
