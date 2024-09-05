package main

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"errors"
	//"fmt"
	//"os"
	"strconv"
)

func checkString(str string) error {
	if len(str) == 0 {
		return nil
	}

	if _, err := strconv.Atoi(string(str[0])); err == nil {
		return errors.New("ошибка 1 символ это число")
	}

	for i := 0; i < len(str)-1; i++ {
		if str[i] == '\\' {
			i += 1
			continue
		}
		if _, err := strconv.Atoi(string(str[i])); err == nil {
			//fmt.Printf("%q looks like a number.\n", v) // это число
			if len(str) > i+1 { // проверяем, не выходит за пределы строки
				if _, err := strconv.Atoi(string(str[i+1])); err == nil { // проверка число ли следующий элемент
					//fmt.Printf("%q looks like a number.\n", v) // это число -- ошибка в стандарном исполнении
					return errors.New("ошибка 2 числа рядом")
				}
			}
		}
	}
	return nil
}

func unpacking(str string) (result string, err error) {

	if len(str) == 0 {
		return "", nil
	}

	err = checkString(str)
	if err != nil {
		return "", errors.New("не подходящая строка")

	}
	slice := make([]rune, 0, len(str))

	for i, _ := range str {
		if number, err := strconv.Atoi(string(str[i])); err == nil {
			if str[i-1] == '\\' {
				if str[i-2] == '\\' {
					for j := 0; j < number; j++ {
						slice = append(slice, rune(str[i-1]))
					}
					continue
				} else {
					slice = append(slice, rune(str[i]))
					continue
				}
			} else {
				for j := 1; j < number; j++ {
					slice = append(slice, rune(str[i-1]))
				}
				continue
			}
		}
		if str[i] == '\\' {
			if len(str) > i+1 && str[i+1] == '\\' {
				slice = append(slice, rune(str[i]))
			}
			i++
			continue
		}
		slice = append(slice, rune(str[i]))
	}

	return string(slice[:]), err
}

// func main() {
// 	t := "\\\\"
// 	//t := ""
// 	result, err := unpacking(t)

// 	if err != nil {
// 		fmt.Println("Ошибка")
// 		os.Exit(1)
// 	}
// 	fmt.Println(result)
// 	fmt.Println("Проверка прошла успешно")
// 	os.Exit(0)
// }
