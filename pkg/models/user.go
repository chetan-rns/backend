package models

import "log"

// import "log"

// // User represents User model in database
// type User struct {
// 	ID         int    `json:"id"`
// 	FirstName  string `json:"username"`
// 	SecondName string `json:"password"`
// 	EMAIL      string `json:"email"`
// }

// UserCredential represents user details in DB
type UserCredential struct {
	ID        int `gorm:"primary_key"`
	Username  string
	FirstName string
	LastName  string
	Email     string
	Token     string
	Rating    []Rating `gorm:"many2many:user_rating;"`
	Resource  []Resource
}

// UserResourceResponse represents all tasks uploaded by user
type UserResourceResponse struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Rating    float64 `json:"rating"`
	Downloads int     `json:"downloads"`
}

// GetAllResourcesByUser will return all tasks uploaded by user
func GetAllResourcesByUser(userID int) []UserResourceResponse {
	sqlStatement := `SELECT ID,NAME,DOWNLOADS,RATING FROM RESOURCE T JOIN USER_RESOURCE
	U ON T.ID=U.RESOURCE_ID WHERE U.USER_ID=$1`
	// rows, err := DB.Query(sqlStatement, userID)
	rows, err := DB.Raw(sqlStatement, userID).Rows()
	if err != nil {
		log.Println(err)
	}
	resources := []UserResourceResponse{}
	for rows.Next() {
		// var id int
		// var name string
		// var rating float64
		// var downloads int
		resource := UserResourceResponse{}
		rows.Scan(&resource.ID, &resource.Name, &resource.Downloads, &resource.Downloads)
		// task := UserTaskResponse{id, name, rating, downloads}
		resources = append(resources, resource)
	}
	return resources
}

// // GetGithubToken will return github token by ID
// func GetGithubToken(userID int) string {
// 	var token string
// 	sqlStatement := `SELECT TOKEN FROM USER_CREDENTIAL WHERE ID=$1`
// 	DB.QueryRow(sqlStatement, userID).Scan(&token)
// 	return token
// }

// // AddResourceRawPath will add a raw path for resource
// func AddResourceRawPath(resourcePath string, resourceID int, resourceType string) {
// 	sqlStatement := `INSERT INTO RESOURCE_RAW_PATH(RESOURCE_ID,RAW_PATH,TYPE) VALUES($1,$2,$3)`
// 	_, err := DB.Exec(sqlStatement, resourceID, resourcePath, resourceType)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }
