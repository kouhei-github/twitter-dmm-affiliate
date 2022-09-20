package repository

import (
	"fmt"
	"gorm.io/gorm"
	"kouhei-github/sample-gin/service"
	"strconv"
)

type AutoFolowingEntity struct {
	gorm.Model
	TwitterUserId string `json:"twitter_user_id"`
	Status        int8   `json:"status"`
	ExpireDate    string
}

func NewAutoFolowingEntity(userId string, status int8, expireAt string) (*AutoFolowingEntity, error) {
	enums := service.GetStatusEnum()
	isStatusNum := service.Contains(enums, status)
	if !isStatusNum {
		myError := service.MyError{Message: "Statusの値が正しくないです。"}
		return &AutoFolowingEntity{}, myError
	}

	carbonDate, err := service.GetCarbonDate(expireAt)

	if err != nil {
		myError := service.MyError{Message: "carbonに変換できませんでした。"}
		return &AutoFolowingEntity{}, myError
	}
	afterTwoWeek := service.AfterWeek(*carbonDate, 1)
	afterYear := afterTwoWeek.Year()
	afterMonth := afterTwoWeek.Month()
	afterDate := afterTwoWeek.Day()
	afterDateTime := strconv.Itoa(afterYear) + "-" + strconv.Itoa(int(afterMonth)) + "-" + strconv.Itoa(afterDate)
	return &AutoFolowingEntity{TwitterUserId: userId, Status: status, ExpireDate: afterDateTime}, nil
}

func (entity *AutoFolowingEntity) Create() error {
	result := db.Create(entity)
	if result.Error != nil {
		myError := service.MyError{Message: "AutoFolowingEntityを保存できませんでした。"}
		return myError
	}
	return nil
}

func BulkInsertAutoFollowing(autoFolow []AutoFolowingEntity) error {
	fmt.Println(autoFolow)
	result := db.Create(&autoFolow)
	if result.Error != nil {
		myErr := service.MyError{
			Message: result.Error.Error(),
		}
		return myErr
	}
	return nil
}

func FindByTwitterUserId(userId string) ([]AutoFolowingEntity, error) {
	var entity []AutoFolowingEntity
	result := db.Where("twitter_user_id = ?", userId).First(&entity)
	if result.Error != nil {
		err := service.MyError{Message: "twitter_user_idで検索できませんでした"}
		return []AutoFolowingEntity{}, err
	}
	return entity, nil
}

func FindByLastRecord() (AutoFolowingEntity, error) {
	var entity AutoFolowingEntity
	result := db.Last(&entity)
	if result.Error != nil {
		myErr := service.MyError{Message: result.Error.Error()}
		return AutoFolowingEntity{}, myErr
	}
	return entity, nil
}

func FindByExpireDate(expireDate string) ([]AutoFolowingEntity, error) {
	var entities []AutoFolowingEntity
	autoEntity := AutoFolowingEntity{ExpireDate: expireDate}
	result := db.Where(&autoEntity).Find(&entities)
	if result.Error != nil {
		myErr := service.MyError{Message: result.Error.Error()}
		return []AutoFolowingEntity{}, myErr
	}
	return entities, nil
}
