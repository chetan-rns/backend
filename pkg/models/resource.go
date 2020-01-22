package models

import (
	"log"
)

// ResourceResponse is a database model representing task and pipeline
type ResourceResponse struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Downloads   int      `json:"downloads"`
	Rating      float64  `json:"rating"`
	Github      string   `json:"github"`
	Tags        []string `json:"tags"`
	Verified    bool     `json:"verified"`
}

// Resource represents reosurce model in DB
type Resource struct {
	ID              int     `gorm:"primary_key;auto_increment"`
	Name            string  `gorm:"default:''"`
	Description     string  `gorm:"default:''"`
	Downloads       int     `gorm:"default:0"`
	Rating          float64 `gorm:"default:0.0"`
	Github          string  `gorm:"default:''"`
	Type            string  `gorm:"default:'task'"`
	Verified        bool    `gorm:"default:false"`
	GithubDetail    GithubDetail
	ResourceRawPath []ResourceRawPath
	Tag             []Tag `gorm:"many2many:resource_tag;"`
}

// Add a new resource
func (resource *Resource) Add() (int, error) {
	r := Resource{}
	res := DB.Create(resource).Scan(&r)
	if res.Error != nil {
		return 0, res.Error
	}
	return r.ID, nil
}

// FetchAllResources will fetch all resources from DB
func FetchAllResources() ([]Resource, error) {
	resources := []Resource{}
	// resource := Resource{}
	res := DB.Find(&resources)
	if res.Error != nil {
		return nil, res.Error
	}
	return resources, nil
}

// AddCatalogResource is called to add resource from catalog
func AddCatalogResource(resource *Resource) (int, error) {
	// sqlStatement := `
	// INSERT INTO RESOURCE (NAME,DOWNLOADS,RATING,GITHUB)
	// VALUES ($1, $2, $3, $4) RETURNING ID`
	// err := DB.QueryRow(sqlStatement, resource.Name, resource.Downloads, resource.Rating, resource.Github).Scan(&resourceID)
	resourceID, err := resource.Add()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return resourceID, nil
}

// AddResource will add a new resource
func AddResource(resource *ResourceResponse, userID int, owner string, respositoryName string, path string) (int, error) {
	var resourceID int
	sqlStatement := `
	INSERT INTO RESOURCE (NAME,DESCRIPTION,DOWNLOADS,RATING,GITHUB,TYPE)
	VALUES (?, ?, ?, ?, ?,?) RETURNING ID`
	DB.Exec(sqlStatement, resource.Name, resource.Description, resource.Downloads, resource.Rating, resource.Github, resource.Type).Row().Scan(&resourceID)
	log.Println("dasdasda")
	// Add Tags separately
	if len(resource.Tags) > 0 {
		for _, tag := range resource.Tags {
			tagExists := true
			// Use existing tags if already exists
			var tagID int
			tagObject := Tag{}
			// sqlStatement = `SELECT ID FROM TAG WHERE NAME=$1`
			res := DB.Model(&Tag{}).Where("name=?", tag).First(&tagObject)
			tagID = tagObject.ID
			// err := DB.QueryRow(sqlStatement, tag).Scan(&tagID)
			log.Println(tagID)
			if res.Error != nil {
				tagExists = false
				log.Println(res.Error)
			}
			// If tag already exists
			if tagExists {
				// addResourceTag(resourceID, tagID)
				resourceTag := ResourceTag{
					ResourceID: resourceID,
					TagID:      tagID,
				}
				resourceTag.Add()
			} else {
				var newTagID int
				newTagID, err := AddTag(tag)
				if err != nil {
					log.Println(err)
				}
				// addResourceTag(resourceID, newTagID)
				resourceTag := ResourceTag{
					ResourceID: resourceID,
					TagID:      newTagID,
				}
				resourceTag.Add()
			}
		}
	}
	githubDetail := GithubDetail{
		ResourceID:     resourceID,
		Owner:          owner,
		RepositoryName: respositoryName,
		Path:           path,
	}
	githubDetail.Add()
	// addGithubDetails(resourceID, owner, respositoryName, path)
	return resourceID, addUserResource(userID, resourceID)
}

