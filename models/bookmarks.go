package models

import "github.com/jinzhu/gorm"

type Bookmarks struct {
	UserID     string `json:"userId"`
	QuestionID int64  `json:"questionId"`

	Base
}

// AddBookmark adds the post id and user id in bookmarks table, returns
// no of rows affected and the error.
func AddBookmark(tx *gorm.DB, userID string, postID int64) (int64, error) {
	bookmark := &Bookmarks{
		UserID:     userID,
		QuestionID: postID,
	}

	var err error
	var existing bool

	if existing, err = isAlreadyBookmarked(tx, userID, postID); err != nil {
		return 0, err
	}

	// Fetching the total no of bookmarks
	var total int64
	if total, err = getBookmarkCount(tx, postID); err != nil {
		return 0, err
	}

	// Check whether it already exists and return if it it is.
	if existing {
		return total, nil
	}

	if db := tx.Save(bookmark); db.Error != nil {
		return 0, db.Error
	}

	return total + 1, nil
}

// DeleteBookmark deletes the post from the bookmarks list of the user and
// returns the current number of bookmark.
func DeleteBookmark(tx *gorm.DB, userID string, postID int64) (int64, error) {
	if db := tx.Delete(Bookmarks{}, "question_id = ? and user_id = ?", postID, userID); db.Error != nil {
		return 0, db.Error
	}

	var err error
	var total int64

	if total, err = getBookmarkCount(tx, postID); err != nil {
		return 0, err
	}

	return total, nil
}

func getBookmarkCount(tx *gorm.DB, postId int64) (int64, error) {
	var count int64

	db := tx.Model(Bookmarks{}).Where("question_id = ?", postId).Count(&count)

	if db.Error != nil {
		return count, db.Error
	}

	return count, nil
}

func isAlreadyBookmarked(tx *gorm.DB, userID string, postID int64) (bool, error) {
	var existing int64
	var db *gorm.DB

	db = tx.Model(Bookmarks{}).
		Where("user_id = ? and question_id = ?", userID, postID).
		Count(&existing)

	if db.Error != nil {
		return false, db.Error
	}

	return existing > 0, nil
}
