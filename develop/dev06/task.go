package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type flags struct {
	fields    string // -f
	delimiter string // -d
	separated bool   // -s
}

func parseFlags() flags {
	// Определение флагов
	fields := flag.String("f", "", "выбрать поля (колонки)")
	delimiter := flag.String("d", "\t", "использовать другой разделитель")
	separated := flag.Bool("s", false, "только строки с разделителем")

	// Парсинг флагов
	flag.Parse()

	if *fields == "" {
		fmt.Println("Необходимо указать поля с помощью флага -f")
		os.Exit(1)
	}

	return flags{
		fields:    *fields,
		delimiter: *delimiter,
		separated: *separated,
	}
}

func main() {
	flags := parseFlags()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		if flags.separated && !strings.Contains(line, flags.delimiter) {
			continue // пропустить строки, не содержащие разделитель
		}

		columns := strings.Split(line, flags.delimiter)
		fieldIndexes := parseFields(flags.fields)
		var result []string

		for _, idx := range fieldIndexes {
			if idx-1 < len(columns) {
				result = append(result, columns[idx-1])
			}
		}
		fmt.Println(strings.Join(result, flags.delimiter))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при чтении ввода: %v\n", err)
		os.Exit(1)
	}
}

func parseFields(fields string) []int {
	indexes := make([]int, 0)
	fieldParts := strings.Split(fields, ",")

	for _, part := range fieldParts {
		var idx int
		fmt.Sscanf(part, "%d", &idx)
		indexes = append(indexes, idx)
	}
	return indexes
}
