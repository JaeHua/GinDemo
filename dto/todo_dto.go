package dto

import "GinVue/model"

type TodoDto struct {
	Title     string `json:"title"`
	Telephone string `json:"telephone"`
}

func ToTodoDto(todo model.Todo) TodoDto {
	return TodoDto{
		Title:     todo.Title,
		Telephone: todo.Telephone,
	}
}
