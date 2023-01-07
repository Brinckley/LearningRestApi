package repository

import (
	todo "Todo"
	"github.com/jmoiron/sqlx"
)

type Authorization interface { // интерфейс для обработки авторизации (базовые функции в виде добавления и создания юзера)
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}

type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input todo.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, input todo.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo.TodoItem, error)
	GetById(userId, itemId int) (todo.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input todo.UpdateItemInput) error
}

type Repository struct { // имплементируем все интерфейсы в структуру Репозиторий
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),     // конструктор на основе sqlx.Db, создаем именно репозиторий для авторизаций
		TodoList:      NewTodoListPostgres(db), // репозиторий для списков дел
		TodoItem:      NewTodoItemPostgres(db), // connecting repository for items
	}
}
