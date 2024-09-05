package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestAnagram(t *testing.T) {
	tests := []struct {
		input          []string
		expectedOutput map[string][]string
		expectedError  error
	}{
		{
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"},
			expectedOutput: map[string][]string{
				"пятак":  {"пятка", "тяпка"},
				"листок": {"слиток", "столик"},
			},
			expectedError: nil,
		},
		{
			input:          []string{},
			expectedOutput: map[string][]string{},
			expectedError:  errors.New("пустая строка"),
		},
		{
			input: []string{"стол", "толс", "тсло"},
			expectedOutput: map[string][]string{
				"стол": {"толс", "тсло"},
			},
			expectedError: nil,
		},
		{
			input: []string{"пятак", "тяпка"},
			expectedOutput: map[string][]string{
				"пятак": {"тяпка"},
			},
			expectedError: nil,
		},
		{
			input: []string{"apple", "banana", "Apple", "BANANA"},
			expectedOutput: map[string][]string{
				"apple":  {"apple"},
				"banana": {"banana"},
			},
			expectedError: nil,
		},
		{
			input:          []string{"ааа", "аа", "а"},
			expectedOutput: map[string][]string{},
			expectedError:  nil,
		},
		{
			input: []string{"aaa", "aa", "a", "aaa"},
			expectedOutput: map[string][]string{
				"aaa": {"aaa"},
			},
			expectedError: nil,
		},
		{
			input: []string{"Кот", "Ток", "кот", "ток"},
			expectedOutput: map[string][]string{
				"кот": {"кот", "ток", "ток"},
			},
			expectedError: nil,
		},
		{
			input:          []string{"авто", "товар", "втор"},
			expectedOutput: map[string][]string{},
			expectedError:  nil,
		},
		{
			input: []string{"test", "sett", "tset", "tets"},
			expectedOutput: map[string][]string{
				"sett": {"test", "tets", "tset"},
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		output, err := anagram(test.input)
		if !reflect.DeepEqual(err, test.expectedError) {
			t.Errorf("anagram(%v) returned error %v, expected %v", test.input, err, test.expectedError)
		}
		if !reflect.DeepEqual(output, test.expectedOutput) {
			if len(output) == 0 && len(test.expectedOutput) == 0 {
				continue
			}
			t.Errorf("anagram(%v) returned %v, expected %v", test.input, output, test.expectedOutput)
		}
	}
}
