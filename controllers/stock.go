package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"log"
	"path"
	"io/ioutil"
	"strings"

	bc "tunaiku/controllers/base"
	"tunaiku/models"
	"tunaiku/utilities/tools"
	service "tunaiku/services"

)

type StockController struct {
	bc.BaseController
}

func (this *StockController) Index() {
	beego.SetStaticPath("/uploads", "tmp/uploads")
	this.Layout = "layouts/default.html"
	this.TplName = "stock/index.html"
	stock, err := service.GetAll(&this.Service)

	if err != nil {
		log.Println("[Error] StockController.Index : ", err)
	}
	this.Data["Stock"] = stock
}

func (this *StockController) New() {
	this.TplName = "stock/form.html"
	this.Layout = "layouts/default.html"
	this.Data["isValid"] = true
}

func (this *StockController) Create() {
	stock := models.Stock{}
	errorMap := []string{}

	if err := this.ParseForm(&stock); err != nil {
		log.Println("[Error] StockController.Create ParseForm : ", err)

		this.Data["HasErrors"] = true
		this.Data["Errors"] = append(errorMap, "Invalid Date! Use dd/mm/yyyy format")
		this.Data["Stock"] = stock

		this.Layout = "layouts/default.html"
		this.TplName = "stock/form.html"
		return
	}

	valid := validation.Validation{}

	valid.Required(stock.Date, "Date").Message("is required")
	valid.Required(stock.Low, "Low").Message("is required")
	valid.Required(stock.High, "High").Message("is required")
	valid.Required(stock.Open, "Open").Message("is required")
	valid.Required(stock.Close, "Close").Message("is required")
	this.Data["HasErrors"] = false

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			errorMap = append(errorMap, err.Key+" "+err.Message)
		}
		this.Data["HasErrors"] = true
		this.Data["Errors"] = errorMap
		this.Data["Stock"] = stock

		this.Layout = "layouts/default.html"
		this.TplName = "stock/form.html"
		return
	}

	newStock, err := service.CreateStock(&this.Service, stock)
	if err != nil {
		log.Println("[Error] StockController.Create : ", err)

		flash := beego.NewFlash()
		flash.Notice("Something went wrong! Please contact administrator for fix!")
		flash.Store(&this.Controller)

		this.Redirect("/stock", 302)
		return
	}

	flash := beego.NewFlash()
	flash.Notice("Success add stock" + newStock.Date.Format("02/01/2006"))
	flash.Store(&this.Controller)

	this.Redirect("/stock", 302)

}

func (this *StockController) Calculate() {
	this.Layout = "layouts/default.html"
	this.TplName = "stock/calculate.html"
	stock, buyDate, sellDate, err := service.Calculate(&this.Service)
	if err != nil {
		log.Println("[Error] StockController.Index : ", err)
	}
	this.Data["Stock"] = stock
	this.Data["BuyDate"] = buyDate
	this.Data["SellDate"] = sellDate

}

func (this *StockController) Export() {
	header := []string{"ID", "Date", "Open", "High", "Low", "Close"}

	stock, err := service.GetAll(&this.Service)
	if err != nil {
		log.Println("[Error] StockController.Export : ", err)
	}
	filename := tools.ExportStock(header, *stock, "", ".csv")

	basepath := path.Base("/csv-data")
	filepath := path.Base(filename)

	fileBytes, _ := ioutil.ReadFile(strings.Join([]string{basepath, filepath}, "/"))

	this.Ctx.ResponseWriter.Header().Set("Content-Type", "text/csv")
	this.Ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment;filename="+filename)
	this.Ctx.ResponseWriter.Write(fileBytes)
}