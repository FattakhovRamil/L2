Что выведет программа? Объяснить вывод программы. Рассказать про внутреннее устройство слайсов и что происходит при передачи их в качестве аргументов функции.

```go
package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s)
}

func modifySlice(i []string) {
	i[0] = "3"
	i = append(i, "4")
	i[1] = "5"
	i = append(i, "6")
}
```

Ответ:
```
[3 2 3]
При i[0] = "3" мы изменили 0ой индекс слайса на значение 3. После этого i = append(i, "4") - командой создался новый слайс с len 4 и cap 4. Дальнешее изменение слайса i никак не затрагивает изначальный
слайс т.к они ссылаются на разные области памяти. 
Слайс представляет собой структуру, содержащую указатель на массив, длину и емкость.
При передаче слайса в функцию передается копия структуры с указателем на тот же массив.
Если операция append вызывает выделение нового массива (из-за превышения емкости), изменения затрагивают новый массив, а исходный массив остается неизменным.
```
