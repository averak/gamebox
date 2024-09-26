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

func TestGameSession_FinishPlaying(t *testing.T) {
	now := time.Now()

	type fields struct {
		Status GameStatus
		Result GameResult
		Wager  int
	}
	type args struct {
		result GameResult
		wallet Wallet
		now    time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    GameSession
		want1   Wallet
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "勝利した場合、掛け金が2倍になる",
			fields: fields{
				Status: GameStatusPlaying,
				Wager:  1,
			},
			args: args{
				result: GameResultWin,
				wallet: Wallet{
					Balance: 0,
				},
				now: now,
			},
			want: GameSession{
				Status:     GameStatusFinished,
				Result:     GameResultWin,
				Wager:      1,
				Payout:     2,
				FinishedAt: now,
			},
			want1: Wallet{
				Balance: 2,
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
				wallet: Wallet{
					Balance: 0,
				},
				now: now,
			},
			want: GameSession{
				Status:     GameStatusFinished,
				Result:     GameResultLose,
				Wager:      1,
				Payout:     0,
				FinishedAt: now,
			},
			want1: Wallet{
				Balance: 0,
			},
			wantErr: assert.NoError,
		},
		{
			name: "ゲームが終了済みの場合 => エラー",
			fields: fields{
				Status: GameStatusFinished,
			},
			args: args{
				result: GameResultWin,
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
				wallet: Wallet{
					Balance: 0,
				},
				now: now,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GameSession{
				Status: tt.fields.Status,
				Result: tt.fields.Result,
				Wager:  tt.fields.Wager,
			}
			got, got1, err := s.FinishPlaying(tt.args.result, tt.args.wallet, tt.args.now)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
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
		want    GameSessionService
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
			want: GameSessionService{
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
			want:    GameSessionService{},
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
			want:    GameSessionService{},
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
		gameID GameID
		wager  int
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
				gameID: GameIDDummy1,
				wager:  100,
				now:    now,
			},
			want: GameSession{
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
				gameID: GameIDDummy1,
				wager:  100,
				now:    now,
			},
			want: GameSession{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrGameAlreadyPlaying)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GameSessionService{
				userID:          tt.fields.userID,
				playingSessions: tt.fields.playingSessions,
			}
			got, err := s.StartPlaying(tt.args.gameID, tt.args.wager, tt.args.now)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}