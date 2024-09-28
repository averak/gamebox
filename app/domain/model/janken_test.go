package model

import (
	"testing"
	"time"

	"github.com/averak/gamebox/testutils/faker"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJankenHand_Battle(t *testing.T) {
	type args struct {
		opponent JankenHand
	}
	tests := []struct {
		name string
		mine JankenHand
		args args
		want GameResult
	}{
		{
			name: "[勝利] グー vs チョキ",
			mine: JankenHandRock,
			args: args{opponent: JankenHandScissors},
			want: GameResultWin,
		},
		{
			name: "[勝利] チョキ vs パー",
			mine: JankenHandScissors,
			args: args{opponent: JankenHandPaper},
			want: GameResultWin,
		},
		{
			name: "[勝利] パー vs グー",
			mine: JankenHandPaper,
			args: args{opponent: JankenHandRock},
			want: GameResultWin,
		},
		{
			name: "[引き分け] グー vs グー",
			mine: JankenHandRock,
			args: args{opponent: JankenHandRock},
			want: GameResultDraw,
		},
		{
			name: "[引き分け] チョキ vs チョキ",
			mine: JankenHandScissors,
			args: args{opponent: JankenHandScissors},
			want: GameResultDraw,
		},
		{
			name: "[引き分け] パー vs パー",
			mine: JankenHandPaper,
			args: args{opponent: JankenHandPaper},
			want: GameResultDraw,
		},
		{
			name: "[敗北] グー vs パー",
			mine: JankenHandRock,
			args: args{opponent: JankenHandPaper},
			want: GameResultLose,
		},
		{
			name: "[敗北] チョキ vs グー",
			mine: JankenHandScissors,
			args: args{opponent: JankenHandRock},
			want: GameResultLose,
		},
		{
			name: "[敗北] パー vs チョキ",
			mine: JankenHandPaper,
			args: args{opponent: JankenHandScissors},
			want: GameResultLose,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.mine.Battle(tt.args.opponent))
		})
	}
}

func TestJankenSession_Choose(t *testing.T) {
	now := time.Now()

	type fields struct {
		GameSessionID uuid.UUID
		Seed          int
		Histories     []JankenHistory
	}
	type args struct {
		session GameSession
		hand    JankenHand
		now     time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    GameSession
		want1   Coins
		want2   JankenHistory
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "勝利 => ゲーム終了 & 賞金を獲得",
			fields: fields(JankenNextRockTemplate),
			args: args{
				session: GameSession{
					ID:     JankenNextRockTemplate.GameSessionID,
					Status: GameStatusPlaying,
					Wager:  100,
				},
				hand: JankenHandPaper,
				now:  now,
			},
			want: GameSession{
				ID:         JankenNextRockTemplate.GameSessionID,
				Status:     GameStatusFinished,
				Result:     GameResultWin,
				Wager:      100,
				Payout:     200,
				FinishedAt: now,
			},
			want1: 200,
			want2: JankenHistory{
				Turn:         1,
				MyHand:       JankenHandPaper,
				OpponentHand: JankenHandRock,
			},
			wantErr: assert.NoError,
		},
		{
			name:   "引き分け => ゲームを続行する",
			fields: fields(JankenNextRockTemplate),
			args: args{
				session: GameSession{
					ID:     JankenNextRockTemplate.GameSessionID,
					Status: GameStatusPlaying,
					Wager:  100,
				},
				hand: JankenHandRock,
				now:  now,
			},
			want: GameSession{
				ID:         JankenNextRockTemplate.GameSessionID,
				Status:     GameStatusPlaying,
				Result:     GameResultUnknown,
				Wager:      100,
				Payout:     0,
				FinishedAt: time.Time{},
			},
			want1: 0,
			want2: JankenHistory{
				Turn:         1,
				MyHand:       JankenHandRock,
				OpponentHand: JankenHandRock,
			},
			wantErr: assert.NoError,
		},
		{
			name:   "敗北 => ゲーム終了 & 賞金を獲得しない",
			fields: fields(JankenNextScissorsTemplate),
			args: args{
				session: GameSession{
					ID:     JankenNextScissorsTemplate.GameSessionID,
					Status: GameStatusPlaying,
					Wager:  100,
				},
				hand: JankenHandPaper,
				now:  now,
			},
			want: GameSession{
				ID:         JankenNextScissorsTemplate.GameSessionID,
				Status:     GameStatusFinished,
				Result:     GameResultLose,
				Wager:      100,
				Payout:     0,
				FinishedAt: now,
			},
			want1: 0,
			want2: JankenHistory{
				Turn:         1,
				MyHand:       JankenHandPaper,
				OpponentHand: JankenHandScissors,
			},
			wantErr: assert.NoError,
		},
		{
			name: "ゲームセッションIDが不一致 => エラー",
			fields: fields{
				GameSessionID: faker.UUIDv5("gs1"),
			},
			args: args{
				session: GameSession{
					ID:     faker.UUIDv5("gs2"),
					Status: GameStatusPlaying,
				},
				hand: JankenHandRock,
				now:  now,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, ErrGameNotPlaying, err)
			},
		},
		{
			name: "プレイ中ではない => エラー",
			fields: fields{
				GameSessionID: faker.UUIDv5("gs1"),
			},
			args: args{
				session: GameSession{
					ID:     faker.UUIDv5("gs1"),
					Status: GameStatusFinished,
				},
				hand: JankenHandRock,
				now:  now,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, ErrGameNotPlaying, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JankenSession{
				GameSessionID: tt.fields.GameSessionID,
				Seed:          tt.fields.Seed,
				Histories:     tt.fields.Histories,
			}
			got, got1, got2, err := j.Choose(tt.args.session, tt.args.hand, tt.args.now)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
		})
	}
}

