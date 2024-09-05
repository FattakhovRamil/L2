package pattern

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

/*
	1) Применимость: Паттерн Состояние применяется в тех случаях, когда поведение объекта зависит от его внутреннего состояния и это поведение часто меняется. Основные проблемы, которые он решает:
		* Избыточные условные операторы: В случаях, когда класс содержит множество условных операторов (if или switch), определяющих поведение в зависимости от текущего состояния.
		* Поддержка и расширение кода: Когда количество состояний и связанных с ними логик увеличивается, условные операторы становятся трудными для поддержки и расширения.
		* Разделение логики: Перенос логики, связанной с различными состояниями, в отдельные классы, что делает код более читаемым и управляемым.

	2) Плюсы и минусы:
		Плюсы:
			+ Избавление от условных операторов: Уменьшение количества больших условных операторов в коде.
			+ Локализация изменений: Логика, связанная с состояниями, сконцентрирована в соответствующих классах, что упрощает внесение изменений и отладку.
			+ Читаемость и поддерживаемость: Улучшение читаемости и поддерживаемости кода за счет разделения логики на отдельные классы.
		Минусы:
			- Сложность: Увеличение количества классов может усложнить структуру программы.
			- Избыточность: В случаях с малым количеством состояний и редкими изменениями, использование паттерна может быть избыточным и усложнить код.

	3) Реальные примеры использования на практике:
		* Управление документами: В системах управления документами, где документ может находиться в разных состояниях (черновик, модерация, опубликован), каждое из которых определяет поведение метода "опубликовать".
		* Музыкальные проигрыватели: В приложениях для воспроизведения музыки, где разные состояния (проигрывание, пауза, заблокировано) изменяют функциональность кнопок управления.
		* Навигационные системы: В навигационных системах, которые меняют свое поведение в зависимости от состояния (планирование маршрута, в пути, остановка).

	4) Как связан с другими паттернами
		* Мост: Оба паттерна используют композицию для делегирования работы другим объектам, но Мост фокусируется на разделении абстракции и реализации для независимого изменения.
		* Стратегия: Паттерн Стратегия позволяет изменять алгоритмы, не зная о существовании других стратегий. В отличие от него, Состояние позволяет состояниям знать друг о друге и инициировать переходы.
		* Адаптер: Адаптер также использует композицию, но для того, чтобы объекты с несовместимыми интерфейсами могли работать вместе.
		* Шаблонный метод: Оба паттерна позволяют изменять поведение объекта, но Шаблонный метод работает на уровне классов, предоставляя частично реализованный алгоритм, в то время как Состояние работает на уровне объектов.

*/

import "fmt"

type State interface {
	setContext(context Human)
	changeState(state State)
	doSomething()
	preformAction()
}

type Hungry struct {
	humanContext *Human
}

func (s *Hungry) setContext(context Human) {
	*s.humanContext = context
}

func (s *Hungry) changeState(state State) {
	s.humanContext.state = state
}

func (s *Hungry) doSomething() {
	fmt.Println("Human eat")
}

func (s *Hungry) preformAction() {
	fmt.Println("Human clean plants")
}

type Thirsty struct {
	humanContext *Human
}

func (s *Thirsty) setContext(context Human) {
	*s.humanContext = context
}

func (s *Thirsty) changeState(state State) {
	s.humanContext.state = state
}

func (s *Thirsty) doSomething() {
	fmt.Println("Human drink")
}

func (s *Thirsty) preformAction() {
	fmt.Println("Human clean bottle")
}

type Sleepy struct {
	humanContext *Human
}

func (s *Sleepy) setContext(context Human) {
	*s.humanContext = context
}

func (s *Sleepy) changeState(state State) {
	s.humanContext.state = state
}

func (s *Sleepy) doSomething() {
	fmt.Println("Human sleep")
}

func (s *Sleepy) preformAction() {
	fmt.Println("Human lies")
}

type Human struct {
	state State
}

func (h *Human) Human(state State) {
	h.state = state
}

func (h *Human) changeState(state State) {
	h.state = state
}

func (h *Human) doSomething() {
	h.state.doSomething()
}

func (h *Human) preformAction() {
	h.state.preformAction()
}

// func main() {
// 	stateHungry := &Hungry{}
// 	human := &Human{state: stateHungry}
// 	human.doSomething()
// 	human.preformAction()

// 	stateThirsty := &Thirsty{}
// 	human.changeState(stateThirsty)
// 	human.doSomething()
// 	human.preformAction()

// 	stateSleepy := &Sleepy{}
// 	human.changeState(stateSleepy)
// 	human.doSomething()
// 	human.preformAction()

// }
