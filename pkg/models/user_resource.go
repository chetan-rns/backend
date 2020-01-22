package models

import "log"

// UserResource represents user_resource table in DB
type UserResource struct {
	ResourceID int
	UserID     int
}

// Add a new user resource
func (userResource *UserResource)Add() error {
	res := DB.Create(userResource)
	if res.Error!=nil{
		log.Println(res.Error)
		return res.Error
	}
	return nil
}
