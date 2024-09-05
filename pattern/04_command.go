package pattern

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

/*
	1) Применимость:
	* Параметризация объектов: Команда позволяет превращать операции в объекты, которые можно передавать, хранить и взаимозаменять внутри других объектов. Это особенно полезно при разработке библиотек, где требуется гибкая настройка поведения без изменения исходного кода.
	* Очередь, выполнение по расписанию, передача по сети: Команды можно сериализовать для сохранения и последующего выполнения. Это подходит для задач, требующих выполнения операций в определённое время, их передачи по сети или логирования.
	* Отмена операций: Паттерн Команда позволяет легко реализовать отмену операций благодаря хранению истории выполненных команд.

	2) Плюсы и минусы:
		Плюсы:
		+ Уменьшение зависимости: Устраняет прямую зависимость между объектами, инициирующими операции, и объектами, выполняющими их.
		+ Отмена и повтор операций: Упрощает реализацию функциональности отмены и повтора действий.
		+ Отложенный запуск: Команды можно сохранять для последующего выполнения.
		+ Сложные команды: Позволяет объединять простые команды в более сложные структуры.
		+ Принцип открытости/закрытости: Легко добавлять новые команды без изменения существующего кода.
		Минусы:
		Усложнение кода: Паттерн Команда требует создания множества дополнительных классов, что может усложнить структуру программы.

	3) Реальные примеры использования на практике:
		* В текстовых редакторах паттерн Команда используется для реализации операций с текстом, таких как копирование, вырезание, вставка, и отмена этих операций.
		* Паттерн применяется для связывания действий интерфейса с бизнес-логикой. Пример: В графических интерфейсах кнопки и горячие клавиши могут быть привязаны к объектам команд, которые взаимодействуют с бизнес-логикой.
		* В системах, где требуется выполнение задач по расписанию или в определённом порядке, команды сериализуются и выполняются в нужное время.

	4) Как связан с другими паттернами
		* Цепочка обязанностей: Передаёт запросы через цепочку обработчиков, в то время как Команда устанавливает одностороннюю связь между отправителями и получателями.
		* Посредник: Устраняет прямые связи между отправителями и получателями, заставляя их взаимодействовать через посредника.
		* Наблюдатель: Позволяет отправителю уведомлять множество получателей о событии, с возможностью динамической подписки и отписки.
		* Снимок: Комбинируется с паттерном Команда для реализации отмены операций, где команды сохраняют состояние объекта до выполнения действия.
		* Стратегия: Хотя и схож с Командой, Стратегия фокусируется на взаимозаменяемых способах выполнения одного действия, а не на параметризации различных действий.
		* Прототип: Используется для копирования команд перед добавлением в историю выполненных операций.
		* Посетитель: Можно рассматривать как расширенную версию Команды, работающую с несколькими видами получателей.
*/

// import "fmt"

// TextCommand интерфейс для всех команд
type TextCommand interface {
	execute()
}

// BaseTextCommand базовая структура для команд, содержит текст и ссылку на получателя (TextEditor)
type BaseTextCommand struct {
	text     string
	receiver *TextEditor
}

// PasteTextCommand конкретная команда для вставки текста
type PastTextCommand struct {
	BaseTextCommand
}

// execute метод для вставки текста
func (c *PastTextCommand) execute() {
	c.BaseTextCommand.receiver.textField = c.BaseTextCommand.receiver.textField + c.BaseTextCommand.text
}

// CopyTextCommand конкретная команда для копирования текста
type CopyTextCommand struct {
	BaseTextCommand
}

// execute метод для копирования текста
func (c *CopyTextCommand) execute() {
	c.BaseTextCommand.receiver.textField += " " + c.BaseTextCommand.text
}

// DeleteTextCommand конкретная команда для удаления текста
type DeleteTextCommand struct {
	BaseTextCommand
}

// execute метод для удаления текста
func (c *DeleteTextCommand) execute() {
	c.BaseTextCommand.receiver.textField = ""
}

// TextEditor получатель команд, содержит текстовое поле
type TextEditor struct {
	textField string
}

// getTextField метод для получения значения текстового поля
func (te *TextEditor) getTextField() string {
	return te.textField
}

// TextInvoker отправитель команд, вызывает их выполнение
type TextInvoker struct {
	command TextCommand
}

// setCommand метод для установки команды в отправителе
func (ti *TextInvoker) setCommand(cm TextCommand) {
	ti.command = cm
}

// executeCommand метод для выполнения установленной команды
func (ti *TextInvoker) executeCommand() {
	ti.command.execute()
}

// func main() {

// 	// создаем получателя, команду (связать с получателем), отправителя (связать с командой)
// 	textEditor := &TextEditor{}
// 	textExample := "past text example"
// 	pastTextCommand := &PastTextCommand{BaseTextCommand{
// 		text:     textExample,
// 		receiver: textEditor,
// 	}}

// 	textInvoker := &TextInvoker{command: pastTextCommand}
// 	textInvoker.executeCommand()

// 	fmt.Println(textEditor.getTextField())

// 	copyTextCommand := &CopyTextCommand{BaseTextCommand{
// 		text:     "copy text example",
// 		receiver: textEditor,
// 	}}
// 	textInvoker.setCommand(copyTextCommand)
// 	textInvoker.executeCommand()
// 	fmt.Println(textEditor.getTextField())

// 	deleteTextCommand := &DeleteTextCommand{BaseTextCommand{

// 		receiver: textEditor,
// 	}}
// 	textInvoker.setCommand(deleteTextCommand)
// 	textInvoker.executeCommand()
// 	fmt.Println(textEditor.getTextField())
// }
