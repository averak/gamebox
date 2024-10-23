package aop

import (
	"context"
	"errors"
	"time"

	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/core/game_context"
	"github.com/averak/gamebox/app/infrastructure/connect/interceptor"
	"github.com/averak/gamebox/app/infrastructure/connect/mdval"
	"github.com/google/uuid"
)

type preconditionParams struct {
	RequestID      uuid.UUID
	IdempotencyKey uuid.UUID
	// リクエストを受け取った現実世界での時刻
	RequestedTime time.Time
	// サーバが現在時刻として認識する時刻 (現実世界と乖離しうる)
	Now time.Time
}

func (p preconditionParams) GameContext() game_context.GameContext {
	return game_context.NewGameContext(
		p.IdempotencyKey,
		p.Now,
	)
}

func fixPreconditionParams(ctx context.Context, incomingMD mdval.IncomingMD) (preconditionParams, error) {
	builder := newPreconditionParamsBuilder()
	builder.requestID(uuid.New())
	fixIdempotencyKey(incomingMD, builder)
	err := fixTimestamp(config.Get(), incomingMD, builder, time.Now())
	if err != nil {
		return preconditionParams{}, err
	}

	interceptor.AddLogHint(ctx, "requestID", builder.raw.RequestID)
	interceptor.AddLogHint(ctx, "gameContext", builder.raw.GameContext())
	return builder.build()
}

func fixIdempotencyKey(incomingMD mdval.IncomingMD, builder *preconditionParamsBuilder) {
	idempotencyKey, ok := incomingMD.Get(mdval.IdempotencyKey)
	if ok {
		ik, err := uuid.Parse(idempotencyKey)
		if err != nil {
			builder.idempotencyKey(uuid.New())
		} else {
			builder.idempotencyKey(ik)
		}
	} else {
		builder.idempotencyKey(uuid.New())
	}
}

func fixTimestamp(conf *config.Config, incomingMD mdval.IncomingMD, builder *preconditionParamsBuilder, now time.Time) error {
	builder.requestedTime(now)
	builder.now(now)
	if conf.GetDebug() {
		adjustedTimeStr, ok := incomingMD.Get(mdval.DebugAdjustedTimeKey)
		if ok {
			adjustedTime, err := time.Parse(time.RFC3339, adjustedTimeStr)
			if err != nil {
				return err
			}
			builder.now(adjustedTime)
		}
	}
	return nil
}

// preconditionParams の設定漏れがないか検証するために、ビルダーを定義しています。
type preconditionParamsBuilder struct {
	raw preconditionParams
}

func newPreconditionParamsBuilder() *preconditionParamsBuilder {
	return &preconditionParamsBuilder{}
}

func (b preconditionParamsBuilder) build() (preconditionParams, error) {
	if b.raw.RequestID == uuid.Nil {
		return preconditionParams{}, errors.New("requestID is not set")
	}
	if b.raw.IdempotencyKey == uuid.Nil {
		return preconditionParams{}, errors.New("idempotencyKey is not set")
	}
	if b.raw.RequestedTime.IsZero() {
		return preconditionParams{}, errors.New("requestedTime is not set")
	}
	if b.raw.Now.IsZero() {
		return preconditionParams{}, errors.New("now is not set")
	}
	return b.raw, nil
}

func (b *preconditionParamsBuilder) requestID(id uuid.UUID) *preconditionParamsBuilder {
	b.raw.RequestID = id
	return b
}

func (b *preconditionParamsBuilder) idempotencyKey(key uuid.UUID) *preconditionParamsBuilder {
	b.raw.IdempotencyKey = key
	return b
}

func (b *preconditionParamsBuilder) requestedTime(t time.Time) *preconditionParamsBuilder {
	b.raw.RequestedTime = t
	return b
}

func (b *preconditionParamsBuilder) now(t time.Time) *preconditionParamsBuilder {
	b.raw.Now = t
	return b
}
