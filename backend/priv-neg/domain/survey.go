package domain

// DBSurvey - represents a survey stored on the database
type DBSurvey struct {
	UserID  string `gorm:"primary_key"`
	PhotoID string `gorm:"primary_key"`
	RawJSON string
}

// TableName - Returns the table name for the entity.
func (u DBSurvey) TableName() string {
	return "surveys"
}
