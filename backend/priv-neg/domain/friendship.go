package domain

// CacheFriendship - A representation of a cached friend.
type CacheFriendship struct {
	ID string `json:"id"`
	//Clique string `json:"clique"`
	//TieStrength int `json:"tieStrength"`
}

// WebFriendship - A representation of a friendship communicated with a web client.
type WebFriendship struct {
	From string `json:"from"`
	To   string `json:"to"`
	//Clique string `json:"clique"`
	//TieStrength int `json:"tieStrength"`
}

// CacheFriendshipFromWebFriendship - Converts a WebFriendship into a CacheFriendship
func CacheFriendshipFromWebFriendship(w *WebFriendship) *CacheFriendship {
	return &CacheFriendship{
		ID: w.To,
	}
}

// WebFriendshipFromCacheFriendshipAndCacheUser - Converts a CacheFriendship and authenticated user into a WebFriendship
func WebFriendshipFromCacheFriendshipAndCacheUser(c *CacheFriendship, u *CacheUser) *WebFriendship {
	return &WebFriendship{
		From: u.ID,
		To:   c.ID,
	}
}
