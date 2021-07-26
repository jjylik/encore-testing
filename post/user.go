package post

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

type User struct {
	Name  string `json:"name"`
	Id    string `json:"id"`
}

type GetUsersResponse struct {
	Users []*User
}


//encore:api public
func GetUsers(ctx context.Context) (*GetUsersResponse, error) {
	rows, err := sqldb.Query(ctx, `
        select public.user.id, public.user.name from public.user;
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []*User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			return nil, fmt.Errorf("could not scan: %v", err)
		}
		users = append(users, &user)
	}

	if err != nil {
		return nil, err
	}
	return &GetUsersResponse{Users: users}, nil
}
