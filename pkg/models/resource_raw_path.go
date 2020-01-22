package models

import (
	"log"
)

// ResourceRawPath represents ResourceRawPath table in DB
type ResourceRawPath struct {
	ResourceID int
	RawPath    string
	Type       string
}

// Add a new resource raw path
func (resourceRawPath *ResourceRawPath) Add() {
	res := DB.Create(resourceRawPath)
	if res.Error != nil {
		log.Println(res.Error)
	}
}

// GetResourceRawLinksByID will query for resource links by ID
func GetResourceRawLinksByID(resourceID int) []ResourceRawPath {
	resourceRawPaths := []ResourceRawPath{}
	DB.Where("resource_id=?", resourceID).Find(&resourceRawPaths)
	return resourceRawPaths
}

// GetResourceRawLinks will return raw github links by ID
func GetResourceRawLinks(resourceID int) RawLinksResponse {
	// sqlStatement := `SELECT * FROM RESOURCE_RAW_PATH WHERE RESOURCE_ID=$1`
	// rows, err := DB.Query(sqlStatement, resourceID)
	// if err != nil {
	// 	log.Println(err)
	// }
	resourceRawPaths := GetResourceRawLinksByID(resourceID)
	log.Print(resourceRawPaths)
	links := RawLinksResponse{}
	for _, resourceRawPath := range resourceRawPaths {
		// var link string
		// var rawResourceType string
		// var id int
		// rows.Scan(&id, &link, &rawResourceType)
		if resourceRawPath.Type == "task" {
			links.Tasks = append(links.Tasks, resourceRawPath.RawPath)
		} else if resourceRawPath.Type == "pipeline" {
			links.Pipelines = append(links.Pipelines, resourceRawPath.RawPath)
		}
	}
	return links
}
