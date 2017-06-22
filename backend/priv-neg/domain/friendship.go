package domain

// CacheFriendship - A representation of a cached friend. Stored as `<userID>:friends`: friendUserID: {}
type CacheFriendship struct {
	ID          string `json:"id"` // the target user
	TieStrength int    `json:"tieStrength"`
}

// QueueFriendship - Information stored in a queue to be processed.
type QueueFriendship struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// WebFriendship - A representation of a friendship communicated with a web client.
type WebFriendship struct {
	ID          string   `json:"id"`      // the target user
	Cliques     []string `json:"cliques"` // cliques that the session user and target user have in common
	TieStrength int      `json:"tieStrength"`
}

// CacheFriendshipFromWebFriendship - Converts a WebFriendship into a CacheFriendship
func CacheFriendshipFromWebFriendship(w *WebFriendship) *CacheFriendship {
	return &CacheFriendship{
		ID: w.ID,
	}
}

// WebFriendshipFromCacheFriendshipAndCliques - Converts a CacheFriendship and authenticated user into a WebFriendship
func WebFriendshipFromCacheFriendshipAndCliques(c *CacheFriendship, cliques []string) *WebFriendship {
	return &WebFriendship{
		ID:      c.ID,
		Cliques: cliques,
	}
}
