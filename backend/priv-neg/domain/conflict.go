package domain

import "github.com/satori/go.uuid"

// DBConflict - Represents a conflict stored in the database.
type DBConflict struct {
	ID       string  `gorm:"primary_key"`
	Photo    DBPhoto `gorm:"ForeignKey:PhotoID"`
	PhotoID  string
	Target   DBUser `gorm:"ForeignKey:TargetID"`
	TargetID string
	Parties  []DBUser `gorm:"many2many:conflict_parties"`
	Result   string
}

// NewDBConflict - Returns a new DBConflict.
func NewDBConflict() DBConflict {
	return DBConflict{
		ID:      uuid.NewV4().String(),
		Target:  DBUser{},
		Parties: []DBUser{},
		Result:  "indeterminate",
	}
}

// TableName - Returns the table name for the entity.
func (c DBConflict) TableName() string {
	return "conflicts"
}

// CacheConflict - Represents a conflict stored in the Cache.
type CacheConflict struct {
	ID        string   `json:"id"`
	Target    string   `json:"target"`
	Parties   []string `json:"parties"`
	Reasoning []Reason `json:"reasoning"` // reason
	Result    string   `json:"result"`    // result
}

// Reason - Represents a result reason in the cache.
type Reason struct {
	UserID string `json:"id"`
	Vote   int    `json:"vote"`
}

// CacheConflictFromDBConflict - Returns a CacheConflict given a DBConflict
func CacheConflictFromDBConflict(dbConflict DBConflict) CacheConflict {
	cC := CacheConflict{
		ID:        dbConflict.ID,
		Reasoning: []Reason{},
		Target:    dbConflict.TargetID,
	}

	for _, user := range dbConflict.Parties {
		cC.Parties = append(cC.Parties, user.ID)
	}

	return cC
}
