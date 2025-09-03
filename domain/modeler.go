package domain

type Modeler interface {
	Copy() (interface{}, error) // структура копирует себя и выдает ссылку на копию с массивами и другими данными
	Model() Model
}
