// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mongo provides mongo connectivity support.
package mongo

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
	"github.com/davecgh/go-spew/spew"
)

const (
	// MasterSession provides direct access to master database.
	MasterSession = "master"

	// MonotonicSession provides reads to slaves.
	MonotonicSession = "monotonic"
)

var (
	// Reference to the singleton.
	singleton mongoManager
)

type (
	// mongoManager contains dial and session information.
	mongoSession struct {
		mongoDBDialInfo *mgo.DialInfo
		mongoSession    *mgo.Session
	}

	// mongoManager manages a map of session.
	mongoManager struct {
		sessions map[string]mongoSession
	}

	// DBCall defines a type of function that can be used
	// to excecute code against MongoDB.
	DBCall       func(*mgo.Collection) error
	DBCallUpsert func(*mgo.Collection) (*mgo.ChangeInfo, error)
	DBCallPipe   func(*mgo.Collection) (*mgo.Pipe, error)
)

// Startup brings the manager to a running state.
func Startup(sessionID string) error {
	// If the system has already been started ignore the call.
	if singleton.sessions != nil {
		return nil
	}

	// Create the Mongo Manager.
	singleton = mongoManager{
		sessions: make(map[string]mongoSession),
	}
	
	hosts := strings.Split(beego.AppConfig.String("mgo_host"), ",")
	// Create the strong session.
	if err := CreateSession(sessionID, "strong", MasterSession, hosts, "admin", beego.AppConfig.String("mgo_username"), beego.AppConfig.String("mgo_password")); err != nil {
		fmt.Println("[Error] Startup.CreateSession Strong: ", err)
		return err
	}

	// Create the monotonic session.
	if err := CreateSession(sessionID, "monotonic", MonotonicSession, hosts, "admin", beego.AppConfig.String("mgo_username"), beego.AppConfig.String("mgo_password")); err != nil {
		fmt.Println("[Error] Startup.CreateSession Monotonic: ", err)
		return err
	}

	return nil
}

// Shutdown systematically brings the manager down gracefully.
func Shutdown(sessionID string) error {
	// Close the databases
	for _, session := range singleton.sessions {
		CloseSession(sessionID, session.mongoSession)
	}
	return nil
}

// CreateSession creates a connection pool for use.
func CreateSession(sessionID string, mode string, sessionName string, hosts []string, databaseName string, username string, password string) error {
	// Create the database object
	mongoSession := mongoSession{
		mongoDBDialInfo: &mgo.DialInfo{
			Addrs:    hosts,
			Timeout:  60 * time.Second,
			Database: databaseName,
			Username: username,
			Password: password,
		},
	}
	spew.Dump(mongoSession)

	// Establish the master session.
	var err error
	mongoSession.mongoSession, err = mgo.DialWithInfo(mongoSession.mongoDBDialInfo)
	if err != nil {
		fmt.Println("[Error] MongoUtility.Startup DialWithInfo: ", err)
		return err
	}
	switch mode {
	case "strong":
		// Reads and writes will always be made to the master server using a
		// unique connection so that reads and writes are fully consistent,
		// ordered, and observing the most up-to-date data.
		// http://godoc.org/github.com/finapps/mgo#Session.SetMode
		mongoSession.mongoSession.SetMode(mgo.Strong, true)
		break

	case "monotonic":
		// Reads may not be entirely up-to-date, but they will always see the
		// history of changes moving forward, the data read will be consistent
		// across sequential queries in the same session, and modifications made
		// within the session will be observed in following queries (read-your-writes).
		// http://godoc.org/github.com/finapps/mgo#Session.SetMode
		mongoSession.mongoSession.SetMode(mgo.Monotonic, true)
	}

	// Have the session check for errors.
	// http://godoc.org/github.com/finapps/mgo#Session.SetSafe
	mongoSession.mongoSession.SetSafe(&mgo.Safe{})

	// Add the database to the map.
	singleton.sessions[sessionName] = mongoSession

	return nil
}

// CopyMasterSession makes a copy of the master session for client use.
func CopyMasterSession(sessionID string) (*mgo.Session, error) {
	return CopySession(sessionID, MasterSession)
}

// CopyMonotonicSession makes a copy of the monotonic session for client use.
func CopyMonotonicSession(sessionID string) (*mgo.Session, error) {
	return CopySession(sessionID, MonotonicSession)
}

