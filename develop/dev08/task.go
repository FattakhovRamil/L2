package main

import (
	//"bufio"

	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах
*/

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// Пример ввода для тестирования
		//input := "echo hello world | grep world"
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if strings.TrimSpace(input) == "\\quit" {
			fmt.Println("Exiting shell...")
			os.Exit(0)
		}

		args := strings.TrimSpace(input)
		if len(args) == 0 {
			continue
		}

		result, err := processPipeline(args)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println(result)
	}
}

// processPipeline обрабатывает команду с пайпами
func processPipeline(input string) (string, error) {
	commands := strings.Split(input, "|")
	var prevOutput io.Reader = nil

	for i, command := range commands {
		args := strings.Fields(strings.TrimSpace(command))
		if len(args) == 0 {
			continue
		}

		cmd := exec.Command(args[0], args[1:]...)
		if i == 0 {
			// Для первой команды в конвейере
			if prevOutput != nil {
				cmd.Stdin = prevOutput
			}
		} else {
			// Для остальных команд в конвейере
			cmd.Stdin = prevOutput
		}

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to execute command %v: %v", args, err)
		}

		// Для последней команды в конвейере
		if i == len(commands)-1 {
			return out.String(), nil
		}

		// Для промежуточных команд в конвейере
		prevOutput = &out
	}

	return "", nil
}

// func runCommand(args []string, beforeStr string) (result string) {
// 	switch args[0] {
// 	case "cd":
// 		if len(args) == 1 {
// 			pwd, err := os.Getwd()
// 			if err != nil {
// 				fmt.Println(err)
// 				os.Exit(1)
// 			}
// 			result = pwd
// 		} else if beforeStr != "" {
// 			err := os.Chdir(beforeStr)
// 			if err != nil {
// 				fmt.Println("cd:", err)
// 				os.Exit(1)
// 			}
// 		} else {
// 			err := os.Chdir(args[1])
// 			if err != nil {
// 				fmt.Println("cd:", err)
// 				os.Exit(1)
// 			}
// 		}
// 	case "pwd":
// 		pwd, err := os.Getwd()
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 		result = pwd
// 	case "echo":
// 		result = strings.Join(args[1:], " ")
// 	case "kill":
// 		if len(args) < 2 && beforeStr == "" {
// 			fmt.Println("No arguments")
// 		} else {
// 			var programID int
// 			if beforeStr != "" {
// 				fmt.Sscanf(beforeStr, "%d", &programID)
// 			} else {
// 				fmt.Sscanf(args[1], "%d", &programID)
// 			}
// 			err := exec.Command("kill", fmt.Sprintf("%d", programID)).Run()
// 			if err != nil {
// 				fmt.Printf("Failed to kill process: %d\n", programID)
// 				os.Exit(1)
// 			}
// 		}
// 		return ""
// 	case "ps":
// 		var out bytes.Buffer
// 		psCmd := exec.Command("ps")
// 		psCmd.Stdout = &out
// 		psCmd.Stderr = &out
// 		psCmd.Run()
// 		result = out.String()
// 		lines := strings.Split(result, "\n")
// 		result = strings.Join(lines[:len(lines)-2], "\n")
// 	case "\\quit":
// 		fmt.Println("Exit my shell")
// 		os.Exit(0)
// 	default:
// 		return otherCommands(args, beforeStr)
// 	}
// 	return
// }

// func otherCommands(args []string, resuilIn string) (result string) {
// 	if len(args) == 0 {
// 		return ""
// 	}
// 	var out bytes.Buffer
// 	if resuilIn != "" {
// 		args = append(args, resuilIn)
// 	}
// 	cmd := exec.Command(args[0], args[1:]...)
// 	cmd.Stdout = &out
// 	cmd.Stderr = &out
// 	err := cmd.Run()
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	result = out.String()
// 	return result
// }
