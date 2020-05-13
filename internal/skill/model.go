package skill

import (
	"bitbucket.movista.ru/maas/maasapi/internal/models"
)

type Skill struct {
	tableName         struct{}    `pg:"maasapi.skill"`
	ID                string      `pg:",pk" json:"id"`
	Title             string      `pg:"title,unique,notnull" json:"description"`
	Order             int         `pg:"ord,unique,notnull" json:"order"`
	ActionID          string      `pg:"action_id" json:"action_id"`
	Icon              string      `pg:"icon" json:"icon"`
	IosIconURL        string      `pg:"ios_icon_url" json:"ios_icon_url"`
	AndroidIconURL    string      `pg:"android_icon_url" json:"android_icon_url"`
	DisableInVersions string      `pg:"disable_in_versions" json:"-" example:">= 1.0, < 1.4"`
	CreatedAt         models.Time `json:"-"`
}

func (s Skill) ConvertToDataObject() *models.DataObject {
	return &models.DataObject{
		ObjectId: ObjectID,
		Data:     s,
	}
}
