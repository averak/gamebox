package aop

import (
	"context"
	"errors"
	"time"

	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/app/domain/repository"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/averak/gamebox/app/infrastructure/connect/error_response"
	"github.com/averak/gamebox/app/infrastructure/connect/mdval"
	"github.com/averak/gamebox/app/infrastructure/session"
	"github.com/averak/gamebox/protobuf/api/api_errors"
	"github.com/google/uuid"
)

func checkSession(ctx context.Context, conf *config.Config, repo repository.UserRepository, conn transaction.Connection, incomingMD mdval.IncomingMD, now time.Time) (*model.User, error) {
	// 優先度は デバッグ用ヘッダ > セッショントークン とする。
	if conf.GetDebug() {
		spoofingUserID, ok := incomingMD.Get(mdval.DebugSpoofingUserIDKey)
		if ok {
			user, err := setupSpoofingUser(ctx, conn, repo, uuid.MustParse(spoofingUserID))
			if err != nil {
				return nil, err
			}
			return &user, nil
		}
	}

	sessionToken, ok := incomingMD.Get(mdval.SessionTokenKey)
	if !ok {
		return nil, error_response.New(api_errors.ErrorCode_COMMON_INVALID_SESSION, api_errors.ErrorSeverity_ERROR_SEVERITY_WARNING, "session token not found")
	}
	sess, err := session.DecodeSessionToken(sessionToken, []byte(conf.GetApiServer().GetSession().GetSecretKey()), now)
	if err != nil {
		return nil, error_response.New(api_errors.ErrorCode_COMMON_INVALID_SESSION, api_errors.ErrorSeverity_ERROR_SEVERITY_WARNING, "invalid session")
	}

	var user model.User
	err = conn.BeginRoTransaction(ctx, func(ctx context.Context, tx transaction.Transaction) error {
		user, err = repo.Get(ctx, tx, sess.PrincipalID)
		return err
	})
	if err != nil {
		return nil, error_response.New(api_errors.ErrorCode_COMMON_INVALID_SESSION, api_errors.ErrorSeverity_ERROR_SEVERITY_WARNING, "failed to get user")
	}
	return &user, nil
}

// setupSpoofingUser はデバッグ用ヘッダで指定されたユーザーIDが存在しない場合に、ユーザの初期化を行います。
// あくまでもデバッグ機能なので、最低限 API を呼び出せる状態のユーザとして初期化します。
func setupSpoofingUser(ctx context.Context, conn transaction.Connection, repo repository.UserRepository, userID uuid.UUID) (model.User, error) {
	var user model.User
	err := conn.BeginRwTransaction(ctx, func(ctx context.Context, tx transaction.Transaction) error {
		var err error
		user, err = repo.Get(ctx, tx, userID)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				user = model.NewUser(userID, model.UserStatusActive)
				return repo.Save(ctx, tx, user)
			}
			return err
		}
		return nil
	})
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
