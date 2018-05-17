// Copyright 2013 Ardan Studios. All rights reserved.
// Use of service source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

// Package services implements boilerplate code for all services.
package services

import (
	"tunaiku/utilities/helper"
	"tunaiku/utilities/mongo"
	"gopkg.in/mgo.v2"
	"log"
)

//** TYPES
type (
	// Service contains common properties for all services.
	Service struct {
		MongoSession *mgo.Session
		UserID       string
	}
)

//** PUBLIC FUNCTIONS

// Prepare is called before any controller.
func (service *Service) Prepare() (err error) {
	service.MongoSession, err = mongo.CopyMonotonicSession("monotonic")
	if err != nil {
		log.Println("[Error] Service.Prepare : ", err)
		return err
	}

	return err
}

// Finish is called after the controller.
func (service *Service) Finish() (err error) {
	defer helper.CatchPanic(&err, service.UserID, "Service.Finish")

	if service.MongoSession != nil {
		mongo.CloseSession(service.UserID, service.MongoSession)
		service.MongoSession = nil
	}

	return err
}

// DBAction executes the MongoDB literal function
func (service *Service) DBAction(databaseName string, collectionName string, dbCall mongo.DBCall) (err error) {
	return mongo.Execute(service.UserID, service.MongoSession, databaseName, collectionName, dbCall)
}

func (service *Service) DBActionUpsert(databaseName string, collectionName string, dbCall mongo.DBCallUpsert) (err error) {
	return mongo.ExecuteUpsert(service.UserID, service.MongoSession, databaseName, collectionName, dbCall)
}
