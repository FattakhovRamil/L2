package main

import "testing"

func TestCheckString1(t *testing.T) {
	str := "a4bc2d5e"

	err := checkString(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}
}

func TestCheckString2(t *testing.T) {
	str := "abcd"

	err := checkString(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}
}

func TestCheckString3(t *testing.T) {
	str := "45"

	err := checkString(str)

	if err == nil {
		t.Fatalf("ожидалось ошибка, но получено успешное выполнение")
	}
}

func TestCheckString4(t *testing.T) {
	str := "4"

	err := checkString(str)

	if err == nil {
		t.Fatalf("ожидалось ошибка, но получено успешное выполнение")
	}
}

func TestCheckString5(t *testing.T) {
	str := "4a"

	err := checkString(str)

	if err == nil {
		t.Fatalf("ожидалось ошибка, но получено успешное выполнение")
	}
}

func TestCheckString6(t *testing.T) {
	str := "a4"

	err := checkString(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}
}

func TestCheckString7(t *testing.T) {
	str := ""

	err := checkString(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}
}

func TestCheckString8(t *testing.T) {
	str := "qwe\\4\\5"

	err := checkString(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}
}

func TestCheckString9(t *testing.T) {
	str := "qwe\\45"

	err := checkString(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}
}

func TestCheckString10(t *testing.T) {
	str := "qwe\\5"

	err := checkString(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}
}

func TestUnpacking1(t *testing.T) {
	str := "qwe\\5"
	expected := "qwe5"
	result, err := unpacking(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}

	if expected != result {
		t.Fatalf("Строки не равны. Ожидалось %s, а получили %s", expected, result)
	}
}

func TestUnpacking2(t *testing.T) {
	str := "a4bc2d5e"
	expected := "aaaabccddddde"
	result, err := unpacking(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}

	if expected != result {
		t.Fatalf("Строки не равны. Ожидалось %s, а получили %s", expected, result)
	}
}

func TestUnpacking3(t *testing.T) {
	str := "abcd"
	expected := "abcd"
	result, err := unpacking(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}

	if expected != result {
		t.Fatalf("Строки не равны. Ожидалось %s, а получили %s", expected, result)
	}
}

func TestUnpacking4(t *testing.T) {
	str := "45"

	_, err := unpacking(str)

	if err == nil {
		t.Fatalf("ожидалось ошибка, но получено успешное выполнение")
	}

}

func TestUnpacking5(t *testing.T) {
	str := ""
	_, err := unpacking(str)
	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}
}

func TestUnpacking6(t *testing.T) {
	str := "qwe\\4\\5"
	expected := "qwe45"
	result, err := unpacking(str)
	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}

	if expected != result {
		t.Fatalf("Строки не равны. Ожидалось %s, а получили %s", expected, result)
	}
}

func TestUnpacking7(t *testing.T) {
	str := "qwe\\45"
	expected := "qwe44444"
	result, err := unpacking(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}

	if expected != result {
		t.Fatalf("Строки не равны. Ожидалось %s, а получили %s", expected, result)
	}

}

func TestUnpacking8(t *testing.T) {
	str := "qwe\\5"
	expected := "qwe5"
	result, err := unpacking(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}

	if expected != result {
		t.Fatalf("Строки не равны. Ожидалось %s, а получили %s", expected, result)
	}

}

func TestUnpacking9(t *testing.T) {
	str := "a4"
	expected := "aaaa"
	result, err := unpacking(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}

	if expected != result {
		t.Fatalf("Строки не равны. Ожидалось %s, а получили %s", expected, result)
	}

}

func TestUnpacking10(t *testing.T) {
	str := "\\\\"
	expected := "\\"
	result, err := unpacking(str)

	if err != nil {
		t.Fatalf("ожидалось успешное выполнение, но получена ошибка: %v", err)
	}

	if expected != result {
		t.Fatalf("Строки не равны. Ожидалось %s, а получили %s", expected, result)
	}

}

