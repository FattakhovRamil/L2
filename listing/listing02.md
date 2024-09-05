Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)


func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}


func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
```

Ответ:
```
2
1
test() возвращает 2, потому что defer-функция увеличивает x до возвращения значения.
anotherTest() возвращает 1, потому что значение уже было возвращено до выполнения defer-функции.
Defer-функции выполняются в обратном порядке их объявления. То есть, последний отложенный вызов будет выполнен первым.

```