func TestJankenSession_turn(t *testing.T) {
	type fields struct {
		GameSessionID uuid.UUID
		Histories     []JankenHistory
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "履歴が存在しない場合 => 1",
			fields: fields{
				Histories: nil,
			},
			want: 1,
		},
		{
			name: "履歴が存在する場合 => 最新履歴のターン数 + 1",
			fields: fields{
				Histories: []JankenHistory{
					{
						Turn: 1,
					},
					{
						Turn: 2,
					},
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := JankenSession{
				GameSessionID: tt.fields.GameSessionID,
				Histories:     tt.fields.Histories,
			}
			assert.Equal(t, tt.want, j.turn())
		})
	}
}

func TestJankenSession_nextOpponentHand(t *testing.T) {
	type fields struct {
		GameSessionID uuid.UUID
		Seed          int
		Histories     []JankenHistory
	}
	tests := []struct {
		name   string
		fields fields
		want   JankenHand
	}{
		{
			name:   "グーを出す",
			fields: fields(JankenNextRockTemplate),
			want:   JankenHandRock,
		},
		{
			name:   "チョキを出す",
			fields: fields(JankenNextScissorsTemplate),
			want:   JankenHandScissors,
		},
		{
			name:   "パーを出す",
			fields: fields(JankenNextPaperTemplate),
			want:   JankenHandPaper,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := JankenSession{
				GameSessionID: tt.fields.GameSessionID,
				Seed:          tt.fields.Seed,
				Histories:     tt.fields.Histories,
			}

			// 同じターン数であれば同じ手が決定される。
			for range 2 {
				assert.Equal(t, tt.want, j.nextOpponentHand())
			}
		})
	}
}

var (
	// 次のターンに、必ず相手はグーを出すテンプレート。
	JankenNextRockTemplate = JankenSession{
		Seed:      7580503617577668000,
		Histories: []JankenHistory{},
	}
	// 次のターンに、必ず相手はチョキを出すテンプレート。
	JankenNextScissorsTemplate = JankenSession{
		Seed:      5592583862009562302,
		Histories: []JankenHistory{},
	}
	// 次のターンに、必ず相手はパーを出すテンプレート。
	JankenNextPaperTemplate = JankenSession{
		Seed:      4779612329441809388,
		Histories: []JankenHistory{},
	}
)
