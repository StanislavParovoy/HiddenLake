package db

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func InClients(user *models.User, hashname string) bool {
	id := GetUserId(user.Auth.Hashpasw)
	if id < 0 {
		return false
	}
	var (
		public string
		err    error
	)
	row := settings.DB.QueryRow(
		"SELECT Public FROM User WHERE IdUser=$1 AND Hashname=$2",
		id,
		hashname,
	)
	err = row.Scan(&public)
	if err != nil {
		return false
	}
	return public != ""
}
