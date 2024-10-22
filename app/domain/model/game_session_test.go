package model

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/averak/gamebox/testutils/faker"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	GameIDDummy1 GameID = 1
	GameIDDummy2 GameID = 2
)

func TestNewGameSession(t *testing.T) {
	type args struct {
		id         uuid.UUID
		userID     uuid.UUID
		gameID     GameID
		status     GameStatus
		result     GameResult
		wager      Coins
		payout     Coins
		startedAt  time.Time
		finishedAt time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    GameSession
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "wager > 0 の場合 => ゲームセッションを作成できる",
			args: args{
				wager: 1,
			},
			want: GameSession{
				Wager: 1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "wager <= 0 の場合 => エラー",
			args: args{
				wager: 0,
			},
			want:    GameSession{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGameSession(tt.args.id, tt.args.userID, tt.args.gameID, tt.args.status, tt.args.result, tt.args.wager, tt.args.payout, tt.args.startedAt, tt.args.finishedAt)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGameSession_FinishPlaying(t *testing.T) {
	now := time.Now()

	type fields struct {
		Status GameStatus
		Result GameResult
		Wager  Coins
	}
	type args struct {
		result GameResult
		now    time.Time
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		want            Coins
		wantGameSession *GameSession
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "勝利した場合、掛け金が2倍になる",
			fields: fields{
				Status: GameStatusPlaying,
				Wager:  1,
			},
			args: args{
				result: GameResultWin,
				now:    now,
			},
			want: 2,
			wantGameSession: &GameSession{
				Status:     GameStatusFinished,
				Result:     GameResultWin,
				Wager:      1,
				Payout:     2,
				FinishedAt: now,
			},
			wantErr: assert.NoError,
		},
		{
			name: "敗北した場合、掛け金が没収される",
			fields: fields{
				Status: GameStatusPlaying,
				Wager:  1,
			},
			args: args{
				result: GameResultLose,
				now:    now,
			},
			want: 0,
			wantGameSession: &GameSession{
				Status:     GameStatusFinished,
				Result:     GameResultLose,
				Wager:      1,
				Payout:     0,
				FinishedAt: now,
			},
			wantErr: assert.NoError,
		},
		{
			name: "引き分けの場合、掛け金がそのまま返却される",
			fields: fields{
				Status: GameStatusPlaying,
				Wager:  1,
			},
			args: args{
				result: GameResultDraw,
				now:    now,
			},
			want: 1,
			wantGameSession: &GameSession{
				Status:     GameStatusFinished,
				Result:     GameResultDraw,
				Wager:      1,
				Payout:     1,
				FinishedAt: now,
			},
			wantErr: assert.NoError,
		},
		{
			name: "ゲームが終了済みの場合 => エラー",
			fields: fields{
				Status: GameStatusFinished,
				Wager:  1,
			},
			args: args{
				result: GameResultWin,
				now:    now,
			},
			wantGameSession: &GameSession{
				Status: GameStatusFinished,
				Wager:  1,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrGameAlreadyFinished)
			},
		},
		{
			name: "不正なゲーム結果が指定された場合 => エラー",
			fields: fields{
				Status: GameStatusPlaying,
				Wager:  1,
			},
			args: args{
				result: GameResultUnknown,
				now:    now,
			},
			wantGameSession: &GameSession{
				Status: GameStatusPlaying,
				Wager:  1,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GameSession{
				Status: tt.fields.Status,
				Result: tt.fields.Result,
				Wager:  tt.fields.Wager,
			}
			got, err := g.FinishPlaying(tt.args.result, tt.args.now)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantGameSession, g)
		})
	}
}

func TestGameSession_IsPlaying(t *testing.T) {
	type fields struct {
		Status GameStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "プレイ中の場合 => true",
			fields: fields{
				Status: GameStatusPlaying,
			},
			want: true,
		},
		{
			name: "終了済みの場合 => false",
			fields: fields{
				Status: GameStatusFinished,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GameSession{
				Status: tt.fields.Status,
			}
			assert.Equal(t, tt.want, g.IsPlaying())
		})
	}
}

