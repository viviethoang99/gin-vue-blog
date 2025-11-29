package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	STATUS_PUBLIC = iota + 1 // Public
	STATUS_SECRET            // Private
	STATUS_DRAFT             // Draft
)

const (
	TYPE_ORIGINAL  = iota + 1 // Original
	TYPE_REPRINT              // Reprint
	TYPE_TRANSLATE            // Translation
)

// belongTo: An article belongs to a category
// belongTo: An article belongs to a user
// many2many: An article can have multiple tags; multiple articles can share a tag
type Article struct {
	Model

	Title       string `gorm:"type:varchar(100);not null" json:"title"`
	Desc        string `json:"desc"`
	Content     string `json:"content"`
	Img         string `json:"img"`
	Type        int    `gorm:"type:tinyint;comment:Type (1-Original 2-Reprint 3-Translation)" json:"type"` // 1-Original 2-Reprint 3-Translation
	Status      int    `gorm:"type:tinyint;comment:Status (1-Public 2-Private)" json:"status"`    // 1-Public 2-Private
	IsTop       bool   `json:"is_top"`
	IsDelete    bool   `json:"is_delete"`
	OriginalUrl string `json:"original_url"`

	CategoryId int `json:"category_id"`
	UserId     int `json:"-"` // user_auth_id

	Tags     []*Tag    `gorm:"many2many:article_tag;joinForeignKey:article_id" json:"tags"`
	Category *Category `gorm:"foreignkey:CategoryId" json:"category"`
	User     *UserAuth `gorm:"foreignkey:UserId" json:"user"`
}

type ArticleTag struct {
	ArticleId int
	TagId     int
}

type BlogArticleVO struct {
	Article

	CommentCount int64 `json:"comment_count"` // Comment count
	LikeCount    int64 `json:"like_count"`    // Like count
	ViewCount    int64 `json:"view_count"`    // View count

	LastArticle       ArticlePaginationVO  `gorm:"-" json:"last_article"`       // Previous
	NextArticle       ArticlePaginationVO  `gorm:"-" json:"next_article"`       // Next
	RecommendArticles []RecommendArticleVO `gorm:"-" json:"recommend_articles"` // Recommended articles
	NewestArticles    []RecommendArticleVO `gorm:"-" json:"newest_articles"`    // Latest articles
}

type ArticlePaginationVO struct {
	ID    int    `json:"id"`
	Img   string `json:"img"`
	Title string `json:"title"`
}

