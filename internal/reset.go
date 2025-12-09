package internal

import (
	"context"
	"fmt"
)

func HandlerReset(s *State, _ Command) error {
	rowsDeleted, err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not delete all users from the db", err)
	}

	fmt.Printf("Successfuly deleted %d rows from USERS table\n", rowsDeleted)
	return nil
}
