package titles

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common"
	"common/convert"
	"common/data/drivers/postgres"
	"common/data/model"
	"common/data/queriers"
)

type titles struct {
	log *logrus.Entry
	ext sqlx.ExtContext

	expr sq.Sqlizer

	postgres.Inserter[model.Title]
	postgres.Selector[model.Title]
	postgres.Updater[model.UpdateTitleParams, model.Title]
}

func New(ext sqlx.ExtContext, log *logrus.Entry) queriers.TitlesProvider {
	var entity model.Title
	titlesColumns := model.PrependTableName(entity.TableName(), model.Columns(entity, false))
	return &titles{
		log: log.WithField("provider", "titles"),
		ext: ext,

		Inserter: postgres.NewInserter[model.Title](ext, log),
		Selector: postgres.NewSelector[model.Title](ext, log, titlesColumns),
		Updater:  postgres.NewUpdater[model.UpdateTitleParams, model.Title](ext, log),

		expr: common.BasicSqlizer,
	}
}

func (t titles) ByStatus(status string) queriers.TitlesProvider {
	t.expr = sq.And{t.expr, sq.Eq{"titles.status": status}}
	return t
}

func (t titles) InsertUniqueBatch(ctx context.Context, entities []model.Title) error {
	if len(entities) == 0 {
		return nil
	}

	sql := `
		INSERT INTO titles (title, summary, hash, url, release_date, status) VALUES (:title, :summary, :hash, :url, :release_date, :status) 
        	ON CONFLICT (hash) DO NOTHING RETURNING *`

	rows, err := sqlx.NamedQueryContext(ctx, t.ext, t.ext.Rebind(sql), entities)
	if err != nil {
		return errors.Wrap(err, "failed to insert entity into table: titles")
	}

	idx := 0
	for rows.Next() {
		err := rows.StructScan(&entities[idx])
		if err != nil {
			return errors.Wrap(err, "failed to scan entity")
		}
		idx++
	}
	return nil
}

func (t titles) Update(ctx context.Context, title model.UpdateTitleParams) ([]model.Title, error) {
	title.UpdatedAt = convert.ToPtr(common.CurrentTimestamp())

	t.Updater = t.Updater.WithExpr(t.expr)
	return t.Updater.Update(ctx, title)
}
