package api

import "time"

type VersionInfo struct {
	Name    string
	Short   string
	Version string
	Time    time.Time
}
