package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {

	slice := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"}

	result, err := anagram(slice)
	if err != nil {
		fmt.Println("Errpor: ", err)
	}
	fmt.Println(result, err)
	anagram(slice)
}

func anagram(words []string) (map[string][]string, error) {
	if len(words) < 1 {
		return nil, errors.New("пустая строка")
	}

	// Привести все слова к нижнему регистру
	loweredWords := toLowerCaseArray(words)

	// Карта для хранения анаграмм
	anagramMap := make(map[string][]string)

	// Заполняем карту отсортированными строками как ключами и списком слов как значениями
	for _, word := range loweredWords {
		sortedWord := sortString(word)
		anagramMap[sortedWord] = append(anagramMap[sortedWord], word)
	}

	result := make(map[string][]string)
	for _, anagrams := range anagramMap {
		if len(anagrams) > 1 {
			sort.Strings(anagrams) // Сортируем каждое множество анаграмм
			result[anagrams[0]] = anagrams[1:]
		}
	}

	return result, nil
}

func sortString(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	return strings.Join(s, "")
}

func toLowerCaseArray(words []string) []string {
	loweredWords := make([]string, len(words))
	for i, word := range words {
		loweredWords[i] = strings.ToLower(word)
	}
	return loweredWords
}
