package echo_repoimpl

import (
	"context"

	"github.com/averak/gamebox/app/adapter/dao"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/averak/gamebox/app/infrastructure/trace"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Repository struct{}

func NewRepository() repository.EchoRepository {
	return &Repository{}
}

func (r Repository) Save(ctx context.Context, tx transaction.Transaction, echos ...model.Echo) error {
	ctx, span := trace.StartSpan(ctx, "echo_repoimpl.Save")
	defer span.End()

	dtos := make([]*dao.Echo, len(echos))
	for i, echo := range echos {
		dtos[i] = &dao.Echo{
			ID:        echo.ID.String(),
			Message:   echo.Message,
			Timestamp: echo.Timestamp,
		}
	}
	_, err := dao.EchoSlice(dtos).UpsertAll(ctx, tx, true, dao.EchoPrimaryKeyColumns, boil.Infer(), boil.Infer())
	if err != nil {
		return err
	}
	return nil
}
