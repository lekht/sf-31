package postgres

import (
	"context"
	"sf-31/pkg/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresDB struct {
	db *pgxpool.Pool
}

const (
	sqlGetAll     = `SELECT * FROM posts ORDER BY id;`
	sqlAddPost    = `INSERT INTO posts(title, content) VALUES ($1,$2) RETURNING id;`
	sqlUpdatePost = `UPDATE posts SET title = $1, content = $2 WHERE id = $3;`
	sqlDeletePost = `DELETE FROM posts WHERE id = $1;`
)

func New(connstr string) (*PostgresDB, error) {
	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	p := PostgresDB{
		db: db,
	}
	return &p, nil
}

func (p *PostgresDB) Posts() ([]storage.Post, error) {
	var posts []storage.Post
	rows, err := p.db.Query(context.Background(), sqlGetAll)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ps storage.Post
		err = rows.Scan(
			&ps.ID,
			&ps.AuthorID,
			&ps.Title,
			&ps.Content,
			&ps.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, ps)
	}
	return posts, rows.Err()
}

func (p *PostgresDB) AddPost(ps storage.Post) error {
	_, err := p.db.Exec(context.Background(),
		sqlAddPost,
		ps.Title,
		ps.Content,
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) UpdatePost(ps storage.Post) error {
	_, err := p.db.Exec(context.Background(),
		sqlUpdatePost,
		ps.Title,
		ps.Content,
		ps.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) DeletePost(ps storage.Post) error {
	_, err := p.db.Exec(context.Background(),
		sqlDeletePost,
		ps.ID,
	)
	if err != nil {
		return nil
	}
	return nil
}
