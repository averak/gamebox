package main

import "github.com/averak/gamebox/entrypoint/batch_job/job"

var registry = map[string]job.BatchJob{
	"purge_old_echos":         job.NewPurgeOldEchos(),
	"purge_old_game_sessions": job.NewPurgeOldGameSessions(),
}
