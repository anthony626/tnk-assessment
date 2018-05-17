package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type (
	Stock struct {
		ID          bson.ObjectId 	`bson:"_id,omitempty" json:"id"`
		Open        int 	       	`bson:"open" json:"open" form:"open" valid:"Required"`
		Close  		int     	   	`bson:"close" json:"close" form:"close" valid:"Required"`
		High        int        		`bson:"high" json:"high" form:"high"`
		Low     	int        		`bson:"low" json:"low" form:"low"`
		Date    	time.Time     	`bson:"date" json:"-" form:"date,01/02/2006"`
		DateStr		string			`bson:"-" json:"date_str"`
		Action		string 			`bson:"-" json:"action"`
	}
)
