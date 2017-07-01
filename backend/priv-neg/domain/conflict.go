package domain

import "github.com/satori/go.uuid"

type DBConflict struct {
	ID       string  `gorm:"primary_key"`
	Photo    DBPhoto `gorm:"ForeignKey:PhotoID"`
	PhotoID  string
	Targets  []DBUser `gorm:"many2many:conflict_targets"`
	Parties  []DBUser `gorm:"many2many:conflict_parties"`
	Resolved bool
}

func NewDBConflict() DBConflict {
	return DBConflict{
		ID:       uuid.NewV4().String(),
		Targets:  []DBUser{},
		Parties:  []DBUser{},
		Resolved: false,
	}
}

func (c DBConflict) TableName() string {
	return "conflicts"
}

type CacheConflict struct {
	ID       string   `json:"id"`
	Targets  []string `json:"targets"`
	Parties  []string `json:"parties"`
	Resolved bool     `json:"resolved"`
}

func CacheConflictFromDBConflict(dbConflict DBConflict) CacheConflict {
	cC := CacheConflict{
		ID:       dbConflict.ID,
		Resolved: dbConflict.Resolved,
	}
	for _, user := range dbConflict.Targets {
		cC.Targets = append(cC.Targets, user.ID)
	}
	for _, user := range dbConflict.Parties {
		cC.Parties = append(cC.Parties, user.ID)
	}

	return cC
}
