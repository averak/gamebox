CREATE TABLE "echos"
(
    "id"         UUID         NOT NULL,
    "message"    VARCHAR(255) NOT NULL,
    "timestamp"  TIMESTAMP    NOT NULL,
    "created_at" TIMESTAMP    NOT NULL,
    "updated_at" TIMESTAMP    NOT NULL,
    PRIMARY KEY ("id")
);
