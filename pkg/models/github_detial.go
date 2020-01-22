package models

import (
	"log"
)

// GithubDetail represents github
type GithubDetail struct {
	ResourceID     int
	Owner          string
	RepositoryName string
	Path           string
	ReadmePath     string `gorm:"default:''"`
}

// Add github details to DB
func (githubDetail *GithubDetail) Add() {
	res := DB.Create(githubDetail)
	if res.Error != nil {
		log.Println(res.Error)
	}
}

// Update github details to DB
func (githubDetail *GithubDetail) Update(parameter string, value interface{}) {
	res := DB.Model(githubDetail).Update(parameter, value)
	if res.Error != nil {
		log.Println(res.Error)
	}
}

// GetResourceGithubDetails will return resource path and github details
func GetResourceGithubDetails(resourceID int) GithubDetail {
	// sqlStatement := `SELECT * FROM GITHUB_DETAIL WHERE RESOURCE_ID=$1`
	githubDetail := GithubDetail{}
	// DB.QueryRow(sqlStatement, resourceID).Scan(&githubDetails.ResourceID, &githubDetails.Owner, &githubDetails.RepositoryName, &githubDetails.Path, &githubDetails.ReadmePath)
	res := DB.Where("resource_id=?", resourceID).Find(&githubDetail)
	if res.Error != nil {
		log.Println(res.Error)
	}
	return githubDetail
}
