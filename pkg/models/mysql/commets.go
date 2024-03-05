package mysql

import (
	"ass1Go/pkg/models"
	"database/sql"
	"errors"
)

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Insert(userID, newsID int, text string) (int, error) {
	stmt := `INSERT INTO comments (user_id, news_id, text)
    VALUES(?, ?, ?)`

	result, err := m.DB.Exec(stmt, userID, newsID, text)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *CommentModel) Get(id int) (*models.Comment, error) {
	stmt := `
        SELECT c.id, c.user_id, c.news_id, c.text, u.name AS user_name
        FROM comments c
        INNER JOIN users u ON c.user_id = u.id
        WHERE c.id = ?
    `
	row := m.DB.QueryRow(stmt, id)
	c := &models.Comment{}
	err := row.Scan(&c.ID, &c.UserID, &c.NewsID, &c.Text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return c, nil
}

func (m *CommentModel) GetAll() ([]*models.Comment, error) {
	stmt := `
        SELECT c.id, c.user_id, c.news_id, c.text, u.name AS user_name
        FROM comments c
        INNER JOIN users u ON c.user_id = u.id
    `
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []*models.Comment{}

	for rows.Next() {
		c := &models.Comment{}
		err = rows.Scan(&c.ID, &c.UserID, &c.NewsID, &c.Text)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *CommentModel) GetAllForSnippetID(snippetID int) ([]*models.Comment, error) {
	// Prepare the SQL statement to retrieve comments for the given snippet ID
	query := `
        SELECT id, user_id, news_id, text
        FROM comments
        WHERE news_id = ?
    `

	// Execute the query and retrieve the rows
	rows, err := m.DB.Query(query, snippetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to store the retrieved comments
	var comments []*models.Comment

	// Iterate through the rows and scan each comment into a Comment struct
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.NewsID, &comment.Text)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return the retrieved comments
	return comments, nil
}