// func addResourceTag(resourceID int, tagID int) {
// 	sqlStatement := `INSERT INTO RESOURCE_TAG(RESOURCE_ID,TAG_ID) VALUES($1,$2)`
// 	_, err := DB.Exec(sqlStatement, resourceID, tagID)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func addGithubDetails(resourceID int, owner string, respositoryName string, path string) {
// 	sqlStatement := `INSERT INTO GITHUB_DETAIL(RESOURCE_ID,OWNER,REPOSITORY_NAME,PATH) VALUES($1,$2,$3,$4)`
// 	_, err := DB.Exec(sqlStatement, resourceID, owner, respositoryName, path)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func updateGithubYAMLDetails(resourceID int, path string) {
// 	sqlStatement := `UPDATE GITHUB_DETAIL SET PATH=$1 WHERE RESOURCE_ID=$2`
// 	_, err := DB.Exec(sqlStatement, path, resourceID)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func updateGithubREADMEDetails(resourceID int, path string) {
// 	sqlStatement := `UPDATE GITHUB_DETAIL SET README_PATH=$1 WHERE RESOURCE_ID=$2`
// 	_, err := DB.Exec(sqlStatement, path, resourceID)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

func addUserResource(userID int, resourceID int) error {
	// sqlStatement := `INSERT INTO USER_RESOURCE(RESOURCE_ID,USER_ID) VALUES($1,$2)`
	// _, err := DB.Exec(sqlStatement, resourceID, userID)
	// if err != nil {
	// 	return err
	// }
	// return nil
	userResource := UserResource{
		ResourceID: resourceID,
		UserID:     userID,
	}
	return userResource.Add()
}

// CheckSameResourceUpload will checkif the user submitted the same resource again
func CheckSameResourceUpload(userID int, name string) bool {
	sqlStatement := `SELECT T.NAME FROM RESOURCE T JOIN USER_RESOURCE U ON T.ID=U.RESOURCE_ID WHERE U.USER_ID=?`
	// rows, err := DB.Query(sqlStatement, userID)
	// if err != nil {
	// 	log.Println(err)
	// }
	rows, err := DB.Raw(sqlStatement, userID).Rows()
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var resourceName string
		err := rows.Scan(&resourceName)
		if err != nil {
			log.Println(err)
		}
		if resourceName == name {
			return true
		}
	}
	return false
}

// GetAllResources will return all the tasks
func GetAllResources() []ResourceResponse {
	// resources := []Resource{}
	// sqlStatement := `
	// SELECT * FROM RESOURCE ORDER BY ID`
	// rows, err := DB.Query(sqlStatement)
	resources := []Resource{}
	allResources := []ResourceResponse{}
	resources, err := FetchAllResources()
	if err != nil {
		log.Println(err)
	}
	// defer rows.Close()
	// for rows.Next() {
	// 	resource := Resource{}
	// 	err = rows.Scan(&resource.ID, &resource.Name, &resource.Description, &resource.Downloads, &resource.Rating, &resource.Github, &resource.Type, &resource.Verified)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	resources = append(resources, resource)
	// }
	resourceIndexMap := make(map[int]int)
	// sqlStatement = `SELECT ID FROM RESOURCE ORDER BY ID`
	// rows, err = DB.Query(sqlStatement)
	// if err != nil {
	// 	log.Println(err)
	// }
	resourceIndex := 0
	// for rows.Next() {
	// 	var id int
	// 	err = rows.Scan(&id)
	// 	resourceIndexMap[id] = resourceIndex
	// 	resourceIndex = resourceIndex + 1
	// }
	for _, resource := range resources {
		resourceResponse := ResourceResponse{
			ID:          resource.ID,
			Name:        resource.Name,
			Type:        resource.Type,
			Description: resource.Description,
			Downloads:   resource.Downloads,
			Rating:      resource.Rating,
			Github:      resource.Github,
			Verified:    resource.Verified,
		}
		allResources = append(allResources, resourceResponse)
		resourceIndexMap[resource.ID] = resourceIndex
		resourceIndex = resourceIndex + 1
	}

	sqlStatement := `SELECT R.ID,TG.NAME FROM TAG TG JOIN RESOURCE_TAG TT ON TT.TAG_ID=TG.ID JOIN RESOURCE R ON R.ID=TT.RESOURCE_ID`
	rows, err := DB.Raw(sqlStatement).Rows()
	// rows, err = DB.Query(sqlStatement)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var tag string
		var resourceID int
		err := rows.Scan(&resourceID, &tag)
		if err != nil {
			log.Println(err)
		}
		allResources[resourceIndexMap[resourceID]].Tags = append(allResources[resourceIndexMap[resourceID]].Tags, tag)
	}
	defer rows.Close()
	return allResources
}

