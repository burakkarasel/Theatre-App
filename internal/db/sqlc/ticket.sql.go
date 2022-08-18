// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: ticket.sql

package db

import (
	"context"
)

const createTicket = `-- name: CreateTicket :one
INSERT INTO tickets(movie_id, ticket_owner, child, adult, total)
VALUES($1, $2, $3, $4, $5)
RETURNING id, movie_id, ticket_owner, child, adult, total, created_at
`

type CreateTicketParams struct {
	MovieID     int64  `json:"movie_id"`
	TicketOwner string `json:"ticket_owner"`
	Child       int16  `json:"child"`
	Adult       int16  `json:"adult"`
	Total       int64  `json:"total"`
}

func (q *Queries) CreateTicket(ctx context.Context, arg CreateTicketParams) (Ticket, error) {
	row := q.db.QueryRowContext(ctx, createTicket,
		arg.MovieID,
		arg.TicketOwner,
		arg.Child,
		arg.Adult,
		arg.Total,
	)
	var i Ticket
	err := row.Scan(
		&i.ID,
		&i.MovieID,
		&i.TicketOwner,
		&i.Child,
		&i.Adult,
		&i.Total,
		&i.CreatedAt,
	)
	return i, err
}

const deleteTickets = `-- name: DeleteTickets :exec
DELETE FROM tickets
WHERE id = $1
`

func (q *Queries) DeleteTickets(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTickets, id)
	return err
}

const getTicket = `-- name: GetTicket :one
SELECT id, movie_id, ticket_owner, child, adult, total, created_at
FROM tickets
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetTicket(ctx context.Context, id int64) (Ticket, error) {
	row := q.db.QueryRowContext(ctx, getTicket, id)
	var i Ticket
	err := row.Scan(
		&i.ID,
		&i.MovieID,
		&i.TicketOwner,
		&i.Child,
		&i.Adult,
		&i.Total,
		&i.CreatedAt,
	)
	return i, err
}

const listTickets = `-- name: ListTickets :many
SELECT id, movie_id, ticket_owner, child, adult, total, created_at
FROM tickets
WHERE ticket_owner = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListTicketsParams struct {
	TicketOwner string `json:"ticket_owner"`
	Limit       int32  `json:"limit"`
	Offset      int32  `json:"offset"`
}

func (q *Queries) ListTickets(ctx context.Context, arg ListTicketsParams) ([]Ticket, error) {
	rows, err := q.db.QueryContext(ctx, listTickets, arg.TicketOwner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Ticket{}
	for rows.Next() {
		var i Ticket
		if err := rows.Scan(
			&i.ID,
			&i.MovieID,
			&i.TicketOwner,
			&i.Child,
			&i.Adult,
			&i.Total,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
