package user

import (
	"github.com/vzau/thoth/pkg/database"
	dbTypes "github.com/vzau/types/database"
)

func GetUser(cid uint64) (*dbTypes.User, error) {
	user := &dbTypes.User{}
	if err := database.DB.Where("id = ?", cid).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetRoles(cid uint64) ([]*dbTypes.Role, error) {
	user, err := GetUser(cid)
	if err != nil {
		return nil, err
	}

	return user.Roles, nil
}
