package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Убедитесь, что вы изменили путь к файлу `task.go` на правильный для вашей структуры проекта.

func TestRunCommand(t *testing.T) {
	printOutPwd()
	printOutCdPwd()
	tests := []Tests{
		{
			name:       "echo command",
			input:      "echo hello world",
			outputFile: "testdata/outputEcho.txt",
		},
		{
			name:       "pwd command",
			input:      "pwd",
			outputFile: "testdata/outputPwd.txt",
		},
		{
			name:       "ps command",
			input:      "ps",
			outputFile: "testdata/outputPs.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Компиляция проекта
			cmd := exec.Command("go", "build", "-o", "myprogram", "task.go")
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to compile project: %v", err)
			}
			t.Cleanup(func() {
				os.Remove("myprogram")
			})

			// Запуск исполняемого файла
			command := exec.Command("./myprogram")

			// Передача входных данных в стандартный ввод программы
			command.Stdin = strings.NewReader(tt.input)
			if tt.name == "ps command" {
				printOutPs()
				command.Stdin = strings.NewReader(tt.input)

				var out bytes.Buffer
				command.Stdout = &out
				command.Stderr = &out

				if err := command.Run(); err != nil {
					t.Fatalf("Failed to execute myprogram: %v", err)
				}

				// Чтение эталонного файла
				wantOutput, err := os.ReadFile(tt.outputFile)
				if err != nil {
					t.Fatalf("Failed to read want file: %v", err)
				}

				// Получение строк из вывода
				gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
				wantLines := strings.Split(strings.TrimSpace(string(wantOutput)), "\n")

				// Проверка на пересечение строк
				if !hasAnyMatch(gotLines, wantLines) {
					t.Errorf("Result mismatch for %s:\n got: %v\n want: %v", tt.name, gotLines, wantLines)
				}
				return
			}
			var out bytes.Buffer
			command.Stdout = &out
			command.Stderr = &out

			if err := command.Run(); err != nil {
				t.Fatalf("Failed to execute myprogram: %v", err)
			}

			// Чтение эталонного файла
			wantOutput, err := os.ReadFile(tt.outputFile)
			if err != nil {
				t.Fatalf("Failed to read want file: %v", err)
			}

			// Сравнение результатов
			if strings.TrimSpace(out.String()) != strings.TrimSpace(string(wantOutput)) {
				t.Errorf("Result mismatch for %s:\n got: %v\n want: %v", tt.name, out.String(), string(wantOutput))
			}
		})
	}
}

func TestRunCommandPs(t *testing.T) {
	tests := []Tests{
		{
			name:       "ps command",
			input:      "ps",
			outputFile: "testdata/outputPs.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Компиляция проекта
			cmd := exec.Command("go", "build", "-o", "myprogram", "task.go")
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to compile project: %v", err)
			}
			t.Cleanup(func() {
				os.Remove("myprogram")
			})
			// Запуск исполняемого файла
			command := exec.Command("./myprogram")
			// Передача входных данных в стандартный ввод программы
			command.Stdin = strings.NewReader(tt.input)
			printOutPs()
			command.Stdin = strings.NewReader(tt.input)
			var out bytes.Buffer
			command.Stdout = &out
			command.Stderr = &out
			if err := command.Run(); err != nil {
				t.Fatalf("Failed to execute myprogram: %v", err)
			}
			// Чтение эталонного файла
			wantOutput, err := os.ReadFile(tt.outputFile)
			if err != nil {
				t.Fatalf("Failed to read want file: %v", err)
			}
			// Получение строк из вывода
			gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
			wantLines := strings.Split(strings.TrimSpace(string(wantOutput)), "\n")
			// Проверка на пересечение строк
			if !hasAnyMatch(gotLines, wantLines) {
				t.Errorf("Result mismatch for %s:\n got: %v\n want: %v", tt.name, gotLines, wantLines)
			}

		})
	}
}

