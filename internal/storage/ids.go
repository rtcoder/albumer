package storage

import "albumer/internal/ids"

func newID() string {
	return ids.NewUUID()
}
