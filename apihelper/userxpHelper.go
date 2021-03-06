package apihelper

import (
	"errors"
	"log"

	"github.com/classAndrew/govalor/models"
)

// UpdateUserTotalXP modifies the DB table
func UpdateUserTotalXP(input models.UserTotalXP) error {
	var user models.UserTotalXP
	if err := models.DB.Where("name = ?", input.Name).Where("guild = ?", input.Guild).First(&user).Error; err != nil {
		return err
	}
	// user.XP = input.XP
	// models.DB.Save(&user)
	models.DB.Model(&models.UserTotalXP{}).Where("name = ?", input.Name).Where("guild = ?", input.Guild).Update(input)
	// updates is broken? replaces everyone with one username
	// models.DB.Model(&user).Update(input)
	return nil
}

// UpdateUserTotalXPTX modifies the DB table but in a single transaction
func UpdateUserTotalXPTX(members []models.UserTotalXP) {
	tx := models.DB.Begin()
	// stmt := "UPDATE user_total_xps SET last_xp = xp"

	// err := tx.Exec(stmt).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	log.Fatalln(err)
	// }
	for i := 0; i < len(members); i++ {

		stmt := "UPDATE user_total_xps SET xp = ?, last_xp = ? WHERE uuid = ?"

		err := tx.Exec(stmt, members[i].XP, members[i].LastXP, members[i].UUID).Error
		if err != nil {
			tx.Rollback()
			log.Fatalln(err)
		}
	}

	err := tx.Commit().Error
	if err != nil {
		log.Fatalln(err)
	}
}

// FindSpecificUserTotalXP fetches specfic a user i.e. with guild name
func FindSpecificUserTotalXP(guild string, name string, user *models.UserTotalXP) error {
	if err := models.DB.Where("name = ?", name).Where("guild = ?", guild).First(user).Error; err != nil {
		return err
	}
	return nil
}

// GetGuildMembersXP does above but all guild members
func GetGuildMembersXP() []models.UserTotalXP {
	var res []models.UserTotalXP
	models.DB.Table("user_total_xps").Find(&res)
	return res
}

// CreateUserTotalXP creates record. errors if record exists
func CreateUserTotalXP(guild string, name string, xp int64, uuid string) error {
	var user models.UserTotalXP
	err := FindSpecificUserTotalXP(guild, name, &user)

	if err == nil {
		return errors.New("record already exists")
	} else if err.Error() == "record not found" {
		// this is good, no existing users found so create new record
		models.DB.Create(&models.UserTotalXP{Guild: guild, Name: name, XP: xp, LastXP: xp, UUID: uuid})
		return nil
	} else {
		log.Println(err.Error())
		return err
	}
}
