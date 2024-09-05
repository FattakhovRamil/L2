package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Flags struct {
	column               int  // -k
	numeric              bool // -n
	reverse              bool // -r
	unique               bool // -u
	month                bool // -M
	ignoreTrailingSpaces bool // -b
	checkSorted          bool // -c
	humanNumeric         bool // -h
}

func parseFlags() (Flags, string, string) {
	// Определение флагов
	flagSet := flag.NewFlagSet("sort", flag.ExitOnError)

	column := flagSet.Int("k", 1, "Указание колонки для сортировки (по умолчанию 1)")
	numeric := flagSet.Bool("n", false, "Сортировать по числовому значению")
	reverse := flagSet.Bool("r", false, "Сортировать в обратном порядке")
	unique := flagSet.Bool("u", false, "Не выводить повторяющиеся строки")
	month := flagSet.Bool("M", false, "Сортировать по названию месяца")
	ignoreTrailingSpaces := flagSet.Bool("b", false, "Игнорировать хвостовые пробелы")
	checkSorted := flagSet.Bool("c", false, "Проверять отсортированы ли данные")
	humanNumeric := flagSet.Bool("h", false, "Сортировать по числовому значению с учётом суффиксов")

	// Парсинг флагов
	flagSet.Parse(os.Args[1:])

	// Получение оставшихся аргументов (файлов)
	args := flagSet.Args()
	if len(args) < 1 {
		fmt.Println("Необходимо указать входной файл")
		os.Exit(1)
	}
	inputFile := args[0]
	var outputFile string
	if len(args) > 1 {
		outputFile = args[1]
	}

	return Flags{
		column:               *column,
		numeric:              *numeric,
		reverse:              *reverse,
		unique:               *unique,
		month:                *month,
		ignoreTrailingSpaces: *ignoreTrailingSpaces,
		checkSorted:          *checkSorted,
		humanNumeric:         *humanNumeric,
	}, inputFile, outputFile
}

func main() {
	flags, inputFile, outputFile := parseFlags()

	// Чтение входного файла
	lines, err := readFile(inputFile)
	if err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		os.Exit(1)
	}

	// Функция для получения нужной колонки
	getColumn := func(line string) string {
		if flags.ignoreTrailingSpaces {
			line = strings.TrimRight(line, " \t")
		}
		words := strings.Fields(line)
		if flags.column-1 < len(words) {
			return words[flags.column-1]
		}
		return ""
	}

	// Проверка отсортированности данных, если указан флаг -c
	if flags.checkSorted {
		if !isSorted(lines, getColumn, flags) {
			fmt.Println("Файл не отсортирован")
			os.Exit(1)
		}
		fmt.Println("Файл отсортирован")
		os.Exit(0)
	}

	// Применение флагов сортировки
	sort.SliceStable(lines, func(i, j int) bool {
		key1 := getColumn(lines[i])
		key2 := getColumn(lines[j])

		if flags.month {
			return compareMonths(key1, key2, flags.reverse)
		}

		if flags.humanNumeric {
			num1, err1 := parseHumanNumeric(key1)
			num2, err2 := parseHumanNumeric(key2)
			if err1 == nil && err2 == nil {
				if flags.reverse {
					return num1 > num2
				}
				return num1 < num2
			}
		}

		if flags.numeric {
			num1, err1 := strconv.Atoi(key1)
			num2, err2 := strconv.Atoi(key2)
			if err1 == nil && err2 == nil {
				if flags.reverse {
					return num1 > num2
				}
				return num1 < num2
			}
		}

		if flags.reverse {
			return key1 > key2
		}
		return key1 < key2
	})

	// Удаление дубликатов, если указан флаг -u
	if flags.unique {
		uniqueLines := []string{}
		seen := map[string]bool{}
		for _, line := range lines {
			normalizedLine := line
			if flags.ignoreTrailingSpaces {
				normalizedLine = strings.TrimRight(line, " \t")
			}
			if !seen[normalizedLine] {
				uniqueLines = append(uniqueLines, line)
				seen[normalizedLine] = true
			}
		}
		lines = uniqueLines
	}

	// Пример вывода результата
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

func readFile(inputFile string) ([]string, error) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		return []string{}, errors.New("ошибка при чтении файла")
	}

	// Разделение строк
	lines := strings.Split(string(input), "\n")
	return lines, nil
}

func compareMonths(month1, month2 string, reverse bool) bool {
	months := map[string]int{
		"январь": 1, "февраль": 2, "март": 3, "апрель": 4,
		"май": 5, "июнь": 6, "июль": 7, "август": 8,
		"сентябрь": 9, "октябрь": 10, "ноябрь": 11, "декабрь": 12,
	}
	m1, ok1 := months[strings.ToLower(month1)]
	m2, ok2 := months[strings.ToLower(month2)]
	if !ok1 || !ok2 {
		if reverse {
			return month1 > month2
		}
		return month1 < month2
	}
	if reverse {
		return m1 > m2
	}
	return m1 < m2
}

func parseHumanNumeric(value string) (int64, error) {
	suffixes := map[byte]int64{
		'K': 1e3,
		'M': 1e6,
		'G': 1e9,
		'T': 1e12,
	}
	if len(value) == 0 {
		return 0, fmt.Errorf("invalid value")
	}
	lastChar := value[len(value)-1]
	if multiplier, ok := suffixes[lastChar]; ok {
		number, err := strconv.ParseInt(value[:len(value)-1], 10, 64)
		if err != nil {
			return 0, err
		}
		return number * multiplier, nil
	}
	return strconv.ParseInt(value, 10, 64)
}

func isSorted(lines []string, getColumn func(string) string, flags Flags) bool {
	for i := 1; i < len(lines); i++ {
		key1 := getColumn(lines[i-1])
		key2 := getColumn(lines[i])
		if flags.ignoreTrailingSpaces {
			key1 = strings.TrimSpace(key1)
			key2 = strings.TrimSpace(key2)
		}
		if flags.month {
			if compareMonths(key1, key2, flags.reverse) {
				continue
			}
			return false
		}
		if flags.humanNumeric {
			num1, err1 := parseHumanNumeric(key1)
			num2, err2 := parseHumanNumeric(key2)
			if err1 == nil && err2 == nil {
				if flags.reverse && num1 < num2 {
					return false
				} else if !flags.reverse && num1 > num2 {
					return false
				}
				continue
			}
		}
		if flags.numeric {
			num1, err1 := strconv.Atoi(key1)
			num2, err2 := strconv.Atoi(key2)
			if err1 == nil && err2 == nil {
				if flags.reverse && num1 < num2 {
					return false
				} else if !flags.reverse && num1 > num2 {
					return false
				}
				continue
			}
		}
		if flags.reverse && key1 < key2 {
			return false
		} else if !flags.reverse && key1 > key2 {
			return false
		}
	}
	return true
}
