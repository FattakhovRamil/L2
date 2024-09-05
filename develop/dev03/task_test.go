package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		inputFile string
		wantFile  string
	}{
		{
			name:      "Basic sort",
			args:      []string{},
			inputFile: "testdata/input_basic.txt",
			wantFile:  "testdata/output_basic.txt",
		},
		{
			name:      "-b -u flag",
			args:      []string{"-b", "-u"},
			inputFile: "testdata/input_ignore_spaces.txt",
			wantFile:  "testdata/output_ignore_spaces.txt",
		},
		{
			name:      "-n -u flag",
			args:      []string{"-n", "-u"},
			inputFile: "testdata/input_numeric.txt",
			wantFile:  "testdata/output_numeric.txt",
		},
		{
			name:      "-n -u flag",
			args:      []string{"-n", "-u"},
			inputFile: "testdata/input_numeric.txt",
			wantFile:  "testdata/output_numeric.txt",
		},
		{
			name:      "-k -r flag",
			args:      []string{"-k", "2", "-r"},
			inputFile: "testdata/input_column.txt",
			wantFile:  "testdata/output_column.txt",
		},
		{
			name:      "-M flag",
			args:      []string{"-M"},
			inputFile: "testdata/input_month.txt",
			wantFile:  "testdata/output_month.txt",
		},
		{
			name:      "-c flag",
			args:      []string{"-c"},
			inputFile: "testdata/input_check_sorted.txt",
			wantFile:  "testdata/output_check_sorted.txt",
		},
		{
			name:      "-h flag",
			args:      []string{"-h"},
			inputFile: "testdata/input_human_numeric.txt",
			wantFile:  "testdata/output_human_numeric.txt",
		},
		// Добавьте больше тестов для других флагов и комбинаций
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Компиляция проекта
			cmd := exec.Command("go", "build", "-o", "sorttool", "task.go")
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to compile project: %v", err)
			}
			t.Cleanup(func() {
				os.Remove("sorttool")
			})
			// Запуск исполняемого файла с флагами
			cmd = exec.Command("./sorttool", append(tt.args, tt.inputFile)...)
			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Failed to execute sorttool: %v", err)
			}

			// Чтение эталонного файла
			wantOutput, err := os.ReadFile(tt.wantFile)
			if err != nil {
				t.Fatalf("Failed to read want file: %v", err)
			}

			// Сравнение результатов
			if strings.TrimSpace(string(output)) != strings.TrimSpace(string(wantOutput)) {
				t.Errorf("Result mismatch for %s:\n got: %v\n want: %v", tt.name, string(output), string(wantOutput))
			}
		})
	}
}
