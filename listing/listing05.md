Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
} 

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
error
Т.к наш интерфейс error имеет тип значения error, а значение типа nil, то при проверке на nil, он выведет ошибку

```
