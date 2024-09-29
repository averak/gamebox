-----------------------------------
-- ユーザ
-----------------------------------

CREATE TABLE "users"
(
    "id"         UUID      NOT NULL,
    "status"     INT       NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "user_game_sessions"
(
    "id"          UUID      NOT NULL,
    "user_id"     UUID      NOT NULL,
    "game_id"     INT       NOT NULL,
    "status"      INT       NOT NULL,
    "result"      INT       NOT NULL,
    "wager"       INT       NOT NULL,
    "payout"      INT       NOT NULL,
    "started_at"  TIMESTAMP NOT NULL,
    "finished_at" TIMESTAMP NULL,
    "created_at"  TIMESTAMP NOT NULL,
    "updated_at"  TIMESTAMP NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk__user_game_sessions__user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);
CREATE INDEX "idx__user_game_sessions__user_id__status" ON "user_game_sessions" ("user_id", "status");

CREATE TABLE "user_janken_sessions"
(
    "game_session_id" UUID      NOT NULL,
    "seed"            INT       NOT NULL,
    "created_at"      TIMESTAMP NOT NULL,
    "updated_at"      TIMESTAMP NOT NULL,
    PRIMARY KEY ("game_session_id"),
    CONSTRAINT "fk__user_janken_sessions__game_session_id" FOREIGN KEY ("game_session_id") REFERENCES "user_game_sessions" ("id") ON DELETE CASCADE
);

CREATE TABLE "user_janken_session_histories"
(
    "game_session_id" UUID      NOT NULL,
    "turn"            INT       NOT NULL,
    "my_hand"         INT       NOT NULL,
    "opponent_hand"   INT       NOT NULL,
    "created_at"      TIMESTAMP NOT NULL,
    "updated_at"      TIMESTAMP NOT NULL,
    PRIMARY KEY ("game_session_id", "turn"),
    CONSTRAINT "fk__user_janken_session_histories__game_session_id" FOREIGN KEY ("game_session_id") REFERENCES "user_janken_sessions" ("game_session_id") ON DELETE CASCADE
);

-----------------------------------
-- その他
-----------------------------------

CREATE TABLE "echos"
(
    "id"         UUID         NOT NULL,
    "message"    VARCHAR(255) NOT NULL,
    "timestamp"  TIMESTAMP    NOT NULL,
    "created_at" TIMESTAMP    NOT NULL,
    "updated_at" TIMESTAMP    NOT NULL,
    PRIMARY KEY ("id")
);
