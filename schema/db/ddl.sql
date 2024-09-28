-----------------------------------
-- ユーザ
-----------------------------------

CREATE TABLE "users"
(
    "id"         UUID      NOT NULL,
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
