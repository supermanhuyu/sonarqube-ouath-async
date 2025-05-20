package sonar

import (
	_ "github.com/pkg/errors"
	orm "sonarqube-ouath-async/global"
)

type Project struct {
	Kee         string `json:"kee" gorm:"kee"`                 // Identifier
	Qualifier   string `json:"qualifier" gorm:"qualifier"`     //
	Name        string `json:"name" gorm:"name"`               //
	Description string `json:"description" gorm:"description"` //
	Private     bool   `json:"private" gorm:"private"`         //
	Tags        string `json:"tags" gorm:"tags"`               //
	Ncloc       int    `json:"ncloc" gorm:"ncloc"`

	//Uuid        string `json:"uuid" gorm:"uuid"`               //
	//CreatedAt   int    `json:"created_at" gorm:"created_at"`   //
	//UpdatedAt   int    `json:"updated_at" gorm:"updated_at"`   //
}

func (*Project) TableName() string {
	return "projects"
}

func (e *Project) Get() (Project, error) {
	var doc Project
	table := orm.SonarDb.Table(e.TableName())

	if e.Kee != "" {
		table = table.Where("kee = ?", e.Kee)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *Project) GetPage(pageSize int, pageIndex int) ([]Project, int64, error) {
	var count int64
	orm.SonarDb.Select("*").Table(e.TableName()).Count(&count)

	var docs []Project
	table := orm.SonarDb.Select("*").Table(e.TableName())
	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&docs).Error; err != nil {
		return nil, 0, err
	}
	return docs, count, nil
}

func (e *Project) List() ([]Project, error) {
	var docs []Project
	table := orm.SonarDb.Select("*").Table(e.TableName())
	if err := table.Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}