// GetResourceResponse will convert resource object to ResourceResponse object
func (resource *Resource) GetResourceResponse() ResourceResponse {
	resourceResponse := ResourceResponse{
		ID:          resource.ID,
		Name:        resource.Name,
		Type:        resource.Type,
		Description: resource.Description,
		Downloads:   resource.Downloads,
		Rating:      resource.Rating,
		Github:      resource.Github,
		Verified:    resource.Verified,
	}
	return resourceResponse
}

// GetResourceByID returns a resource with requested ID
func GetResourceByID(id int) ResourceResponse {
	resourceResponse := ResourceResponse{}
	resource := Resource{}
	var resourceTagMap map[int][]string
	resourceTagMap = make(map[int][]string)
	resourceTagMap = getResourceTagMap()
	// sqlStatement := `
	// SELECT * FROM RESOURCE WHERE ID=$1;`
	// err := DB.QueryRow(sqlStatement, id).Scan(&resource.ID, &resource.Name, &resource.Description, &resource.Downloads, &resource.Rating, &resource.Github, &resource.Type, &resource.Verified)
	// if err != nil {
	// 	return Resource{}
	// }
	resourceResponse = ResourceResponse{}
	DB.First(&resource, id)
	resourceResponse = resource.GetResourceResponse()
	resourceResponse.Tags = resourceTagMap[resourceResponse.ID]
	return resourceResponse
}

// // GetTaskNameFromID returns name from given ID
// func GetTaskNameFromID(taskID string) string {
// 	id, err := strconv.Atoi(taskID)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	sqlStatement := `SELECT NAME FROM TASK WHERE ID=$1`
// 	var taskName string
// 	err = DB.QueryRow(sqlStatement, id).Scan(&taskName)
// 	if err != nil {
// 		return ""
// 	}
// 	return taskName
// }

// // IncrementDownloads will increment the number of downloads
// func IncrementDownloads(taskID string) {
// 	id, err := strconv.Atoi(taskID)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	log.Println(id)
// 	sqlStatement := `UPDATE RESOURCE SET DOWNLOADS = DOWNLOADS + 1 WHERE ID=$1`
// 	_, err = DB.Exec(sqlStatement, id)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func updateAverageRating(resourceID int, rating float64) error {
// 	sqlStatement := `UPDATE RESOURCE SET RATING=$2 WHERE ID=$1`
// 	_, err := DB.Exec(sqlStatement, resourceID, rating)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	return nil
// }

// UpdateRating will update the rating in resource
func (resource *Resource) UpdateRating(resourceID int, rating float64) error {
	res := DB.Model(resource).Where("id=?", resourceID).Update("rating", rating)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func getRatingFromResource(resourceID int) float64 {
	// var rating float64
	// sqlStatement := `SELECT RATING FROM RESOURCE WHERE ID=$1`
	// err := DB.QueryRow(sqlStatement, resourceID).Scan(&rating)
	resource := Resource{}
	res := DB.Where("id=?", resourceID).First(&resource)
	if res.Error != nil {
		log.Println(res.Error)
	}
	return resource.Rating
}

// // GetResourceIDFromName will return resource ID from name
// func GetResourceIDFromName(name string) (int, error) {
// 	sqlStatement := `SELECT ID FROM RESOURCE WHERE NAME=$1`
// 	var resourceID int
// 	err := DB.QueryRow(sqlStatement, name).Scan(&resourceID)
// 	if err != nil {
// 		log.Println(err)
// 		return 0, err
// 	}
// 	return resourceID, nil
// }

// func resourceExists(resourceName string) bool {
// 	sqlStatement := `SELECT EXISTS(SELECT 1 FROM RESOURCE WHERE NAME=$1 AND VERIFIED=$2)`
// 	var exists bool
// 	err := DB.QueryRow(sqlStatement, resourceName, true).Scan(&exists)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return exists
// }

// // DeleteResource will delete a resource
// func DeleteResource(resourceID int) error {
// 	sqlStatement := `DELETE FROM RESOURCE WHERE ID=$1`
// 	_, err := DB.Exec(sqlStatement, resourceID)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	return nil
// }
