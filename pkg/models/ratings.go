package models

import (
	"fmt"
	"log"
)

// Rating represents Rating model in database
type Rating struct {
	ID         int `gorm:"primary_key;auto_increment" json:"id"`
	ResourceID int `json:"resource_id"`
	OneStar    int `gorm:"default=0" json:"one"`
	TwoStar    int `gorm:"default=0" json:"two"`
	ThreeStar  int `gorm:"default=0" json:"three"`
	FourStar   int `gorm:"default=0" json:"four"`
	FiveStar   int `gorm:"default=0" json:"five"`
}

// PrevStarRequest represents previous stars
type PrevStarRequest struct {
	UserID     int `json:"user_id"`
	ResourceID int `json:"resource_id"`
}

// GetRatingDetialsByResourceID retrieves rating details of a resource
func GetRatingDetialsByResourceID(resourceID int) Rating {
	// sqlStatement := `SELECT * FROM RATING WHERE RESOURCE_ID=$1`
	rating := Rating{}
	// err := DB.QueryRow(sqlStatement, resourceID).Scan(&taskRating.ID, &taskRating.ResourceID, &taskRating.OneStar, &taskRating.TwoStar, &taskRating.ThreeStar, &taskRating.FourStar, &taskRating.FiveStar)
	DB.Where("resource_id=?", resourceID).First(&rating).Row().Scan(&rating.ID, &rating.ResourceID, &rating.OneStar, &rating.TwoStar, &rating.ThreeStar, &rating.FourStar, &rating.FiveStar)
	// DB.Raw(sqlStatement, resourceID).Row().Scan(&rating.ID, &rating.ResourceID, &rating.OneStar, &rating.TwoStar, &rating.ThreeStar, &rating.FourStar, &rating.FiveStar)
	// if res.Error != nil {
	// 	log.Println(res.Error)
	// }
	// log.Println(rating, resourceID)
	return rating
}

func calculateAverageRating(resourceID int) float64 {
	rating := Rating{}
	// sqlStatement := `SELECT * FROM RATING WHERE RESOURCE_ID=$1`
	// err := DB.QueryRow(sqlStatement, resourceID).Scan(&rating.ID, &rating.ResourceID, &rating.OneStar, &rating.TwoStar, &rating.ThreeStar, &rating.FourStar, &rating.FiveStar)
	rating = GetRatingDetialsByResourceID(resourceID)

	totalResponses := float64(rating.OneStar + rating.TwoStar + rating.ThreeStar + rating.FourStar + rating.FiveStar)
	averageRating := float64(rating.OneStar+rating.TwoStar*2+rating.ThreeStar*3+rating.FourStar*4+rating.FiveStar*5) / (totalResponses)
	log.Println(averageRating)
	return averageRating
}

func getStarsInString(stars int) string {
	switch stars {
	case 1:
		return "onestar"
	case 2:
		return "twostar"
	case 3:
		return "threestar"
	case 4:
		return "fourstar"
	case 5:
		return "fivestar"
	}
	return ""
}

func addStars(resourceID int, stars int, prevStars int) error {
	starsString := getStarsInString(stars)
	sqlStatement := fmt.Sprintf("INSERT INTO RATING(%v,RESOURCE_ID) VALUES($1,$2) ON CONFLICT (RESOURCE_ID) DO UPDATE SET %v=RATING.%v+1", starsString, starsString, starsString)
	// _, err := DB.Exec(sqlStatement, 1, taskID)
	res := DB.Exec(sqlStatement, 1, resourceID)
	if res.Error != nil {
		log.Println(res.Error)
		return res.Error
	}
	return nil
}

func updateStars(resourceID int, stars int, prevStars int) {
	starsString := getStarsInString(stars)
	sqlStatement := fmt.Sprintf("UPDATE RATING SET %v=%v+1 WHERE RESOURCE_ID=$1", starsString, starsString)
	// _, err := DB.Exec(sqlStatement, resourceID)
	// rating := GetRatingDetialsByResourceID(resourceID)
	// log.Println("Before:", rating)
	// curStars := getCurrentStarsFromRating(rating, stars) + 1
	// DB.Model(&rating).Update(starsString, curStars)
	res := DB.Exec(sqlStatement, resourceID)
	// res := DB.Model(&Rating{}).Where("resource_id = ?", resourceID).Update(starsString, curStars)
	if res.Error != nil {
		log.Println(res.Error)
	}
	// log.Println(starsString, curStars)
	log.Println("After:", GetRatingDetialsByResourceID(resourceID))
	deleteOldStars(resourceID, prevStars)
}

func deleteOldStars(resourceID int, prevStars int) {
	starsString := getStarsInString(prevStars)
	sqlStatement := fmt.Sprintf("UPDATE RATING SET %v=%v-1 WHERE RESOURCE_ID=$1", starsString, starsString)
	// _, err := DB.Exec(sqlStatement, resourceID)
	// rating := GetRatingDetialsByResourceID(resourceID)
	// curStars := getCurrentStarsFromRating(rating, prevStars) - 1
	// res := DB.Model(&rating).Where("resource_id=?", resourceID).Update(starsString, curStars)
	res := DB.Exec(sqlStatement, resourceID)
	if res.Error != nil {
		log.Println(res.Error)
	}
}

func getCurrentStarsFromRating(rating Rating, stars int) int {
	switch stars {
	case 1:
		return rating.OneStar
	case 2:
		return rating.TwoStar
	case 3:
		return rating.ThreeStar
	case 4:
		return rating.FourStar
	case 5:
		return rating.FiveStar
	}
	return 0
}

// AddRating add's rating provided by user
func AddRating(userID int, resourceID int, stars int, prevStars int) interface{} {
	// sqlStatement := `INSERT INTO USER_RATING(USER_ID,RESOURCE_ID,STARS) VALUES($1,$2,$3)`
	// _, err := DB.Exec(sqlStatement, userID, resourceID, stars)
	userRating := UserRating{
		UserID:     userID,
		ResourceID: resourceID,
		Stars:      stars,
	}
	err := userRating.Add()
	if err != nil {
		log.Println(err)
		return map[string]interface{}{"status": false, "message": "Use PUT method to update existing rating"}
	}
	err = addStars(resourceID, stars, prevStars)
	if err != nil {
		return map[string]interface{}{"status": false, "message": "Not able to add stars to ratings table"}
	}
	averageRating := calculateAverageRating(resourceID)
	resource := Resource{}
	err = resource.UpdateRating(resourceID, averageRating)
	// err = updateAverageRating(resourceID, averageRating)
	if err != nil {
		return map[string]interface{}{"status": false, "message": "Unable to update average rating in task"}
	}
	return updatedRatings(userID, resourceID)
}
