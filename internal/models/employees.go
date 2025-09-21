package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Employee struct {
	LastName  string
	FirstName string
	Title     string
	City      string
	Email     string
}

type EmployeeModel struct {
	DB *sql.DB
}

// This will return the 10 most recently created snippets.
func (m *EmployeeModel) List() ([]Employee, error) {
	// Write the SQL statement we want to execute.
	stmt := `select LastName, FirstName, Title, City, Email from employees`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []Employee

	for rows.Next() {
		var e Employee
		err = rows.Scan(&e.LastName, &e.FirstName, &e.Title, &e.City, &e.Email)
		if err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return employees, nil
}
