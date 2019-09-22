package models

// GameOrder 主体对客体做行为
type GameOrder struct {
	// 主体
	FromGroup int   `json:"fromGroup"`
	FromID    int64 `json:"fromId"`

	// 客体
	ToGroup int   `json:"toGroup"`
	ToID    int64 `json:"toId"`

	// 分类行为
	Type int `json:"type"`
	ID   int `json:"id"`

	// 时间
	TimeCreate  int64 `json:"timeCreate"`
	TimeCurrent int64 `json:"timeCurrent"`

	// 附加信息
	Data interface{} `json:"data"`
}
