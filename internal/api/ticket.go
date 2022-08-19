package api

import (
	"database/sql"
	"net/http"

	db "github.com/burakkarasel/Theatre-API/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

// CreateTicketRequest holds the json data of the createTicket
type CreateTicketRequest struct {
	TicketOwner string `json:"ticket_owner" binding:"required,min=6"`
	MovieID     int64  `json:"movie_id" binding:"required,min=1"`
	Total       int64  `json:"total" binding:"required,gt=0"`
	Child       int16  `json:"child" binding:"min=0"`
	Adult       int16  `json:"adult" binding:"min=0"`
}

// CreateTicketResponse holds the data for createTicket response
type CreateTicketResponse struct {
	Ticket db.Ticket `json:"ticket"`
	Movie  db.Movie  `json:"movie"`
}

func (server *Server) createTicket(ctx *gin.Context) {
	// first i check for bindings
	var req CreateTicketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// then i check that i am not creating a ticket for no one
	if req.Adult == 0 && req.Child == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrInvalidTicket))
		return
	}

	// then i create args for the DB operation
	arg := db.CreateTicketParams{
		MovieID:     req.MovieID,
		TicketOwner: req.TicketOwner,
		Total:       req.Total,
		Child:       req.Child,
		Adult:       req.Adult,
	}

	// then i get the movie of the ticket and check for error
	m, err := server.store.GetMovie(ctx, req.MovieID)

	if err != nil {
		if err == sql.ErrNoRows {
			// if error is no rows i return 404 and the error
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		// otherwise i return 500 and the error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// then i create the ticket
	t, err := server.store.CreateTicket(ctx, arg)

	// if any error occurs i return 500 and the error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// if no error occurs i return ok and create ticket response
	ctx.JSON(http.StatusOK, CreateTicketResponse{Ticket: t, Movie: m})
}

// GetTicketRequest holds uri data of the request
type GetTicketRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// GetTicketResponse holds the json data of the response
type GetTicketResponse struct {
	Ticket db.Ticket `json:"ticket"`
	Movie  db.Movie  `json:"movie"`
}

// getTicket takes ID and returns the relevant Ticket
func (server *Server) getTicket(ctx *gin.Context) {
	// first i check for bindings
	var req GetTicketRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// then i get the ticket from DB
	t, err := server.store.GetTicket(ctx, req.ID)

	if err != nil {
		// if err is no rows i return 404 and the error
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		// otherwise i return 500 and the error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// then i get the movie from db
	m, err := server.store.GetMovie(ctx, t.MovieID)

	// if any error occurs i return 500 and the error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// if no error occurs i return OK and get ticket response
	ctx.JSON(http.StatusOK, GetTicketResponse{Movie: m, Ticket: t})
}

// ListTicketRequest holds the query data of the request
type ListTicketsRequest struct {
	TicketOwner string `form:"ticket_owner" bindings:"required,min=6"`
	PageID      int32  `form:"page_id" bindings:"required,min=1"`
	PageSize    int32  `form:"page_size" bindings:"required,min=5,max=10"`
}

func (server *Server) listTickets(ctx *gin.Context) {
	// first i check for the bindings
	var req ListTicketsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// then i create args to call DB func
	arg := db.ListTicketsParams{
		TicketOwner: req.TicketOwner,
		Limit:       req.PageSize,
		Offset:      (req.PageID - 1) * req.PageSize,
	}

	tickets, err := server.store.ListTickets(ctx, arg)

	// if any error occurs i check the error
	if err != nil {
		if err == sql.ErrNoRows {
			// if error is no rows i return 404 and the error
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		// otherwise i return 500 and the error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// then i get each ticket's movie for the response
	var result []GetTicketResponse
	for _, t := range tickets {
		m, err := server.store.GetMovie(ctx, t.MovieID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var res = GetTicketResponse{
			Movie:  m,
			Ticket: t,
		}

		result = append(result, res)
	}

	// then i return the result with OK
	ctx.JSON(http.StatusOK, result)
}

// DeleteTicketRequest holds the uri data of the request
type DeleteTicketRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// deleteTicket deletes the ticket for given ID in DB
func (server *Server) deleteTicket(ctx *gin.Context) {
	// first i check bindings
	var req DeleteTicketRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// then i delete the ticket
	err := server.store.DeleteTicket(ctx, req.ID)

	// if any error occurs i check the error message
	if err != nil {
		// if error is no rows i return 404 and the error
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		// otherwise i return 500 and the error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// if no error occurs i return OK and no data
	ctx.JSON(http.StatusOK, nil)
}
