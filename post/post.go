package post

import (
	"context"
	"fmt"

	"encore.app/firebaseauth"

	"encore.dev/beta/auth"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
)

type UserPost struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type PostsResponse struct {
	Posts []*UserPost `json:"posts"`
}

type AddPostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type AddPostResponse struct {
	Id int `json:"id"`
}

func toPostListResponse(rows *sqldb.Rows) (*PostsResponse, error) {
	posts := []*UserPost{}
	for rows.Next() {
		var post UserPost
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Username); err != nil {
			return nil, fmt.Errorf("could not scan: %v", err)
		}
		posts = append(posts, &post)
	}
	return &PostsResponse{Posts: posts}, nil
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
	return toPostListResponse(rows)
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
	return toPostListResponse(rows)
}

// encore:api auth
func AddPost(ctx context.Context, params *AddPostRequest) (*AddPostResponse, error) {
	var id int
	user_id, _ := auth.UserID()
	authdata, _ := auth.Data().(*firebaseauth.AuthData)
	tx, _ := sqldb.Begin(ctx)
	rlog.Info("add post", "name", authdata.Name)
	_, err := sqldb.ExecTx(tx, ctx, `INSERT INTO public.user (id, name) values ($2, $1) ON CONFLICT (id) DO UPDATE SET name = $1;`, authdata.Name, user_id)
	if err != nil {
		return nil, err
	}
	err = sqldb.QueryRowTx(tx, ctx, `
        INSERT INTO public.post (title, content, user_id)
        VALUES ($1, $2, $3)
        RETURNING id
    `, params.Title, params.Content, user_id).Scan(&id)
	if err != nil {
		sqldb.Rollback(tx)
		return nil, err
	}
	sqldb.Commit(tx)
	return &AddPostResponse{Id: id}, nil
}
