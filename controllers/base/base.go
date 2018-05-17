// Copyright 2013 Ardan Studios. All rights reserved.
// Use of baseController source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

// Package baseController implements boilerplate code for all baseControllers.
package base

import (
	"runtime"
	"fmt"
	"log"

	"tunaiku/services"
	"tunaiku/utilities/mongo"
	"github.com/astaxie/beego"
)

//** TYPES

type (
	// BaseController composes all required types and behavior.
	BaseController struct {
		beego.Controller
		services.Service
	}
)

//** INTERCEPT FUNCTIONS

// Prepare is called prior to the baseController method.
func (baseController *BaseController) Prepare() {
	beego.ReadFromRequest(&baseController.Controller)

	baseController.UserID = baseController.GetString("userID")

	if baseController.UserID == "" {
		baseController.UserID = baseController.GetString(":userID")
	}

	if baseController.UserID == "" {
		baseController.UserID = "Unknown"
	}

	if err := baseController.Service.Prepare(); err != nil {
		log.Println("[Error] BaseController.PrepareService : ", err)
		baseController.ServeError(err)
		return
	}
}

// Finish is called once the baseController method completes.
func (baseController *BaseController) Finish() {
	defer func() {
		if baseController.MongoSession != nil {
			mongo.CloseSession(baseController.UserID, baseController.MongoSession)
			baseController.MongoSession = nil
		}
	}()
}

//** EXCEPTIONS

// ServeError prepares and serves an Error exception.
func (baseController *BaseController) ServeError(err error) {
	baseController.Data["json"] = struct {
		Error string `json:"Error"`
	}{err.Error()}
	baseController.Ctx.Output.SetStatus(500)
	baseController.ServeJSON()
}

// ServeValidationErrors prepares and serves a validation exception.
func (baseController *BaseController) ServeValidationErrors(Errors []string) {
	baseController.Data["json"] = struct {
		Errors []string `json:"Errors"`
	}{Errors}
	baseController.Ctx.Output.SetStatus(409)
	baseController.ServeJSON()
}

//** CATCHING PANICS

// CatchPanic is used to catch any Panic and log exceptions. Returns a 500 as the response.
func (baseController *BaseController) CatchPanic(functionName string) {
	if r := recover(); r != nil {
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)
		baseController.ServeError(fmt.Errorf("%v", r))
	}
}

//** AJAX SUPPORT

// AjaxResponse returns a standard ajax response.
func (baseController *BaseController) AjaxResponse(resultCode int, resultString string, data interface{}) {
	response := struct {
		Result       int
		ResultString string
		ResultObject interface{}
	}{
		Result:       resultCode,
		ResultString: resultString,
		ResultObject: data,
	}

	baseController.Data["json"] = response
	baseController.ServeJSON()
}
