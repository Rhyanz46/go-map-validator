package map_validator

import (
	"testing"
)

func TestIsPossibleXSS(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"<script>alert('XSS')</script>", true},
		{"<img src='x' onerror='alert(1)'>", true},
		{"<div style=\"behavior:url(#default#VML);\">", true},
		{"<a href=\"javascript:alert('XSS')\">Click me</a>", true},
		{"<iframe src=\"http://example.com\"></iframe>", true},
		{"<svg onload=\"alert('XSS')\"></svg>", true},
		{"<div>Safe content</div>", false},
		{"Hello, world!", false},
		{"<b>Bold text</b>", false},
	}

	for _, test := range tests {
		result := isPossibleXSS(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', expected %v but got %v", test.input, test.expected, result)
		}
	}
}

func TestIsPossibleXSS2(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"<script>alert('XSS')</script>", true},
		{"<img src='x' onerror='alert(1)'>", true},
		{"<div style=\"behavior:url(#default#VML);\">", true},
		{"<a href=\"javascript:alert('XSS')\">Click me</a>", true},
		{"<iframe src=\"http://example.com\"></iframe>", true},
		{"<svg onload=\"alert('XSS')\"></svg>", true},
		{"<div>Safe content</div>", false},
		{"Hello, world!", false},
		{"<b>Bold text</b>", false},
	}

	for _, test := range tests {
		result := isPossibleXSS(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', expected %v but got %v", test.input, test.expected, result)
		}
	}
}

func TestNotIsPossibleXSS3(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"<script>alert('XSS')</script>", true},
		{"<img src='x' onerror='alert(1)'>", true},
		{"<div style=\"behavior:url(#default#VML);\">", true},
		{"<a href=\"javascript:alert('XSS')\">Click me</a>", true},
		{"<iframe src=\"http://example.com\"></iframe>", true},
		{"<svg onload=\"alert('XSS')\"></svg>", true},
		{"<div>Safe content</div>", false},
		{"Hello, world!", false},
		{"<b>Bold text</b>", false},
	}

	for _, test := range tests {
		result := isPossibleXSS(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', expected %v but got %v", test.input, test.expected, result)
		}
	}
}
