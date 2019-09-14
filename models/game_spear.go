package models

// GameSpear ...
type GameSpear struct {
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

func initSystemGameSpear() {
	CreateGameSpear(15625045984, 1000, 100, 0, 0, 0, 0, 0)
	CreateGameSpear(15622222222, 1000, 100, 0, 0, 0, 0, 0)
	CreateGameSpear(15666666666, 1000, 100, 0, 0, 0, 0, 0)
	CreateGameSpear(15688888888, 1000, 100, 0, 0, 0, 0, 0)
}

// CreateGameSpear ...
func CreateGameSpear(ID int64, Strength int, Mana int, Metal int, Wood int, Water int, Fire int, Earth int) {
	gameSpear := GameSpear{
		ID:       ID,
		Strength: Strength,
		Mana:     Mana,
		Metal:    Metal,
		Wood:     Wood,
		Water:    Water,
		Fire:     Fire,
		Earth:    Earth,
	}
	err := dbOrmDefault.Model(&GameSpear{}).Save(gameSpear).Error
	if err != nil {
		MConfig.MLogger.Error(err.Error())
	}
}

// QueryGameSpear ...
func QueryGameSpear(ID int64) (*GameSpear, error) {
	gameSpear := &GameSpear{
		ID: ID,
	}
	err := dbOrmDefault.Model(&GameSpear{}).Find(gameSpear).Error
	if err != nil {
		return nil, err
	}
	return gameSpear, nil
}

// UpdateGameSpear ...
func UpdateGameSpear(gameSpear *GameSpear) error {
	err := dbOrmDefault.Model(&GameSpear{}).Update(gameSpear).Error
	return err
}
