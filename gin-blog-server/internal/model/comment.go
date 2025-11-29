package model

import (
	"gorm.io/gorm"
)

const (
	TYPE_ARTICLE = iota + 1 // Article
	TYPE_LINK               // Friend link
	TYPE_TALK               // Talk
)

/*
If the comment type is an article, then topic_id is the article ID.
If the comment type is a friend link, topic_id is not required.
*/

type Comment struct {
	Model
	UserId      int    `json:"user_id"`       // Comment author
	ReplyUserId int    `json:"reply_user_id"` // Replied-to user
	TopicId     int    `json:"topic_id"`      // Article being commented
	ParentId    int    `json:"parent_id"`     // Parent comment (the comment being replied to)
	Content     string `gorm:"type:varchar(500);not null" json:"content"`
	Type        int    `gorm:"type:tinyint(1);not null;comment:Comment type (1.Article 2.Friend link 3.Talk)" json:"type"` // Comment type 1.Article 2.Friend link 3.Talk
	IsReview    bool   `json:"is_review"`

	// Belongs To
	User      *UserAuth `gorm:"foreignKey:UserId" json:"user"`
	ReplyUser *UserAuth `gorm:"foreignKey:ReplyUserId" json:"reply_user"`
	Article   *Article  `gorm:"foreignKey:TopicId" json:"article"`
}

type CommentVO struct {
	Comment
	LikeCount  int         `json:"like_count" gorm:"-"`
	ReplyCount int         `json:"reply_count" gorm:"-"`
	ReplyList  []CommentVO `json:"reply_list" gorm:"-"`
}

// Add a comment
func AddComment(db *gorm.DB, userId, typ, topicId int, content string, isReview bool) (*Comment, error) {
	comment := Comment{
		UserId:   userId,
		TopicId:  topicId,
		Content:  content,
		Type:     typ,
		IsReview: isReview,
	}
	result := db.Create(&comment)
	return &comment, result.Error
}

// Reply to a comment
func ReplyComment(db *gorm.DB, userId, replyUserId, parentId int, content string, isReview bool) (*Comment, error) {
	var parent Comment
	result := db.First(&parent, parentId)
	if result.Error != nil {
		return nil, result.Error
	}

	comment := Comment{
		UserId:      userId,
		Content:     content,
		ReplyUserId: replyUserId,
		ParentId:    parentId,
		IsReview:    isReview,
		TopicId:     parent.TopicId, // Topic same as parent comment
		Type:        parent.Type,    // Type same as parent comment
	}
	result = db.Create(&comment)
	return &comment, result.Error
}

// Get admin comment list
func GetCommentList(db *gorm.DB, page, size, typ int, isReview *bool, nickname string) (data []Comment, total int64, err error) {
	
	// SELECT UID FROM user_info WHERE nikename LIKE nickname
	var uid int
	if nickname != "" {
		result := db.Model(&UserInfo{}).Where("nickname LIKE ?",nickname).Pluck("id",&uid)
		if result.Error != nil{
			return nil,0,result.Error
		}
		db = db.Where("user_id = ?",uid)
	}

	if typ != 0 {
		db = db.Where("type = ?", typ)
	}
	if isReview != nil {
		db = db.Where("is_review = ?", *isReview)
	}

	result := db.Model(&Comment{}).
		Count(&total).
		Preload("User").Preload("User.UserInfo").
		Preload("ReplyUser").Preload("ReplyUser.UserInfo").
		Preload("Article").
		Order("id DESC").
		Scopes(Paginate(page, size)).
		Find(&data)

	return data, total, result.Error
}

// Get blog comment list
func GetCommentVOList(db *gorm.DB, page, size, topic, typ int) (data []CommentVO, total int64, err error) {
	var list []Comment

	tx := db.Model(&Comment{})
	if typ != 0 {
		tx = tx.Where("type = ?", typ)
	}
	if topic != 0 {
		tx = tx.Where("topic_id = ?", topic)
	}

	// Get top-level comments
	tx.Where("parent_id = 0").
		Count(&total).
		Preload("User").Preload("User.UserInfo").
		// Preload("ReplyUser").Preload("ReplyUser.UserInfo").
		Order("id DESC").
		Scopes(Paginate(page, size))
	if err := tx.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	// Get replies for each top-level comment
	for _, v := range list {
		replyList := make([]CommentVO, 0)

		tx := db.Model(&Comment{})
		tx.Where("parent_id = ?", v.ID).
			Preload("User").Preload("User.UserInfo").
			// Preload("ReplyUser").Preload("ReplyUser.UserInfo")
			Order("id DESC")
		if err := tx.Find(&replyList).Error; err != nil {
			return nil, 0, err
		}

		data = append(data, CommentVO{
			ReplyCount: len(replyList),
			Comment:    v,
			ReplyList:  replyList,
		})
	}

	return data, total, nil
}

// Get reply list by comment ID
func GetCommentReplyList(db *gorm.DB, id, page, size int) (data []Comment, err error) {
	result := db.Model(&Comment{}).
		Where(&Comment{ParentId: id}).
		Preload("User").Preload("User.UserInfo").
		Order("id DESC").
		Scopes(Paginate(page, size)).
		Find(&data)
	return data, result.Error
}

// Get the number of comments for an article
func GetArticleCommentCount(db *gorm.DB, articleId int) (count int64, err error) {
	result := db.Model(&Comment{}).
		Where("topic_id = ? AND type = 1 AND is_review = 1", articleId).
		Count(&count)
	return count, result.Error
}
