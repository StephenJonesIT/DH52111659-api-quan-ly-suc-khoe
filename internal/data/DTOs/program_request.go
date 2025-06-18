package dtos

type CreateProgramRequest struct {
	Title 		string    		`json:"title" binding:"required"`
	Description string 			`json:"desciption"`
	Duration	int				`json:"duration" binding:"required,gt=0"`
	DiseaseIDs  []int			`json:"disease_ids"`	
	GoalIDs		[]int			`json:"goal_ids"`
	Levels		[]CreateLevel 	`json:"levels"`
}

type CreateLevel struct {
	Name  			string 				`json:"name" binding:"required"`
	Description 	string 				`json:"description"`
	PointRequire  	int					`json:"point_require" binding:"required,gt=0"`
	Activities 		[]CreateActivity 	`json:"activities"`
}

type CreateActivity struct {
	Title 	string 			`json:"title" binding:"required"`
	Description string		`json:"description,omitempty"`
	Duration	int			`json:"duration" binding:"required,gt=0"`
	PointReward	int			`json:"point_reward" binding:"required,gt=0"`
	Type        string 		`json:"type" binding:"required,oneof=Activity MiniGame Challenge"`
	RepeatDay	[]string	`json:"repeat_days"`
}

type UpdateProgramRequest struct{
	Title 		string    		`json:"title" binding:"required"`
	Description string 			`json:"desciption,omitempty"`
	Duration	int				`json:"duration" binding:"required,gt=0"`
	DiseaseIDs  []int			`json:"disease_ids"`	
	GoalIDs		[]int			`json:"goal_ids"`
}