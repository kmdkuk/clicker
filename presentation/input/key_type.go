package input

// KeyType represents the type of key input
type KeyType int

const (
	KeyTypeUp       KeyType = iota // Up
	KeyTypeDown                    // Down
	KeyTypeLeft                    // Left
	KeyTypeRight                   // Right
	KeyTypeDecision                // Decision
	KeyTypeNone                    // No input or other keys
)
