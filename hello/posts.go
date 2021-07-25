package hello

import (
	"context"
	"fmt"

	"encore.app/firebaseauth"

	"encore.dev/beta/auth"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
)

type UserPost struct {
	Title    string
	Content  string
	Id       int
	Username string
}

type PostsResponse struct {
	Posts []*UserPost
}

type AddPostRequest struct {
	Title   string
	Content string
}

type AddPostResponse struct {
	Id int
}

// encore:api public
func GetPosts(ctx context.Context) (*PostsResponse, error) {
	rows, err := sqldb.Query(ctx, `
	select public.post.id, public.post.title, public.post.content, public.user.name from public.post, public.user 
	where public.post.user_id = public.user.id;
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []*UserPost{}
	for rows.Next() {
		var post UserPost
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Username); err != nil {
			return nil, fmt.Errorf("could not scan: %v", err)
		}
		posts = append(posts, &post)
	}

	if err != nil {
		return nil, err
	}
	return &PostsResponse{Posts: posts}, nil
}

// encore:api auth
func GetMyPosts(ctx context.Context) (*PostsResponse, error) {
	user_id, _ := auth.UserID()
	rows, err := sqldb.Query(ctx, `
	select public.post.id, public.post.title, public.post.content, public.user.name from public.post, public.user 
	where public.post.user_id = public.user.id and public.post.user_id = $1;
    `, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []*UserPost{}
	for rows.Next() {
		var post UserPost
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Username); err != nil {
			return nil, fmt.Errorf("could not scan: %v", err)
		}
		posts = append(posts, &post)
	}

	if err != nil {
		return nil, err
	}
	return &PostsResponse{Posts: posts}, nil
}

// encore:api auth
func AddPost(ctx context.Context, params *AddPostRequest) (*AddPostResponse, error) {
	var id int
	user_id, _ := auth.UserID()
	authdata, _ := auth.Data().(*firebaseauth.AuthData)
	rlog.Info("adding post", "name", authdata.Name)
	err := sqldb.QueryRow(ctx, `
        INSERT INTO public.post (title, content, user_id)
        VALUES ($1, $2, $3)
        RETURNING id
    `, params.Title, params.Content, user_id).Scan(&id)
	_, err = sqldb.Exec(ctx, `UPDATE public.user SET name = $1 WHERE id = $2;`, authdata.Name, user_id)
	if err != nil {
		return nil, err
	}
	return &AddPostResponse{Id: id}, nil
}
