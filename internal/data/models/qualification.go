package models

import "time"

type Qualification struct {
	QualificationID int `json:"qualification_id" gorm:"column:qualification_id;primaryKey"`
	QualificationRequest
}

type QualificationRequest struct {
	ExpertID      	int    		`json:"expert_id" gorm:"column:expert_id;not null"`
	DegreeType    	string 		`json:"degree_type" gorm:"column:degree_type;not null"`
	Major         	string 		`josn:"major" gorm:"column:major;not null"`
	Institution   	string 		`json:"institution" gorm:"column:institution;not null"`
	YearGraduated 	*time.Time 	`json:"year_graduated,omitempty" gorm:"column:year_graduated"`
	CertificateURL 	string 		`json:"certificate_url" gorm:"column:certificate_url"`
	Verified 		bool 		`json:"verified" gorm:"column:verified;default:false"`
}

func (Qualification) TableName() string{
	return "qualifications"
}

func(QualificationRequest) TableName() string{
	return Qualification{}.TableName()
}