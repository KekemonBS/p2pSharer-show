package models

import "crypto/sha256"

// Folder related info
type Folder struct {
	Name     string
	Size     int
	FileTree string
	CID      string
}

// User related info
type User struct {
	Name     string
	PassHash [sha256.Size]byte
}
