package dtos

type CreateActivityRequest struct {
	LevelID		string		`json:"level_id" binding:"required"`
	Title 		string 		`json:"title" binding:"required"`
	Description string		`json:"description,omitempty"`
	Duration	int			`json:"duration" binding:"required,gt=0"`
	PointReward	int			`json:"point_reward" binding:"required,gt=0"`
	Type        string 		`json:"type" binding:"required,oneof=Activity MiniGame Challenge"`
	RepeatDays	[]string	`json:"repeat_days"`
}
