package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCut(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		inputData  string
		wantOutput string
	}{
		{
			name:       "Basic cut with tab delimiter",
			args:       []string{"-f", "1,3"},
			inputData:  "one\ttwo\tthree\nfour\tfive\tsix",
			wantOutput: "one\tthree\nfour\tsix",
		},
		{
			name:       "Cut with comma delimiter",
			args:       []string{"-f", "2", "-d", ","},
			inputData:  "one,two,three\nfour,five,six",
			wantOutput: "two\nfive",
		},
		{
			name:       "Cut with separator flag",
			args:       []string{"-f", "2", "-s"},
			inputData:  "one\ttwo\nthree\tfour\nno delimiter here",
			wantOutput: "two\nfour",
		},
		{
			name:       "Cut with non-default delimiter and fields",
			args:       []string{"-f", "1,3", "-d", ";"},
			inputData:  "one;two;three\nfour;five;six",
			wantOutput: "one;three\nfour;six",
		},
		{
			name:       "Cut with non-default delimiter, fields, and separator flag",
			args:       []string{"-f", "2", "-d", ",", "-s"},
			inputData:  "one,two,three\nfour,five,six\nno delimiter here",
			wantOutput: "two\nfive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создание временного файла для ввода
			inputFile, err := os.CreateTemp("", "input")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(inputFile.Name())

			_, err = inputFile.WriteString(tt.inputData)
			if err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			inputFile.Close()

			// Компиляция проекта
			cmd := exec.Command("go", "build", "-o", "cuttool", "task.go")
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to compile project: %v", err)
			}
			t.Cleanup(func() {
				os.Remove("cuttool")
			})

			// Запуск исполняемого файла с флагами
			cmd = exec.Command("./cuttool", tt.args...)
			cmd.Stdin, err = os.Open(inputFile.Name())
			if err != nil {
				t.Fatalf("Failed to open temp file: %v", err)
			}
			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Failed to execute cuttool: %v", err)
			}

			// Сравнение результатов
			if strings.TrimSpace(string(output)) != strings.TrimSpace(tt.wantOutput) {
				t.Errorf("Result mismatch for %s:\n got: %v\n want: %v", tt.name, string(output), tt.wantOutput)
			}
		})
	}
}
