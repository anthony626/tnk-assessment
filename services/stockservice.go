package services

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"

	"tunaiku/models"
	"log"
	"math"
)

func GetAll(service *Service) (*[]models.Stock, error) {
	stock := []models.Stock{}

	f := func(collection *mgo.Collection) error {
		return collection.Find(nil).Sort("date").All(&stock)
	}

	if err := service.DBAction(beego.AppConfig.String("mgo_database"), "stock", f); err != nil {
		log.Println("[Error] StockService.GetAll : ", err)
		if err != mgo.ErrNotFound {
			return nil, err
		}
	}

	for i, item := range stock {
		stock[i].DateStr = item.Date.Format("01/02/2006")
	}

	return &stock, nil
}

func CreateStock(service *Service, stock models.Stock) (*models.Stock, error) {
	f := func(collection *mgo.Collection) error {
		return collection.Insert(&stock)
	}

	if err := service.DBAction(beego.AppConfig.String("mgo_database"), "stock", f); err != nil {
		log.Println("[Error] StockService.CreateStock : ", err)
		return &stock, err
	}

	return &stock, nil
}

func Calculate(service *Service) (*[]models.Stock, string, string, error) {
	stock, err := GetAll(service)
	if err != nil {
		log.Println("[Error] StockService.Calculate GetAll", err)
		return nil, "", "", err
	}

	arr := *stock
	buyDate := arr[0].Date
	sellDate := arr[1].Date
	buyValue := math.MaxInt32
	sellValue := 0

	for i, item := range arr {
		tmp1 := 0
		tmp2 := 0
		tmp1 = item.High - item.Close
		tmp2 = item.Close - item.Low
		if tmp1 < 0 {
			tmp1 = tmp1 * -1
		}
		if tmp2 < 0 {
			tmp2 = tmp2 * -1
		}

		if tmp1 < tmp2 {
			arr[i].Action = "sell"
			if sellValue < tmp1 {
				sellValue = tmp1
				if arr[i].Date.After(buyDate) {
					sellDate = arr[i].Date
				}
			}
		} else {
			arr[i].Action = "buy"
			if buyValue > tmp2 {
				buyValue = tmp2
				buyDate = arr[i].Date
			}
		}
	}

	return &arr, buyDate.Format("01/02/2006"), sellDate.Format("01/02/2006"), err

}