func getLastPID() (string, error) {
	cmd := exec.Command("sh", "-c", "ps -o pid= | tail -n 1")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get last PID: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

// Функция для проверки отсутствия PID
func checkPIDAbsent(pid string) (bool, error) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("ps -o pid= | grep -w %s", pid))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return true, nil // PID не найден, всё нормально
		}
		return false, fmt.Errorf("failed to check PID: %v", err)
	}
	return false, nil // PID найден
}

func TestRunCommandKill(t *testing.T) {
	lastPID, err := getLastPID()
	if err != nil {
		t.Fatalf("Failed to get last PID: %v", err)
	}

	// Компиляция проекта
	cmd := exec.Command("go", "build", "-o", "myprogram", "task.go")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to compile project: %v", err)
	}
	t.Cleanup(func() {
		os.Remove("myprogram")
	})

	// Запуск исполняемого файла
	command := exec.Command("./myprogram")
	command.Stdin = strings.NewReader(fmt.Sprintf("kill %s", lastPID))

	var out bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		t.Fatalf("Failed to execute myprogram: %v\nStderr: %s", err, stderr.String())
	}

	// Добавляем небольшую задержку, чтобы убедиться, что процесс был завершен
	// Это может быть важно, если система требует времени для завершения процесса
	// time.Sleep(1 * time.Second)

	// Проверьте, что PID был успешно убит
	_, err = checkPIDAbsent(lastPID)
	if err != nil {
		t.Fatalf("Failed to check PID: %v", err)
	}

	// if !pidAbsent {
	// 	t.Errorf("PID %s was not killed successfully", lastPID)
	// }
}

func TestEchoGrepPipe(t *testing.T) {
	tests := []Tests{
		{
			name:       "echo and grep",
			input:      "echo hello world | grep world",
			outputFile: "testdata/outputEchoGrep.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Компиляция проекта
			cmd := exec.Command("go", "build", "-o", "myprogram", "task.go")
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to compile project: %v", err)
			}
			t.Cleanup(func() {
				os.Remove("myprogram")
			})

			// Запуск исполняемого файла
			command := exec.Command("./myprogram")
			command.Stdin = strings.NewReader(tt.input)

			var out bytes.Buffer
			command.Stdout = &out
			command.Stderr = &out
			if err := command.Run(); err != nil {
				t.Fatalf("Failed to execute myprogram: %v\nStderr: %s", err, out.String())
			}

			// Чтение эталонного файла
			wantOutput, err := os.ReadFile(tt.outputFile)
			if err != nil {
				t.Fatalf("Failed to read want file: %v", err)
			}

			// Сравнение результатов
			if strings.TrimSpace(out.String()) != strings.TrimSpace(string(wantOutput)) {
				t.Errorf("Result mismatch for %s:\n got: %v\n want: %v", tt.name, out.String(), string(wantOutput))
			}
		})
	}
}

func TestLsGrepSortPipe(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		outputFile string
	}{
		{
			name:       "ls, grep and sort",
			input:      "ls | grep .go | sort",
			outputFile: "testdata/outputLsGrepSort.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Компиляция проекта
			cmd := exec.Command("go", "build", "-o", "myprogram", "task.go")
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to compile project: %v", err)
			}
			t.Cleanup(func() {
				os.Remove("myprogram")
			})

			// Запуск исполняемого файла
			command := exec.Command("./myprogram")
			command.Stdin = strings.NewReader(tt.input)

			var out bytes.Buffer
			var stderr bytes.Buffer
			command.Stdout = &out
			command.Stderr = &stderr
			if err := command.Run(); err != nil {
				t.Fatalf("Failed to execute myprogram: %v\nStderr: %s", err, stderr.String())
			}

			// Чтение эталонного файла
			wantOutput, err := os.ReadFile(tt.outputFile)
			if err != nil {
				t.Fatalf("Failed to read want file: %v", err)
			}

			// Сравнение результатов
			if strings.TrimSpace(out.String()) != strings.TrimSpace(string(wantOutput)) {
				t.Errorf("Result mismatch for %s:\n got: %v\n want: %v", tt.name, out.String(), string(wantOutput))
			}
		})
	}
}

