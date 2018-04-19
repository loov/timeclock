package pgdb

import "github.com/loov/timeclock/user"

type Users struct {
	*Database
}

func (db *Database) Users() user.Database { return &Users{db} }

func (db *Users) Register(user user.User, provider string, key []byte) (user.ID, error) {
	return 0, todo
}

func (db *Users) FindCredentials(name string) error {
	return todo
}

func (db *Users) List() ([]user.User, error) {
	rows, err := db.Query(`
		SELECT 
			Users.ID, Users.Alias, Users.Name, Users.Email, Users.Root,
			Roles.Admin, Roles.Accountant, Roles.Worker
		FROM Users
		LEFT JOIN Roles ON Users.ID = Roles.UserID
		ORDER BY Users.Name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		err := rows.Scan(
			&u.ID, &u.Alias, &u.Name, &u.Email, &u.Root,
			&u.Admin, &u.Accountant, &u.Worker,
		)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}
