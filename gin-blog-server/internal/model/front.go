package model

import (
	"gorm.io/gorm"
)

type FrontHomeVO struct {
	ArticleCount  int64             `json:"article_count"`  // Number of articles
	UserCount     int64             `json:"user_count"`     // Number of users
	MessageCount  int64             `json:"message_count"`  // Number of messages
	CategoryCount int64             `json:"category_count"` // Number of categories
	TagCount      int64             `json:"tag_count"`      // Number of tags
	ViewCount     int64             `json:"view_count"`     // Page views
	Config        map[string]string `json:"blog_config"`    // Blog configuration
	// PageList      []Page            `json:"page_list"`      // Page list
}

func GetFrontStatistics(db *gorm.DB) (data FrontHomeVO, err error) {
	result := db.Model(&Article{}).Where("status = ? AND is_delete = ?", 1, 0).Count(&data.ArticleCount)
	if result.Error != nil {
		return data, result.Error
	}

	result = db.Model(&UserAuth{}).Count(&data.UserCount)
	if result.Error != nil {
		return data, result.Error
	}

	result = db.Model(&Message{}).Count(&data.MessageCount)
	if result.Error != nil {
		return data, result.Error
	}

	result = db.Model(&Category{}).Count(&data.CategoryCount)
	if result.Error != nil {
		return data, result.Error
	}

	result = db.Model(&Tag{}).Count(&data.TagCount)
	if result.Error != nil {
		return data, result.Error
	}

	data.Config, err = GetConfigMap(db)
	if err != nil {
		return data, err
	}

	return data, nil
}
