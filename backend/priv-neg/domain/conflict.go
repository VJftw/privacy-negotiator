package domain

type DBConflict struct {
	Photo   DBPhoto `gorm:"ForeignKey:PhotoID"`
	PhotoID string
	Targets []DBUser `gorm:"many2many:conflict_targets"`
	Parties []DBUser `gorm:"many2many:conflict_parties"`
}

func (c DBConflict) TableName() string {
	return "conflicts"
}

type CacheConflict struct {
	Targets []string `json:"targets"`
	Parties []string `json:"parties"`
}