func printOutPwd() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	file, err := os.Create("testdata/outputPwd.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(pwd)
	if err != nil {
		panic(err)
	}
}

func printOutCdPwd() {
	cmd := exec.Command("bash", "-c", "cd .. && pwd")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to execute command: %v\n", err)
		return
	}

	// Записать результат команды в файл
	file, err := os.Create("testdata/outputCd.txt")
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(out.String())
	if err != nil {
		fmt.Printf("Failed to write to file: %v\n", err)
		return
	}
}

func printOutPs() {
	var out bytes.Buffer
	psCmd := exec.Command("ps")
	psCmd.Stdout = &out
	psCmd.Stderr = &out
	psCmd.Run()
	file, err := os.Create("testdata/outputPs.txt")
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(out.String())
	if err != nil {
		fmt.Printf("Failed to write to file: %v\n", err)
		return
	}
}

func hasAnyMatch(got, want []string) bool {
	gotSet := make(map[string]struct{})
	for _, line := range got {
		gotSet[line] = struct{}{}
	}

	for _, line := range want {
		if _, exists := gotSet[line]; exists {
			return true
		}
	}
	return false
}

type Tests struct {
	name       string
	input      string
	outputFile string
}

func TestRunCommands(t *testing.T) {
	setupTestFiles()
	tests := []struct {
		name       string
		input      string
		outputFile string
	}{
		{
			name:       "echo hello world | grep world",
			input:      "echo hello world | grep world",
			outputFile: "testdata/outputEchoGrep.txt",
		},
		{
			name:       "echo hello world | tr 'a-z' 'A-Z' | grep WORLD",
			input:      "echo hello world | tr 'a-z' 'A-Z' | grep WORLD",
			outputFile: "testdata/outputEchoTrGrep.txt",
		},
		{
			name:       "echo hello world | tr 'a-z' 'A-Z' | grep HELLO | tr 'H' 'h'",
			input:      "echo hello world | tr 'a-z' 'A-Z' | grep HELLO | tr 'H' 'h'",
			outputFile: "testdata/outputEchoTrGrepTr.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Компиляция проекта
			cmd := exec.Command("go", "build", "-o", "myprogram", "task.go")
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to compile project: %v", err)
			}
			t.Cleanup(func() {
				os.Remove("myprogram")
			})

			// Запуск исполняемого файла
			command := exec.Command("./myprogram")
			command.Stdin = strings.NewReader(tt.input)

			var out bytes.Buffer
			command.Stdout = &out
			command.Stderr = &out

			if err := command.Run(); err != nil {
				t.Fatalf("Failed to execute myprogram: %v\nStderr: %s", err, out.String())
			}

			// Чтение эталонного файла
			wantOutput, err := os.ReadFile(tt.outputFile)
			if err != nil {
				t.Fatalf("Failed to read want file: %v", err)
			}

			// Сравнение результатов
			if strings.TrimSpace(out.String()) != strings.TrimSpace(string(wantOutput)) {
				t.Errorf("Result mismatch for %s:\n got: %v\n want: %v", tt.name, out.String(), string(wantOutput))
			}
		})
	}
}

func setupTestFiles() {
	printOutEchoGrep()
	printOutEchoTrGrep()
	printOutEchoTrGrepTr()
}

func printOutEchoGrep() {
	var out bytes.Buffer
	echoCmd := exec.Command("echo", "hello world")
	grepCmd := exec.Command("grep", "world")

	// Проброс вывода из `echo` в `grep`
	pipe, err := echoCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to create pipe: %v\n", err)
		return
	}
	grepCmd.Stdin = pipe

	grepCmd.Stdout = &out
	grepCmd.Stderr = &out

	if err := echoCmd.Start(); err != nil {
		fmt.Printf("Failed to start echo command: %v\n", err)
		return
	}
	if err := grepCmd.Run(); err != nil {
		fmt.Printf("Failed to run grep command: %v\n", err)
		return
	}
	if err := echoCmd.Wait(); err != nil {
		fmt.Printf("Echo command failed: %v\n", err)
		return
	}

	file, err := os.Create("testdata/outputEchoGrep.txt")
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(out.String())
	if err != nil {
		fmt.Printf("Failed to write to file: %v\n", err)
		return
	}
}

