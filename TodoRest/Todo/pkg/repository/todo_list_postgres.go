package repository

import (
	todo "Todo"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
)

type TodoListPostgres struct { // struct to connect to db of todolist
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{
		db: db,
	}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin() // ready for transactions (start point for rollback)
	if err != nil {
		return 0, err
	}

	var id int // returning id of new list
	createListQuery := fmt.Sprintf("insert into %s (title, description) values ($1, $2) returning id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback() // if there was a mistake -> rolling back where nothing has been changed
		return 0, err
	}

	// setting dependency for many to many connections
	createUserListQuery := fmt.Sprintf("insert into %s (user_id, list_id) values ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUserListQuery, userId, id)
	if err != nil {
		tx.Rollback() // if there was a mistake -> rolling back where nothing has been changed
		return 0, err
	}

	return id, tx.Commit() // committing all changes to the table because everything was ok
}

// block of sending sql queries to db
func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList // list to parse lists in

	query := fmt.Sprintf("select tl.id, tl.title, tl.description from %s tl inner join %s ul on tl.id = ul.list_id where ul.user_id = $1",
		todoListsTable, usersListsTable)      // query to get all lists this user
	err := r.db.Select(&lists, query, userId) // sending error upper

	return lists, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	query := fmt.Sprintf("select tl.id, tl.title, tl.description from %s tl inner join %s ul on tl.id = ul.list_id where ul.user_id = $1 and ul.list_id = $2",
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)
	return list, err
}

func (r *TodoListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("delete from %s tl using %s ul where tl.id=ul.list_id and ul.user_id=$1 and ul.list_id=$2",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListPostgres) Update(userId, listId int, input todo.UpdateListInput) error {
	setValues := make([]string, 0) // names of the values
	args := make([]interface{}, 0) // interface slice for new arguments
	argId := 1                     // needed to properly number $ in the sql query later

	if input.Title != nil { // title $1
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil { // description $1 or $2
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, listId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)

	return err
}
