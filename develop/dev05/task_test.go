package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGrep(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		inputFile   string
		outputFile  string
		patternWord string
	}{
		{
			name:        "Basic grep",
			args:        []string{},
			inputFile:   "testdata/inputBasic.txt",
			outputFile:  "testdata/outputBasic.txt",
			patternWord: "text",
		},
		{
			name:        "Basic grep without inputFile",
			args:        []string{},
			inputFile:   "Sometext for test without inputFile\nasl;dkjfaosdjf\nadftest\noiadhfjgopidhfg\n08374593847y5qp3948",
			outputFile:  "testdata/inputBasicWithoutOutputFile.txt",
			patternWord: "test",
		},

		{
			name:        "-A flag",
			args:        []string{"-A", "2"},
			inputFile:   "testdata/inputA.txt",
			outputFile:  "testdata/outputA.txt",
			patternWord: "text",
		},
		{
			name:        "-B flag",
			args:        []string{"-B", "2"},
			inputFile:   "testdata/inputB.txt",
			outputFile:  "testdata/outputB.txt",
			patternWord: "text",
		},
		{
			name:        "-C flag",
			args:        []string{"-C", "2"},
			inputFile:   "testdata/inputC.txt",
			outputFile:  "testdata/outputC.txt",
			patternWord: "text",
		},
		{
			name:        "-c flag",
			args:        []string{"-c"},
			inputFile:   "testdata/inputCount.txt",
			outputFile:  "testdata/outputCount.txt",
			patternWord: "text",
		},
		{
			name:        "-i flag",
			args:        []string{"-i"},
			inputFile:   "testdata/inputIgnoreCase.txt",
			outputFile:  "testdata/outputIgnoreCase.txt",
			patternWord: "text",
		},
		{
			name:        "-v flag",
			args:        []string{"-v"},
			inputFile:   "testdata/inputInvert.txt",
			outputFile:  "testdata/outputInvert.txt",
			patternWord: "text",
		},
		{
			name:        "-F flag",
			args:        []string{"-F"},
			inputFile:   "testdata/inputFixed.txt",
			outputFile:  "testdata/outputFixed.txt",
			patternWord: "exact",
		},
		{
			name:        "-n flag",
			args:        []string{"-n"},
			inputFile:   "testdata/inputLineNum.txt",
			outputFile:  "testdata/outputLineNum.txt",
			patternWord: "text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Компиляция проекта
			cmd := exec.Command("go", "build", "-o", "greptool", "task.go")
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to compile project: %v", err)
			}
			t.Cleanup(func() {
				os.Remove("greptool")
			})
			// Запуск исполняемого файла с флагами
			cmd = exec.Command("./greptool", append(tt.args, tt.patternWord, tt.inputFile)...)
			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Failed to execute greptool: %v", err)
			}

			// Чтение эталонного файла
			wantOutput, err := os.ReadFile(tt.outputFile)
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
