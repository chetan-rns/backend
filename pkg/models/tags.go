package models

import "log"

// import "log"

// Tag is a model representing tags associated with tasks
type Tag struct {
	ID   int    `gorm:"primary_key;auto_increment" json:"id"`
	Name string `json:"name"`
}

// GetAllTags will query for all tags
func GetAllTags() []Tag {
	tags := []Tag{}
	// sqlStatement := `
	// SELECT * FROM TAG;`
	// rows, err := DB.Query(sqlStatement)
	// for rows.Next() {
	// 	tag := Tag{}
	// 	err = rows.Scan(&tag.ID, &tag.Name)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	tags = append(tags, tag)
	// }
	res := DB.Find(&tags)
	if res.Error != nil {
		log.Println(res.Error)
	}
	return tags
}

// AddTag will add a new tag
// func AddTag(tag string) (int, error) {
// 	var newTagID int
// 	sqlStatement := `INSERT INTO TAG(NAME) VALUES($1) RETURNING ID`
// 	err := DB.QueryRow(sqlStatement, tag).Scan(&newTagID)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return newTagID, nil
// }

// AddTag will add a new tag to DB
func AddTag(name string) (int, error) {
	// var tagID int
	tag := Tag{}
	res := DB.Create(name).Scan(&tag)

	if res.Error != nil {
		log.Println(res.Error)
		return 0, res.Error
	}
	return tag.ID, nil
}