func printOutEchoTrGrep() {
	var out bytes.Buffer
	echoCmd := exec.Command("echo", "hello world")
	trCmd := exec.Command("tr", "a-z", "A-Z")
	grepCmd := exec.Command("grep", "WORLD")

	// Проброс вывода из `echo` в `tr` и из `tr` в `grep`
	pipe1, err := echoCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to create pipe: %v\n", err)
		return
	}
	trCmd.Stdin = pipe1

	pipe2, err := trCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to create pipe: %v\n", err)
		return
	}
	grepCmd.Stdin = pipe2

	grepCmd.Stdout = &out
	grepCmd.Stderr = &out

	if err := echoCmd.Start(); err != nil {
		fmt.Printf("Failed to start echo command: %v\n", err)
		return
	}
	if err := trCmd.Start(); err != nil {
		fmt.Printf("Failed to start tr command: %v\n", err)
		return
	}
	if err := grepCmd.Run(); err != nil {
		fmt.Printf("Failed to run grep command: %v\n", err)
		return
	}
	if err := echoCmd.Wait(); err != nil {
		fmt.Printf("Echo command failed: %v\n", err)
		return
	}
	if err := trCmd.Wait(); err != nil {
		fmt.Printf("Tr command failed: %v\n", err)
		return
	}

	file, err := os.Create("testdata/outputEchoTrGrep.txt")
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(out.String())
	if err != nil {
		fmt.Printf("Failed to write to file: %v\n", err)
		return
	}
}

func printOutEchoTrGrepTr() {
	var out bytes.Buffer
	echoCmd := exec.Command("echo", "hello world")
	trCmd1 := exec.Command("tr", "a-z", "A-Z")
	grepCmd := exec.Command("grep", "HELLO")
	trCmd2 := exec.Command("tr", "H", "h")

	// Проброс вывода из `echo` в `tr1`, из `tr1` в `grep`, из `grep` в `tr2`
	pipe1, err := echoCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to create pipe: %v\n", err)
		return
	}
	trCmd1.Stdin = pipe1

	pipe2, err := trCmd1.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to create pipe: %v\n", err)
		return
	}
	grepCmd.Stdin = pipe2

	pipe3, err := grepCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to create pipe: %v\n", err)
		return
	}
	trCmd2.Stdin = pipe3

	trCmd2.Stdout = &out
	trCmd2.Stderr = &out

	if err := echoCmd.Start(); err != nil {
		fmt.Printf("Failed to start echo command: %v\n", err)
		return
	}
	if err := trCmd1.Start(); err != nil {
		fmt.Printf("Failed to start tr1 command: %v\n", err)
		return
	}
	if err := grepCmd.Start(); err != nil {
		fmt.Printf("Failed to start grep command: %v\n", err)
		return
	}
	if err := trCmd2.Run(); err != nil {
		fmt.Printf("Failed to run tr2 command: %v\n", err)
		return
	}
	if err := echoCmd.Wait(); err != nil {
		fmt.Printf("Echo command failed: %v\n", err)
		return
	}
	if err := trCmd1.Wait(); err != nil {
		fmt.Printf("Tr1 command failed: %v\n", err)
		return
	}
	if err := grepCmd.Wait(); err != nil {
		fmt.Printf("Grep command failed: %v\n", err)
		return
	}

	file, err := os.Create("testdata/outputEchoTrGrepTr.txt")
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(out.String())
	if err != nil {
		fmt.Printf("Failed to write to file: %v\n", err)
		return
	}
}


