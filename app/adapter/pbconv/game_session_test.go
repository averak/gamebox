package pbconv

import (
	"testing"
	"time"

	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/protobuf/resource"
	"github.com/averak/gamebox/testutils/faker"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestToGameSessionPb(t *testing.T) {
	now := time.Now()

	type args struct {
		sess model.GameSession
	}
	tests := []struct {
		name string
		args args
		want *resource.GameSession
	}{
		{
			name: "変換できる",
			args: args{
				sess: model.GameSession{
					ID:         faker.UUIDv5("gs1"),
					GameID:     model.GameIDSolitaire,
					Status:     model.GameStatusFinished,
					Result:     model.GameResultWin,
					Wager:      100,
					Payout:     200,
					StartedAt:  now.Add(-time.Hour),
					FinishedAt: now,
				},
			},
			want: &resource.GameSession{
				SessionId:  faker.UUIDv5("gs1").String(),
				GameId:     resource.GameID_GAME_ID_SOLITAIRE,
				Status:     resource.GameStatus_GAME_STATUS_FINISHED,
				Result:     resource.GameResult_GAME_RESULT_WIN,
				Wager:      100,
				Payout:     200,
				StartedAt:  timestamppb.New(now.Add(-time.Hour)),
				FinishedAt: timestamppb.New(now),
			},
		},
		{
			name: "FinishedAt がゼロ値の場合 => finishedAt は nil",
			args: args{
				sess: model.GameSession{
					ID:         faker.UUIDv5("gs1"),
					FinishedAt: time.Time{},
				},
			},
			want: &resource.GameSession{
				SessionId:  faker.UUIDv5("gs1").String(),
				StartedAt:  timestamppb.New(time.Time{}),
				FinishedAt: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToGameSessionPb(tt.args.sess)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestToGameIDPb(t *testing.T) {
	type args struct {
		id model.GameID
	}
	tests := []struct {
		name string
		args args
		want resource.GameID
	}{
		{
			args: args{
				id: model.GameIDSolitaire,
			},
			want: resource.GameID_GAME_ID_SOLITAIRE,
		},
		{
			args: args{
				id: model.GameIDBlackjack,
			},
			want: resource.GameID_GAME_ID_BLACKJACK,
		},
		{
			args: args{
				id: model.GameIDJanken,
			},
			want: resource.GameID_GAME_ID_JANKEN,
		},
		{
			name: "ゼロ値の場合 => GAME_ID_UNSPECIFIED",
			args: args{
				id: model.GameID(0),
			},
			want: resource.GameID_GAME_ID_UNSPECIFIED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ToGameIDPb(tt.args.id))
		})
	}
}

func TestToGameStatusPb(t *testing.T) {
	type args struct {
		status model.GameStatus
	}
	tests := []struct {
		name string
		args args
		want resource.GameStatus
	}{
		{
			args: args{
				status: model.GameStatusPlaying,
			},
			want: resource.GameStatus_GAME_STATUS_PLAYING,
		},
		{
			args: args{
				status: model.GameStatusFinished,
			},
			want: resource.GameStatus_GAME_STATUS_FINISHED,
		},
		{
			name: "ゼロ値の場合 => GAME_STATUS_UNSPECIFIED",
			args: args{
				status: model.GameStatus(0),
			},
			want: resource.GameStatus_GAME_STATUS_UNSPECIFIED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ToGameStatusPb(tt.args.status))
		})
	}
}

func TestToGameResultPb(t *testing.T) {
	type args struct {
		result model.GameResult
	}
	tests := []struct {
		name string
		args args
		want resource.GameResult
	}{
		{
			args: args{
				result: model.GameResultWin,
			},
			want: resource.GameResult_GAME_RESULT_WIN,
		},
		{
			args: args{
				result: model.GameResultLose,
			},
			want: resource.GameResult_GAME_RESULT_LOSE,
		},
		{
			args: args{
				result: model.GameResultDraw,
			},
			want: resource.GameResult_GAME_RESULT_DRAW,
		},
		{
			name: "ゼロ値の場合 => GAME_RESULT_UNSPECIFIED",
			args: args{
				result: model.GameResult(0),
			},
			want: resource.GameResult_GAME_RESULT_UNSPECIFIED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ToGameResultPb(tt.args.result))
		})
	}
}

func TestToGameID(t *testing.T) {
	type args struct {
		pb resource.GameID
	}
	tests := []struct {
		name    string
		args    args
		want    model.GameID
		wantErr assert.ErrorAssertionFunc
	}{
		{
			args: args{
				pb: resource.GameID_GAME_ID_SOLITAIRE,
			},
			want:    model.GameIDSolitaire,
			wantErr: assert.NoError,
		},
		{
			args: args{
				pb: resource.GameID_GAME_ID_BLACKJACK,
			},
			want:    model.GameIDBlackjack,
			wantErr: assert.NoError,
		},
		{
			args: args{
				pb: resource.GameID_GAME_ID_JANKEN,
			},
			want:    model.GameIDJanken,
			wantErr: assert.NoError,
		},
		{
			name: "未指定の場合 => エラー",
			args: args{
				pb: resource.GameID_GAME_ID_UNSPECIFIED,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrGameIDNotExists)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToGameID(tt.args.pb)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
