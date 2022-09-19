package repository

import (
	"fmt"
	"gorm.io/gorm"
	"kouhei-github/sample-gin/service"
	"unicode/utf8"
)

type MerchandiseEntity struct {
	gorm.Model
	Image            string `json:"image" binding:"required"`
	Name             string `json:"name" binding:"required"`
	Detail           string `json:"detail" binding:"required"`
	Status           int    `json:"status" binding:"required"`
	Carriage         int    `json:"carriage" binding:"required"`
	RequestRequired  int    `json:"request_required" binding:"required"`
	SellPrice        int    `json:"sell_price" binding:"required"`
	IsUpload         string `json:"is_upload"`
	IsPurchased      string `json:"is_purchased"`
	DeliveryEntityID uint   `json:"delivery_id" binding:"required"`
	DeliveryEntity   DeliveryEntity
}

func Upload(entity MerchandiseEntity, myError error) {
	fmt.Println(myError)
}

func NewMerchandiseEntity(
	image string,
	name string,
	detail string,
	status int,
	carriage int,
	requestRequired int,
	sellPrice int,
	deliveryEntityID uint,
) (*MerchandiseEntity, error) {
	if utf8.RuneCountInString(image) <= 1 {
		err := service.MyError{Message: "画像を入力してください"}
		return &MerchandiseEntity{}, err
	}
	if utf8.RuneCountInString(name) <= 1 {
		err := service.MyError{Message: "商品名を入力してください"}
		return &MerchandiseEntity{}, err
	}
	if utf8.RuneCountInString(detail) <= 1 {
		err := service.MyError{Message: "商品の説明を入力してください"}
		return &MerchandiseEntity{}, err
	}
	entity := MerchandiseEntity{
		Image:            image,
		Name:             name,
		Detail:           detail,
		Status:           status,
		Carriage:         carriage,
		RequestRequired:  requestRequired,
		SellPrice:        sellPrice,
		DeliveryEntityID: deliveryEntityID,
	}
	return &entity, nil
}

func CreateMerchandiseList(merchandises []MerchandiseEntity) error {
	result := db.Create(&merchandises)
	if result.Error != nil {
		myErr := service.MyError{
			Message: result.Error.Error(),
		}
		return myErr
	}
	return nil
}

func UpdateMerchandiseList(merchandises []MerchandiseEntity) error {
	result := db.Save(&merchandises)
	if result.Error != nil {
		myErr := service.MyError{
			Message: result.Error.Error(),
		}
		return myErr
	}
	return nil
}

func FindByUploadTarget(search string) ([]MerchandiseEntity, error) {
	var merchandiseEntities []MerchandiseEntity
	var merchandiseEntity MerchandiseEntity
	result := db.Model(merchandiseEntity).Where("is_upload = ?", search).Find(&merchandiseEntities)
	if result.Error != nil {
		myErr := service.MyError{Message: result.Error.Error()}
		return []MerchandiseEntity{}, myErr
	}
	fmt.Println(merchandiseEntities[0])
	return merchandiseEntities, nil
}

type validateResponse struct {
	Result bool `json:"result"`
}
