package models

import "log"

// import "log"

// UserRating represents relationship between User and Rating
type UserRating struct {
	UserID     int `json:"user_id"`
	ResourceID int `json:"resource_id"`
	Stars      int `json:"stars"`
}

// UpdateResourceRating will update existing rating
func UpdateResourceRating(userID int, resourceID int, stars int, prevStars int) UpdatedRatingResponse {
	// sqlStatement := `UPDATE USER_RATING SET STARS=$3 WHERE RESOURCE_ID=$2 AND USER_ID=$1`
	// _, err := DB.Exec(sqlStatement, userID, resourceID, stars)
	res := DB.Model(&UserRating{}).Where("resource_id=? AND user_id=?", resourceID, userID).Update("stars", stars)
	if res.Error != nil {
		log.Println(res.Error)
	}
	updateStars(resourceID, stars, prevStars)
	averageRating := calculateAverageRating(resourceID)
	resource := Resource{}
	resource.UpdateRating(resourceID, averageRating)
	// updateAverageRating(resourceID, averageRating)
	return updatedRatings(userID, resourceID)
}

// GetUserRating queries for user rating by id
func GetUserRating(userID int, resourceID int) UserRating {
	userRating := UserRating{}
	// sqlStatement := `SELECT * FROM USER_RATING WHERE RESOURCE_ID=$2 AND USER_ID=$1`
	// err := DB.QueryRow(sqlStatement, userID, resourceID).Scan(&userRating.UserID, &userRating.ResourceID, &userRating.Stars)
	res := DB.Where("resource_id=? AND user_id=?", resourceID, userID).Find(&userRating)
	if res.Error != nil {
		log.Println(res.Error)
	}
	return userRating
}

func updatedRatings(userID int, resourceID int) UpdatedRatingResponse {
	rating := Rating{}
	// sqlStatement := `SELECT * FROM RATING WHERE RESOURCE_ID=$1`
	// err := DB.QueryRow(sqlStatement, resourceID).Scan(&rating.ID, &rating.ResourceID, &rating.OneStar, &rating.TwoStar, &rating.ThreeStar, &rating.FourStar, &rating.FiveStar)
	rating = GetRatingDetialsByResourceID(resourceID)
	// if err != nil {
	// 	log.Println(err)
	// }
	updatedRatingResponse := UpdatedRatingResponse{
		OneStar:    rating.OneStar,
		TwoStar:    rating.TwoStar,
		ThreeStar:  rating.ThreeStar,
		FourStar:   rating.FourStar,
		FiveStar:   rating.FiveStar,
		ResourceID: rating.ResourceID,
		Average:    getRatingFromResource(resourceID),
	}
	// updatedRatingResponse.OneStar = rating.OneStar
	// updatedRatingResponse.TwoStar = rating.TwoStar
	// updatedRatingResponse.ThreeStar = rating.ThreeStar
	// updatedRatingResponse.FourStar = rating.FourStar
	// updatedRatingResponse.FiveStar = rating.FiveStar
	// updatedRatingResponse.ResourceID = resourceID
	// updatedRatingResponse.Average = getRatingFromResource(resourceID)
	return updatedRatingResponse
}

// Add a new user rating
func (userRating *UserRating) Add() error {
	res := DB.Create(userRating)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
