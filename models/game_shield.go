package models

// GameShield ...
type GameShield struct {
	ID int64 `gorm:"primary_key,id" json:"id,omitempty"`

	Strength int `gorm:"column:strength" json:"strength"`
	Mana     int `gorm:"column:mana" json:"mana"`

	//  Five elements
	Metal int `gorm:"column:metal" json:"metal"`
	Wood  int `gorm:"column:wood" json:"wood"`
	Water int `gorm:"column:water" json:"water"`
	Fire  int `gorm:"column:fire" json:"fire"`
	Earth int `gorm:"column:earth" json:"earth"`
}

// CreateGameShield ...
func CreateGameShield(ID int64, Strength int, Mana int, Metal int, Wood int, Water int, Fire int, Earth int) {
	gameShield := GameShield{
		ID:       ID,
		Strength: Strength,
		Mana:     Mana,
		Metal:    Metal,
		Wood:     Wood,
		Water:    Water,
		Fire:     Fire,
		Earth:    Earth,
	}
	err := dbOrmDefault.Model(&GameShield{}).Save(gameShield).Error
	if err != nil {
		MConfig.MLogger.Error(err.Error())
	}
}

// QueryGameShield ...
func QueryGameShield(ID int64) (*GameShield, error) {
	gameShield := &GameShield{
		ID: ID,
	}
	err := dbOrmDefault.Model(&GameShield{}).Find(gameShield).Error
	if err != nil {
		return nil, err
	}
	return gameShield, nil
}

// UpdateGameShield ...
func UpdateGameShield(gameShield *GameShield) error {
	err := dbOrmDefault.Model(&GameShield{}).Update(gameShield).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}
