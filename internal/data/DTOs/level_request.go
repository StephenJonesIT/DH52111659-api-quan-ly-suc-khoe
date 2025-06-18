package dtos


type CreateLevelRequest struct {
	ProgramID    string		`json:"program_id" binding:"required"`
	Name         string 	`json:"name" binding:"required"`
	Description  string 	`json:"description"`
	PointRequire int    	`json:"point_require" binding:"required,gt=0"`
}