type RecommendArticleVO struct {
	ID        int       `json:"id"`
	Img       string    `json:"img"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

// Article details
func GetArticle(db *gorm.DB, id int) (data *Article, err error) {
	result := db.Preload("Category").Preload("Tags").
		Where(Article{Model: Model{ID: id}}).
		First(&data)
	return data, result.Error
}

// Get first accessible article (not in recycle bin and public)
func GetBlogArticle(db *gorm.DB, id int) (data *Article, err error) {
	result := db.Preload("Category").Preload("Tags").
		Where(Article{Model: Model{ID: id}}).
		Where("is_delete = 0 AND status = 1"). // *
		First(&data)
	return data, result.Error
}

// Frontend article list (not in recycle bin and public)
func GetBlogArticleList(db *gorm.DB, page, size, categoryId, tagId int) (data []Article, total int64, err error) {
	db = db.Model(Article{})
	db = db.Where("is_delete = 0 AND status = 1") // *

	if categoryId != 0 {
		db = db.Where("category_id = ?", categoryId)
	}
	if tagId != 0 {
		db = db.Where("id IN (SELECT article_id FROM article_tag WHERE tag_id = ?)", tagId)
	}

	db = db.Count(&total)
	result := db.Preload("Tags").Preload("Category").
		Order("is_top DESC, id DESC").
		Scopes(Paginate(page, size)).
		Find(&data)

	return data, total, result.Error
}

func GetArticleList(db *gorm.DB, page, size int, title string, isDelete *bool, status, typ, categoryId, tagId int) (list []Article, total int64, err error) {
	db = db.Model(Article{})

	if title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}
	if isDelete != nil {
		db = db.Where("is_delete", isDelete)
	}
	if status != 0 {
		db = db.Where("status", status)
	}
	if categoryId != 0 {
		db = db.Where("category_id", categoryId)
	}
	if typ != 0 {
		db = db.Where("type", typ)
	}

	db = db.Preload("Category").Preload("Tags").
		Joins("LEFT JOIN article_tag ON article_tag.article_id = article.id").
		Group("id") // Deduplicate
	if tagId != 0 {
		db = db.Where("tag_id = ?", tagId)
	}

	result := db.Count(&total).
		Scopes(Paginate(page, size)).
		Order("is_top DESC, article.id DESC").
		Find(&list)
	return list, total, result.Error
}

// Query n recommended articles (by tags)
func GetRecommendList(db *gorm.DB, id, n int) (list []RecommendArticleVO, err error) {
	// sub1: find tag ID list
	// SELECT tag_id FROM `article_tag` WHERE `article_id` = ?
	sub1 := db.Table("article_tag").
		Select("tag_id").
		Where("article_id", id)
	// sub2: find article IDs for these tags (distinct, excluding current article)
	// SELECT DISTINCT article_id FROM (sub1) t
	// JOIN article_tag t1 ON t.tag_id = t1.tag_id
	// WHERE `article_id` != ?
	sub2 := db.Table("(?) t1", sub1).
		Select("DISTINCT article_id").
		Joins("JOIN article_tag t ON t.tag_id = t1.tag_id").
		Where("article_id != ?", id)
	// Fetch article info by article ID list (first n)
	result := db.Table("(?) t2", sub2).
		Select("id, title, img, created_at").
		Joins("JOIN article a ON t2.article_id = a.id").
		Where("a.is_delete = 0").
		Order("is_top, id DESC").
		Limit(n).
		Find(&list)
	return list, result.Error
}

// Query previous article (id < current article id)
func GetLastArticle(db *gorm.DB, id int) (val ArticlePaginationVO, err error) {
	sub := db.Table("article").Select("max(id)").Where("id < ?", id)
	result := db.Table("article").
		Select("id, title, img").
		Where("is_delete = 0 AND status = 1 AND id = (?)", sub).
		Limit(1).
		Find(&val)
	return val, result.Error
}

// Query next article (id > current article id)
func GetNextArticle(db *gorm.DB, id int) (data ArticlePaginationVO, err error) {
	result := db.Model(&Article{}).
		Select("id, title, img").
		Where("is_delete = 0 AND status = 1 AND id > ?", id).
		Limit(1).
		Find(&data)
	return data, result.Error
}

// Query latest n articles
func GetNewestList(db *gorm.DB, n int) (data []RecommendArticleVO, err error) {
	result := db.Model(&Article{}).
		Select("id, title, img, created_at").
		Where("is_delete = 0 AND status = 1").
		Order("created_at DESC, id ASC").
		Limit(n).
		Find(&data)
	return data, result.Error
}

// Physically delete articles
func DeleteArticle(db *gorm.DB, ids []int) (int64, error) {
	// Delete [article-tag] relations
	result := db.Where("article_id IN ?", ids).Delete(&ArticleTag{})
	if result.Error != nil {
		return 0, result.Error
	}

	// Delete [articles]
	result = db.Where("id IN ?", ids).Delete(&Article{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// Soft delete articles (update)
func UpdateArticleSoftDelete(db *gorm.DB, ids []int, isDelete bool) (int64, error) {
	result := db.Model(Article{}).
		Where("id IN ?", ids).
		Update("is_delete", isDelete)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// Create or edit an article, and maintain relations by category name and tag names
func SaveOrUpdateArticle(db *gorm.DB, article *Article, categoryName string, tagNames []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Create category if it does not exist
		category := Category{Name: categoryName}
		result := db.Model(&Category{}).Where("name", categoryName).FirstOrCreate(&category)
		if result.Error != nil {
			return result.Error
		}
		article.CategoryId = category.ID

		// First create/update the article to get its ID
		if article.ID == 0 {
			result = db.Create(&article)
		} else {
			result = db.Model(&article).Where("id", article.ID).Updates(article)
		}
		if result.Error != nil {
			return result.Error
		}

		// Clear article-tag associations
		result = db.Delete(ArticleTag{}, "article_id", article.ID)
		if result.Error != nil {
			return result.Error
		}

		var articleTags []ArticleTag
		for _, tagName := range tagNames {
			// Create tag if it does not exist
			tag := Tag{Name: tagName}
			result := db.Model(&Tag{}).Where("name", tagName).FirstOrCreate(&tag)
			if result.Error != nil {
				return result.Error
			}
			articleTags = append(articleTags, ArticleTag{
				ArticleId: article.ID,
				TagId:     tag.ID,
			})
		}
		result = db.Create(&articleTags)
		return result.Error
	})
}

func UpdateArticleTop(db *gorm.DB, id int, isTop bool) error {
	result := db.Model(&Article{Model: Model{ID: id}}).Update("is_top", isTop)
	return result.Error
}

func ImportArticle(db *gorm.DB, userAuthId int, title string, content string, img string ,categoryname string,tangname string) error {
	article := Article{
		Title:   title,
		Content: content,
		Img:     img,
		Status:  STATUS_DRAFT,
		Type:    TYPE_ORIGINAL,
		UserId:  userAuthId,
	}
	category := Category{Name : categoryname}
	result := db.Model(&Category{}).Where("name",categoryname).FirstOrCreate(&category)
	if result.Error != nil{
			return result.Error
		}
		article.CategoryId =category.ID
		
	result = db.Create(&article)		
	if result.Error!= nil{
		return result.Error
	}

	var articletag ArticleTag
	tag := Tag{Name: tangname}	
	result = db.Model(&Tag{}).Where("name",tangname).FirstOrCreate(&tag)
	if result.Error!= nil{
		return result.Error
	}

	articletag.ArticleId = article.ID
	articletag.TagId = tag.ID
	result = db.Create(&articletag)
	
	return result.Error
}
