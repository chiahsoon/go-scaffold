package dals

import (
	"github.com/chiahsoon/go_scaffold/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRefreshTokenDAL struct {}

func (dal *UserRefreshTokenDAL) Upsert(tx *gorm.DB, urt *models.UserRefreshToken) error {
	if createRet := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"refresh_token", "created_at", "updated_at", "deleted_id"}),
	}).Create(urt); createRet.Error != nil {
		return models.NewInternalServerError(createRet.Error.Error())
	}
	return nil
}
