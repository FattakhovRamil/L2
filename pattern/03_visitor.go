package pattern

/*
	Реализовать паттерн «посетитель».
	Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
	1) Применимость:
		Паттерн Посетитель применяется, когда нужно выполнить операцию над элементами сложной структуры объектов (например, дерева) и не хочется засорять классы этой структуры новыми операциями. Он решает несколько проблем:
		* Добавление новых операций: Позволяет легко добавлять новые операции, не изменяя сами классы объектов, над которыми эти операции выполняются.
		* Изоляция поведения: Выносит родственные операции в отдельные классы-посетители, уменьшая загрязнение исходных классов.
		* Обход структур: Упрощает выполнение операций над объектами различных классов в составе одной структуры.

	2) Плюсы и минусы:
		Плюсы:
		+ Упрощенное добавление операций: Легко добавлять новые операции, не изменяя классы объектов.
		+ Объединение операций: Родственные операции объединяются в одном классе-посетителе.
		+ Сохранение состояния: Посетитель может аккумулировать состояние при обходе структуры компонентов.
		Минусы:
		- Изменения в иерархии компонентов: Если иерархия компонентов часто меняется, необходимо обновлять всех посетителей, что усложняет поддержку.
		- Нарушение инкапсуляции: Посетителю может потребоваться доступ к приватным полям компонентов, что может нарушить их инкапсуляцию.
		- Усложнение структуры: Неоправданное использование паттерна может привести к усложнению структуры кода.

	3) Реальные примеры использования на практике:
		* Графические редакторы: Например, в графическом редакторе нужно экспортировать разные фигуры (круги, прямоугольники, сложные фигуры) в XML. Используя паттерн Посетитель, можно создать класс-посетитель для экспорта в XML, который будет знать, как экспортировать каждую фигуру.
		* Компиляторы: В компиляторах паттерн Посетитель может применяться для выполнения различных операций над синтаксическим деревом, таких как анализ, оптимизация и генерация кода.
		* Базы данных: В системах, работающих с базами данных, можно использовать паттерн для выполнения различных операций над таблицами и записями, таких как сбор статистики, валидация данных или преобразование данных.

	4) Как связан с другими паттернами
		* Команда (Command): Посетитель можно рассматривать как расширенный аналог Команды, способный работать сразу с несколькими видами получателей.
		* Компоновщик (Composite): Посетитель может обходить и выполнять операции над всей структурой, созданной с использованием паттерна Компоновщик.
		* Итератор (Iterator): Можно использовать совместно с Итератором. Итератор отвечает за обход структуры данных, а Посетитель — за выполнение действий над каждым компонентом.
		* Стратегия (Strategy): Позволяет инкапсулировать различные алгоритмы. В случае Посетителя, разные посетители могут выполнять разные операции над компонентами, что схоже с паттерном Стратегия.
*/

import "fmt"

// ElectricDevice интерфейс для всех электрических устройств
type ElectricDevice interface {
	acceptVisitor(v ElectricityMeter) // Метод acceptVisitor принимает объект типа ElectricityMeter
}

// Структура Phone, представляющая телефон
type Phone struct {
}

// Реализация метода acceptVisitor для Phone
// Визитор вызывает метод measureElectricityPhone
func (p *Phone) acceptVisitor(v ElectricityMeter) {
	v.measureElectricityPhone(p) // Вызывается метод measureElectricityPhone на визиторе
}

// Метод featurePhone возвращает потребление электроэнергии телефоном
func (p *Phone) featurePhone() int {
	return 100
}

// Структура Laptop, представляющая ноутбук
type Laptop struct {
}

// Реализация метода acceptVisitor для Laptop
// Визитор вызывает метод measureElectricityLaptop
func (l *Laptop) acceptVisitor(v ElectricityMeter) {
	v.measureElectricityLaptop(l) // Вызывается метод measureElectricityLaptop на визиторе
}

// Метод featureLaptop возвращает потребление электроэнергии ноутбуком
func (l *Laptop) featureLaptop() int {
	return 300
}

// Структура TV, представляющая телевизор
type TV struct {
}

// Реализация метода acceptVisitor для TV
// Визитор вызывает метод measureElectricityTV
func (tv *TV) acceptVisitor(v ElectricityMeter) {
	v.measureElectricityTV(tv) // Вызывается метод measureElectricityTV на визиторе
}

// Метод featureTV возвращает потребление электроэнергии телевизором
func (tv *TV) featureTV() int {
	return 250
}

// Интерфейс ElectricityMeter для визитора, который считает потребление электроэнергии
type ElectricityMeter interface {
	// Методы для измерения потребления электроэнергии разными устройствами
	measureElectricityPhone(p *Phone)
	measureElectricityLaptop(l *Laptop)
	measureElectricityTV(tv *TV)
}

// Конкретная реализация ElectricityMeter
type ConcreteElecticityMeter struct {
}

// Реализация метода measureElectricityPhone
// Выводит потребление электроэнергии телефоном
func (m *ConcreteElecticityMeter) measureElectricityPhone(p *Phone) {
	fmt.Println(p.featurePhone())
}

// Реализация метода measureElectricityLaptop
// Выводит потребление электроэнергии ноутбуком
func (m *ConcreteElecticityMeter) measureElectricityLaptop(l *Laptop) {
	fmt.Println(l.featureLaptop())
}

// Реализация метода measureElectricityTV
// Выводит потребление электроэнергии телевизором
func (m *ConcreteElecticityMeter) measureElectricityTV(tv *TV) {
	fmt.Println(tv.featureTV())
}

// func main() {
// 	// Создаем устройства
// 	phone := &Phone{}
// 	laptop := &Laptop{}
// 	tv := &TV{}

// 	// Создаем визитора
// 	visitor := &ConcreteElectricityMeter{}

// 	// Визитор посещает каждое устройство
// 	phone.acceptVisitor(visitor)
// 	laptop.acceptVisitor(visitor)
// 	tv.acceptVisitor(visitor)
// }
