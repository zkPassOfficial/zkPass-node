package zkpass

import (
	"strconv"
	"time"
)

var (
	version    = "0.0.1"
	commitHash string
	commitTime string

	Version = func() string {
		if commitHash != "" {
			return version + "-" + commitHash
		}
		return version + "-dev"
	}()

	CommitTime = func() string {
		if commitTime == "" {
			commitTime = strconv.Itoa(int(time.Now().Unix()))
		}
		return commitTime
	}
)
