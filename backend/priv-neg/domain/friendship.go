package domain

// CacheFriend - A representation of a cached friend.
type CacheFriendship struct {
	ID string `json:"id"`
	//Clique string `json:"clique"`
	//TieStrength int `json:"tieStrength"`
}

type WebFriendship struct {
	From string `json:"from"`
	To   string `json:"to"`
	//Clique string `json:"clique"`
	//TieStrength int `json:"tieStrength"`
}

func CacheFriendshipFromWebFriendship(w *WebFriendship) *CacheFriendship {
	return &CacheFriendship{
		ID: w.To,
	}
}

func WebFriendshipFromCacheFriendshipAndCacheUser(c *CacheFriendship, u *CacheUser) *WebFriendship {
	return &WebFriendship{
		From: u.ID,
		To:   c.ID,
	}
}
