package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать -N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type flags struct {
	after      uint // -A
	before     uint // -B
	context    uint // -C
	count      bool // -c
	ignoreCase bool // -i
	invert     bool // -v
	fixed      bool // -F
	lineNum    bool // -n
}

func parseFlags() (flags, string, string, string) {
	// Определение флагов
	flagSet := flag.NewFlagSet("grep", flag.ExitOnError)

	after := flagSet.Uint("A", 0, "\"after\" печатать +N строк после совпадения")
	before := flagSet.Uint("B", 0, "\"before\" печатать -N строк до совпадения")
	context := flagSet.Uint("C", 0, "\"context\" (A+B) печатать ±N строк вокруг совпадения")
	count := flagSet.Bool("c", false, "\"count\" (количество строк)")
	ignoreCase := flagSet.Bool("i", false, "\"ignore-case\" (игнорировать регистр)")
	invert := flagSet.Bool("v", false, "\"invert\" (вместо совпадения, исключать)")
	fixed := flagSet.Bool("F", false, "\"fixed\", точное совпадение со строкой, не паттерн")
	lineNum := flagSet.Bool("n", false, "\"line num\", печатать номер строки")

	// Парсинг флагов
	flagSet.Parse(os.Args[1:])

	// Получение оставшихся аргументов (файлов)
	args := flagSet.Args()
	if len(args) < 1 {
		fmt.Println("Необходимо указать входной файл")
		os.Exit(1)
	}
	patternWord := args[0]
	inputFile := args[1]
	var outputFile string
	if len(args) > 2 {
		outputFile = args[2]
	}

	return flags{
		after:      *after,
		before:     *before,
		context:    *context,
		count:      *count,
		ignoreCase: *ignoreCase,
		invert:     *invert,
		fixed:      *fixed,
		lineNum:    *lineNum,
	}, inputFile, outputFile, patternWord
}

func main() {
	flags, inputFile, outputFile, patternWord := parseFlags()

	data, err := readFlow(inputFile) // Получаем массив строк
	if err != nil {
		fmt.Println("Ошибка:", err)
		os.Exit(1)
	}

	result, err := flagFilter(data, &patternWord, &flags)
	if err != nil {
		fmt.Println("Ошибка:", err)
		os.Exit(1)
	}
	writeFile(outputFile, result)

}

func writeFile(outputFile string, lines []string) {
	output := strings.Join(lines, "\n")
	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(output), 0644)
		if err != nil {
			fmt.Printf("Ошибка при записи файла: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(output)
	}
}

func readFlow(inputFile string) ([]string, error) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		//return []string{}, errors.New("ошибка при чтении файла")
		lines := strings.Split(string(inputFile), "\n") // случай, если у нас вместо названия файла указан сразу текст
		return lines, nil
	}

	// Разделение строк
	lines := strings.Split(string(input), "\n")
	return lines, nil
}

func flagFilter(data []string, patternWord *string, flags *flags) ([]string, error) {
	counter := 0
	result := make([]string, 0, len(data)-1)
	if flags.ignoreCase {
		*patternWord = "(?i)" + *patternWord
	}
	lineCounter := make([]int, 0)
	for i, v := range data {
		line := v
		if flags.ignoreCase {
			line = strings.ToLower(v)
			*patternWord = strings.ToLower(*patternWord)
		}

		matched := false
		if flags.fixed {
			matched = strings.Contains(line, *patternWord)
		} else {
			re, err := regexp.Compile(*patternWord)
			if err != nil {
				fmt.Println("Ошибка компиляции регулярного выражения:", err)
				continue
			}
			matched = re.MatchString(line)
		}

		if matched && !flags.invert {
			if flags.lineNum {
				result = append(result, fmt.Sprintf("%d:%s", i+1, v))
			} else {
				result = append(result, v)
			}
			lineCounter = append(lineCounter, i)
			counter++
		} else if !matched && flags.invert {
			if flags.lineNum {
				result = append(result, fmt.Sprintf("%d:%s", i+1, v))
			} else {
				result = append(result, v)
			}
			lineCounter = append(lineCounter, i)
			counter++
		}
	}

	if flags.after != 0 || flags.before != 0 || flags.context != 0 {
		result = extraLines(flags, data, lineCounter)
		return result, nil
	}

	// Найти все совпадения
	if flags.count {
		return []string{fmt.Sprintf("%d", counter)}, nil
	}
	return result, nil
}

func extraLines(flags *flags, data []string, lineCounter []int) []string {
	before := flags.before
	after := flags.after

	if flags.context != 0 {
		before = flags.context
		after = flags.context
	}

	uniqueLines := make(map[int]struct{})

	for _, lineIndex := range lineCounter {
		// Добавляем строки до совпадения
		for i := lineIndex - int(before); i < lineIndex; i++ {
			if i >= 0 { // Проверка на неотрицательный индекс
				uniqueLines[i] = struct{}{}
			}
		}
		// Добавляем строки после совпадения
		for i := lineIndex; i <= lineIndex+int(after); i++ {
			if i < len(data) { // Проверка на индекс в пределах массива
				uniqueLines[i] = struct{}{}
			}
		}
	}

	// Сортировка уникальных строк по индексам
	sortedIndices := make([]int, 0, len(uniqueLines))
	for idx := range uniqueLines {
		sortedIndices = append(sortedIndices, idx)
	}
	sort.Ints(sortedIndices)

	// Создание результата на основе отсортированных индексов
	result := make([]string, 0, len(sortedIndices))
	for _, idx := range sortedIndices {

		if flags.lineNum {

			result = append(result, fmt.Sprintf("%d:%s", idx+1, data[idx]))
		} else {
			result = append(result, data[idx])
		}

	}

	return result
}
