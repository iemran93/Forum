package database

import "forumProject/internal/models"

func CreateUser(user models.User) error {
	sqlStm := `INSERT INTO users 
	(email, username, password) 
	VALUES (?, ?, ?)
	`

	stm, err := DB.Prepare(sqlStm)
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(user.Email, user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func GetUsers() ([]models.User, error) {
	sqlStm := `SELECT id, 
	email, username,
	password FROM users
	`
	stm, err := DB.Prepare(sqlStm)
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	rows, err := stm.Query()
	if err != nil {
		return nil, err
	}
	users := make([]models.User, 0)
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
