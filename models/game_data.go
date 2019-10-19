package models

// GameData ...
type GameData struct {
	ID int64 `gorm:"primary_key,id" json:"id,omitempty"`

	Name  string `gorm:"column:name;type:varchar(100)" json:"name,omitempty"`
	Level int    `gorm:"column:level" json:"level"`
	Blood int    `gorm:"column:blood" json:"blood"`

	Speed float64 `gorm:"column:speed" json:"speed"`
	X     float64 `gorm:"column:x" json:"x"`
	Z     float64 `gorm:"column:z" json:"z"`

	Spear  *GameSpear  `gorm:"foreignkey:id" json:"spear"`
	Shield *GameShield `gorm:"foreignkey:id" json:"shield"`

	OrderMap map[int]*GameOrder `gorm:"-" json:"orderMap"`
}

func initSystemGameData() {
	CreateDefaultGameData(15625045984, "乂末")
	CreateDefaultGameData(15622222222, "Morge")
	CreateDefaultGameData(15666666666, "Morge")
	CreateDefaultGameData(15688888888, "SghenMorge")
}

func CreateDefaultGameData(ID int64, Name string) {
	CreateGameData(ID, Name, 100, 10000, 50, 0, 0)
}

// CreateGameData ...
func CreateGameData(ID int64, Name string, Level int, Blood int, Speed float64, X float64, Z float64) {
	gameData := GameData{
		ID:    ID,
		Name:  Name,
		Level: Level,
		Blood: Blood,
		Speed: Speed,
		X:     X,
		Z:     Z,
	}
	err := dbOrmDefault.Model(&GameData{}).Save(gameData).Error
	if err != nil {
		MConfig.MLogger.Error(err.Error())
	}

	CreateGameShield(ID, 800, 80, 0, 0, 0, 0, 0)
	CreateGameSpear(ID, 1000, 100, 0, 0, 0, 0, 0)
}

// QueryGameData ...
func QueryGameData(ID int64) (*GameData, error) {
	gameData := &GameData{
		ID: ID,
	}
	err := dbOrmDefault.Model(&GameData{}).Preload("Spear").Preload("Shield").Find(gameData).Error
	if err != nil {
		return nil, err
	}
	return gameData, nil
}

// UpdateGameData ...
func UpdateGameData(gameData *GameData) error {
	// update: nothing will be updated such as "", 0, false are blank values of their types
	err := dbOrmDefault.Model(&GameData{}).Save(gameData).Error
	return err
}