func TestNewGameSessionService(t *testing.T) {
	type args struct {
		userID          uuid.UUID
		playingSessions []GameSession
	}
	tests := []struct {
		name    string
		args    args
		want    GameSessionStartService
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "正常系",
			args: args{
				playingSessions: []GameSession{
					{
						UserID: faker.UUIDv5("u1"),
						GameID: GameIDDummy2,
						Status: GameStatusPlaying,
					},
					{
						UserID: faker.UUIDv5("u1"),
						GameID: GameIDDummy1,
						Status: GameStatusPlaying,
					},
				},
			},
			want: GameSessionStartService{
				playingSessions: []GameSession{
					{
						UserID: faker.UUIDv5("u1"),
						GameID: GameIDDummy2,
						Status: GameStatusPlaying,
					},
					{
						UserID: faker.UUIDv5("u1"),
						GameID: GameIDDummy1,
						Status: GameStatusPlaying,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "異常系: 同じゲームが複数プレイ中",
			args: args{
				playingSessions: []GameSession{
					{
						UserID: faker.UUIDv5("u1"),
						GameID: GameIDDummy1,
						Status: GameStatusPlaying,
					},
					{
						UserID: faker.UUIDv5("u1"),
						GameID: GameIDDummy1,
						Status: GameStatusPlaying,
					},
				},
			},
			want:    GameSessionStartService{},
			wantErr: assert.Error,
		},
		{
			name: "異常系: プレイ中ではないゲームセッションが指定された",
			args: args{
				playingSessions: []GameSession{
					{
						UserID: faker.UUIDv5("u1"),
						GameID: GameIDDummy1,
						Status: GameStatusFinished,
					},
				},
			},
			want:    GameSessionStartService{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGameSessionService(context.TODO(), tt.args.userID, tt.args.playingSessions)
			if !tt.wantErr(t, err, fmt.Sprintf("NewGameSessionService(%v)", tt.args.playingSessions)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewGameSessionService(%v)", tt.args.playingSessions)
		})
	}
}

func TestGameSessionService_StartPlaying(t *testing.T) {
	now := time.Now()

	type fields struct {
		userID          uuid.UUID
		playingSessions []GameSession
	}
	type args struct {
		id     uuid.UUID
		gameID GameID
		wager  Coins
		now    time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    GameSession
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "同じゲームがプレイ中でない場合 => ゲームセッションを開始する",
			fields: fields{
				userID: faker.UUIDv5("u1"),
				playingSessions: []GameSession{
					{
						UserID: faker.UUIDv5("u1"),
						GameID: GameIDDummy2,
						Status: GameStatusPlaying,
					},
				},
			},
			args: args{
				id:     faker.UUIDv5("gs1"),
				gameID: GameIDDummy1,
				wager:  100,
				now:    now,
			},
			want: GameSession{
				ID:        faker.UUIDv5("gs1"),
				UserID:    faker.UUIDv5("u1"),
				GameID:    GameIDDummy1,
				Status:    GameStatusPlaying,
				Wager:     100,
				StartedAt: now,
			},
			wantErr: assert.NoError,
		},
		{
			name: "同じゲームがプレイ中の場合 => エラー",
			fields: fields{
				playingSessions: []GameSession{
					{
						UserID: faker.UUIDv5("u1"),
						GameID: GameIDDummy1,
						Status: GameStatusPlaying,
					},
				},
			},
			args: args{
				id:     faker.UUIDv5("gs1"),
				gameID: GameIDDummy1,
				wager:  100,
				now:    now,
			},
			want: GameSession{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrGameAlreadyPlaying)
			},
		},
		{
			name: "不正な賭け金 => エラー",
			fields: fields{
				userID:          faker.UUIDv5("u1"),
				playingSessions: []GameSession{},
			},
			args: args{
				id:     faker.UUIDv5("gs1"),
				gameID: GameIDDummy1,
				wager:  0,
				now:    now,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrGameWagerIsInvalid)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GameSessionStartService{
				userID:          tt.fields.userID,
				playingSessions: tt.fields.playingSessions,
			}
			got, err := s.StartPlaying(tt.args.id, tt.args.gameID, tt.args.wager, tt.args.now)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
