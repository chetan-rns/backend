package models

import (
	"database/sql"
	"fmt"
	"log"
)

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// )

// // TaskTags represents many-many between Task and Tag models
// type TaskTags struct {
// 	TaskID int `json:"taskID"`
// 	TagID  int `json:"tagID"`
// }

// ResourceTag represents ResourceTag table in DB
type ResourceTag struct {
	ResourceID int
	TagID      int
}

// Add a new resource tag
func (resourceTag *ResourceTag) Add() {
	res := DB.Create(resourceTag)
	if res.Error != nil {
		log.Println(res.Error)
	}
}

// GetAllResourcesWithGivenTags queries for all resources with given tags
func GetAllResourcesWithGivenTags(resourceType string, verified string, tags []string) []ResourceResponse {
	resources := []ResourceResponse{}
	args := make([]interface{}, len(tags))
	for index, value := range tags {
		args[index] = value
	}
	params := `$1`
	for index := 1; index <= len(tags); index++ {
		if index > 1 {
			params = params + fmt.Sprintf(",$%d", index)
		}
	}
	var (
		resourceTagMap map[int][]string
		rows           *sql.Rows
		err            error
	)
	resourceTagMap = getResourceTagMap()
	rows, err = executeTagsQuery(tags, params, args)
	defer rows.Close()
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		resource := Resource{}
		err = rows.Scan(&resource.ID, &resource.Name, &resource.Type, &resource.Description, &resource.Downloads, &resource.Rating, &resource.Github, &resource.Verified)
		if err != nil {
			log.Println(err)
		}
		resourceResponse := resource.GetResourceResponse()
		resourceResponse.Tags = resourceTagMap[resource.ID]
		matchTypeAndVerified(resourceType, verified, resourceResponse, &resources)
	}
	return resources
}

func executeTagsQuery(tags []string, params string, args []interface{}) (*sql.Rows, error) {
	var (
		rows         *sql.Rows
		err          error
		sqlStatement string
	)
	if len(tags) > 0 {
		sqlStatement = `
	SELECT DISTINCT T.ID,T.NAME,T.TYPE,T.DESCRIPTION,T.DOWNLOADS,T.RATING,T.GITHUB,T.VERIFIED
	FROM RESOURCE AS T JOIN RESOURCE_TAG AS TT ON (T.ID=TT.RESOURCE_ID) JOIN TAG
	AS TG ON (TG.ID=TT.TAG_ID AND TG.NAME in (` +
			params + `));`
		rows, err = DB.Raw(sqlStatement, args...).Rows()
	} else {
		sqlStatement = `
	SELECT DISTINCT T.ID,T.NAME,T.TYPE,T.DESCRIPTION,T.DOWNLOADS,T.RATING,T.GITHUB,T.VERIFIED
	FROM RESOURCE T`
		rows, err = DB.Raw(sqlStatement).Rows()
	}
	return rows, err
}

func matchTypeAndVerified(resourceType string, verified string, resource ResourceResponse, resources *[]ResourceResponse) {
	isVerified := getBoolString(resource.Verified)
	if resourceType != "all" && verified != "all" {
		if resourceType == resource.Type && isVerified == verified {
			*resources = append(*resources, resource)
		}
	} else if resourceType == "all" && verified != "all" {
		if isVerified == verified {
			*resources = append(*resources, resource)
		}
	} else if resourceType != "all" && verified == "all" {
		if resourceType == resource.Type {
			*resources = append(*resources, resource)
		}
	} else {
		*resources = append(*resources, resource)
	}
}

func getBoolString(p bool) string {
	if p == true {
		return "true"
	} else if p == false {
		return "false"
	}
	return "all"
}

func getResourceTagMap() map[int][]string {
	sqlStatement := `SELECT DISTINCT T.ID,TG.NAME FROM RESOURCE AS T JOIN RESOURCE_TAG AS TT ON (T.ID=TT.RESOURCE_ID) JOIN TAG AS TG ON (TG.ID=TT.TAG_ID);`
	// rows, err := DB.Query(sqlStatement)
	rows, err := DB.Raw(sqlStatement).Rows()
	// mapping task ID with tag names
	var resourceTagMap map[int][]string
	resourceTagMap = make(map[int][]string)
	for rows.Next() {
		var resourceID int
		var tagName string
		err = rows.Scan(&resourceID, &tagName)
		if err != nil {
			log.Println(err)
		}
		resourceTagMap[resourceID] = append(resourceTagMap[resourceID], tagName)
	}
	return resourceTagMap
}