// CopySession makes a copy of the specified session for client use.
func CopySession(sessionID string, useSession string) (*mgo.Session, error) {
	// Find the session object.
	session := singleton.sessions[useSession]
	if session.mongoSession == nil {
		err := fmt.Errorf("Unable To Locate Session %s", useSession)
		return nil, err
	}

	// Copy the master session.
	mongoSession := session.mongoSession.Copy()

	return mongoSession, nil
}

// CloneMasterSession makes a clone of the master session for client use.
func CloneMasterSession(sessionID string) (*mgo.Session, error) {
	return CloneSession(sessionID, MasterSession)
}

// CloneMonotonicSession makes a clone of the monotinic session for client use.
func CloneMonotonicSession(sessionID string) (*mgo.Session, error) {
	return CloneSession(sessionID, MonotonicSession)
}

// CloneSession makes a clone of the specified session for client use.
func CloneSession(sessionID string, useSession string) (*mgo.Session, error) {
	// Find the session object.
	session := singleton.sessions[useSession]

	if session.mongoSession == nil {
		err := fmt.Errorf("Unable To Locate Session %s", useSession)
		return nil, err
	}

	// Clone the master session.
	mongoSession := session.mongoSession.Clone()

	return mongoSession, nil
}

// CloseSession puts the connection back into the pool.
func CloseSession(sessionID string, mongoSession *mgo.Session) {
	mongoSession.Close()
}

// GetDatabase returns a reference to the specified database.
func GetDatabase(mongoSession *mgo.Session, useDatabase string) *mgo.Database {
	return mongoSession.DB(useDatabase)
}

// GetCollection returns a reference to a collection for the specified database and collection name.
func GetCollection(mongoSession *mgo.Session, useDatabase string, useCollection string) *mgo.Collection {
	return mongoSession.DB(useDatabase).C(useCollection)
}

// CollectionExists returns true if the collection name exists in the specified database.
func CollectionExists(sessionID string, mongoSession *mgo.Session, useDatabase string, useCollection string) bool {
	database := mongoSession.DB(useDatabase)
	collections, err := database.CollectionNames()

	if err != nil {
		return false
	}

	for _, collection := range collections {
		if collection == useCollection {
			return true
		}
	}

	return false
}

// ToString converts the quer map to a string.
func ToString(queryMap interface{}) string {
	json, err := json.Marshal(queryMap)
	if err != nil {
		return ""
	}

	return string(json)
}

// ToStringD converts bson.D to a string.
func ToStringD(queryMap bson.D) string {
	json, err := json.Marshal(queryMap)
	if err != nil {
		return ""
	}

	return string(json)
}

// Execute the MongoDB literal function.
func Execute(sessionID string, mongoSession *mgo.Session, databaseName string, collectionName string, dbCall DBCall) error {
	// Capture the specified collection.
	collection := GetCollection(mongoSession, databaseName, collectionName)
	if collection == nil {
		err := fmt.Errorf("Collection %s does not exist", collectionName)
		return err
	}

	// Execute the MongoDB call.
	err := dbCall(collection)
	if err != nil {
		return err
	}
	return nil
}

// Execute the MongoDB literal function.
func ExecuteUpsert(sessionID string, mongoSession *mgo.Session, databaseName string, collectionName string, dbCall DBCallUpsert) error {
	// Capture the specified collection.
	collection := GetCollection(mongoSession, databaseName, collectionName)
	if collection == nil {
		err := fmt.Errorf("Collection %s does not exist", collectionName)
		return err
	}

	// Execute the MongoDB call.
	_, err := dbCall(collection)
	if err != nil {
		return err
	}
	return nil
}

// Execute the MongoDB literal function.
func ExecutePipe(sessionID string, mongoSession *mgo.Session, databaseName string, collectionName string, dbCall DBCallPipe) error {
	// Capture the specified collection.
	collection := GetCollection(mongoSession, databaseName, collectionName)
	if collection == nil {
		err := fmt.Errorf("Collection %s does not exist", collectionName)
		return err
	}

	// Execute the MongoDB call.
	_, err := dbCall(collection)
	if err != nil {
		return err
	}
	return nil
}
