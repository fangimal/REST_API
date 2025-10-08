package postgresql

import (
	"REST_API/internal/author/model"
	"REST_API/internal/author/storage"
	"REST_API/pkg/client/postgresql"
	"REST_API/pkg/logging"
	"context"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) Create(ctx context.Context, author *model.Author) error {
	q := `
			INSERT INTO author (name)  
			VALUES ($1)
			RETURNING id
		`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, author.Name).Scan(&author.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, sortOptions storage.SortOptions) (u []model.Author, err error) {
	qb := sq.Select("id, name, age, is_alive, created_at").From("public.author")
	if sortOptions != nil {
		qb = qb.OrderBy(sortOptions.GetOrderBy())
	}
	sql, i, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(sql)))

	rows, err := r.client.Query(ctx, sql, i...)
	if err != nil {
		return nil, err
	}

	authors := make([]model.Author, 0)

	for rows.Next() {
		var ath model.Author

		err = rows.Scan(&ath.ID, &ath.Name, &ath.Age, &ath.IsAlive, &ath.CreatedAt)
		if err != nil {
			return nil, err
		}

		authors = append(authors, ath)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (model.Author, error) {
	q := `
			SELECT id, name FROM public.author WHERE id = $1
		`
	var ath model.Author
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	err := r.client.QueryRow(ctx, q, id).Scan(&ath.ID, &ath.Name)
	if err != nil {
		return model.Author{}, err
	}

	return ath, nil
}

func (r *repository) Update(ctx context.Context, user model.Author) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) storage.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
