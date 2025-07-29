package main

import (
	"reflect"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "remove https scheme",
			input:    "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove trailing slash",
			input:    "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "http instead of https",
			input:    "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "only host",
			input:    "https://blog.boot.dev/",
			expected: "blog.boot.dev",
		},
	}
	for _, tc := range tests {
		actual, err := normalizeURL(tc.input)
		if err != nil {
			t.Errorf("Tests '%s' failed: unexpected error: %v", tc.name, err)
		}
		if actual != tc.expected {
			t.Errorf("Test '%s' failed: expected %v, got %v", tc.name, tc.expected, actual)
		}
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputBody string
		inputURL  string
		expected  []string
	}{
		{
			name:      "absolute URL",
			inputBody: `<html><body><a href="https://blog.boot.dev"><span>Boot.dev</span></a></body></html>`,
			inputURL:  "https://blog.boot.dev",
			expected:  []string{"https://blog.boot.dev"},
		},
		{
			name:      "relative URL",
			inputBody: `<html><body><a href="/path/one"><span>Boot.dev</span></a></body></html>`,
			inputURL:  "https://blog.boot.dev",
			expected:  []string{"https://blog.boot.dev/path/one"},
		},
		{
			name: "multiple URL",
			inputBody: `<html><body>
        <a href="https://boot.dev">Boot.dev</a>
        <a href="/learn">Learn</a>
        <a href="/about">About</a>
    </body></html>`,
			inputURL: "https://blog.boot.dev",
			expected: []string{"https://boot.dev", "https://blog.boot.dev/learn", "https://blog.boot.dev/about"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test '%s' failed: unexpected error: %v", tc.name, err)
				return
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test '%s' failed: expected %v, got %v", tc.name, tc.expected, actual)
			}
		})
	}
}
