package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"strconv"
	"strings"
	"time"
	"html/template"

	"github.com/gin-gonic/gin"
)

type Front struct{}

type FAddMessageReq struct {
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar"`
	Content  string `json:"content" binding:"required"`
	Speed    int    `json:"speed"`
}

type FAddCommentReq struct {
	ReplyUserId int    `json:"reply_user_id" form:"reply_user_id"`
	TopicId     int    `json:"topic_id" form:"topic_id"`
	Content     string `json:"content" form:"content"`
	ParentId    int    `json:"parent_id" form:"parent_id"`
	Type        int    `json:"type" form:"type" validate:"required,min=1,max=3" label:"Comment Type"`
}

type FCommentQuery struct {
	PageQuery
	ReplyUserId int    `json:"reply_user_id" form:"reply_user_id"`
	TopicId     int    `json:"topic_id" form:"topic_id"`
	Content     string `json:"content" form:"content"`
	ParentId    int    `json:"parent_id" form:"parent_id"`
	Type        int    `json:"type" form:"type"`
}

type FArticleQuery struct {
	PageQuery
	CategoryId int `form:"category_id"`
	TagId      int `form:"tag_id"`
}

type ArchiveVO struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Created_at time.Time `json:"created_at"`
}

type ArticleSearchVO struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Front page info
func (*Front) GetHomeInfo(c *gin.Context) {
	db := GetDB(c)
	rdb := GetRDB(c)

	data, err := model.GetFrontStatistics(db)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	data.ViewCount, _ = rdb.Get(rctx, g.VIEW_COUNT).Int64()

	ReturnSuccess(c, data)
}

// Get tag list
func (*Front) GetTagList(c *gin.Context) {
	list, _, err := model.GetTagList(GetDB(c), 1, 1000, "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}

// Get category list
func (*Front) GetCategoryList(c *gin.Context) {
	list, _, err := model.GetCategoryList(GetDB(c), 1, 1000, "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}

// Get message list
func (*Front) GetMessageList(c *gin.Context) {
	isReview := true
	list, _, err := model.GetMessageList(GetDB(c), 1, 1000, "", &isReview)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}

// Get friend link list
func (*Front) GetLinkList(c *gin.Context) {
	list, _, err := model.GetLinkList(GetDB(c), 1, 1000, "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, list)
}

/*
// The following APIs require login
*/

// TODO: Add avatar/nickname for messages (allow anonymous messages)
// Save message (add only, no edit)
func (*Front) SaveMessage(c *gin.Context) {
	var req FAddMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	req.Content = template.HTMLEscapeString(req.Content)
	auth, _ := CurrentUserAuth(c)
	db := GetDB(c)

	ipAddress := utils.IP.GetIpAddress(c)
	ipSource := utils.IP.GetIpSource(ipAddress)
	isReview := model.GetConfigBool(db, g.CONFIG_IS_COMMENT_REVIEW)

	info := auth.UserInfo
	message, err := model.SaveMessage(db, info.Nickname, info.Nickname, req.Content, ipAddress, ipSource, req.Speed, isReview)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, message)
}

// Save comment (add only, no edit)
// TODO: Add avatar/nickname for comments (allow anonymous comments)
// TODO: Enable email notifications for users
// TODO: HTMLUtil.Filter to sanitize strings in HTML elements...
func (*Front) SaveComment(c *gin.Context) {
	var req FAddCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	// Filter comment content to prevent XSS
	req.Content = template.HTMLEscapeString(req.Content)
	auth, _ := CurrentUserAuth(c)
	db := GetDB(c)
	isReview := model.GetConfigBool(db, g.CONFIG_IS_COMMENT_REVIEW)

	var comment *model.Comment
	var err error

	if req.ReplyUserId == 0 {
		comment, err = model.AddComment(db, auth.ID, req.Type, req.TopicId, req.Content, isReview)
	} else {
		comment, err = model.ReplyComment(db, auth.ID, req.ReplyUserId, req.ParentId, req.Content, isReview)
	}

	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, comment)
}

// Get comment list
func (*Front) GetCommentList(c *gin.Context) {
	var query FCommentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	data, total, err := model.GetCommentVOList(db, query.Page, query.Size, query.TopicId, query.Type)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	likeCountMap := rdb.HGetAll(rctx, g.COMMENT_LIKE_COUNT).Val()
	for i, comment := range data {
		if len(data[i].ReplyList) > 3 {
			data[i].ReplyList = data[i].ReplyList[:3] // show only 3 replies
		}
		data[i].LikeCount, _ = strconv.Atoi(likeCountMap[strconv.Itoa(comment.ID)])
	}

	ReturnSuccess(c, PageResult[model.CommentVO]{
		List:  data,
		Total: total,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// Get reply list by comment id
func (*Front) GetReplyListByCommentId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	replyList, err := model.GetCommentReplyList(db, id, query.Page, query.Size)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	likeCountMap := rdb.HGetAll(rctx, g.COMMENT_LIKE_COUNT).Val()

	data := make([]model.CommentVO, 0)
	for _, reply := range replyList {
		like, _ := strconv.Atoi(likeCountMap[strconv.Itoa(reply.ID)])
		data = append(data, model.CommentVO{
			Comment:   reply,
			LikeCount: like,
		})
	}

	ReturnSuccess(c, data)
}

// Like a comment
func (*Front) LikeComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rdb := GetRDB(c)
	auth, _ := CurrentUserAuth(c)

	// Record that a user liked a specific comment
	commentLikeUserKey := g.COMMENT_USER_LIKE_SET + strconv.Itoa(auth.ID)
	// If already recorded, like toggles to unlike
	if rdb.SIsMember(rctx, commentLikeUserKey, id).Val() {
		rdb.SRem(rctx, commentLikeUserKey, id)
		rdb.HIncrBy(rctx, g.COMMENT_LIKE_COUNT, strconv.Itoa(id), -1)
	} else { // Not recorded, this is a new like
		rdb.SAdd(rctx, commentLikeUserKey, id)
		rdb.HIncrBy(rctx, g.COMMENT_LIKE_COUNT, strconv.Itoa(id), 1)
	}

	ReturnSuccess(c, nil)
}

/*
// Article related APIs
*/

// Get article list
func (*Front) GetArticleList(c *gin.Context) {
	var query FArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	list, _, err := model.GetBlogArticleList(GetDB(c), query.Page, query.Size, query.CategoryId, query.TagId)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, list)
}

// Get article detail by id
func (*Front) GetArticleInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	// Article detail
	val, err := model.GetBlogArticle(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	article := model.BlogArticleVO{Article: *val}

	// Recommended articles (6)
	article.RecommendArticles, err = model.GetRecommendList(db, id, 6)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// Latest articles (5)
	article.NewestArticles, err = model.GetNewestList(db, 5)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// Update article view count TODO: remove view count when article deleted
	// updateArticleViewCount(c, id)

	// TODO: Update visit count
	// * Each request increments visit count; refresh can inflate visits
	rdb.ZIncrBy(rctx, g.ARTICLE_VIEW_COUNT, 1, strconv.Itoa(id))

	// Previous article
	article.LastArticle, err = model.GetLastArticle(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// Next article
	article.NextArticle, err = model.GetNextArticle(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// Like count, view count
	article.ViewCount = int64(rdb.ZScore(rctx, g.ARTICLE_VIEW_COUNT, strconv.Itoa(id)).Val())
	likeCount, _ := strconv.Atoi(rdb.HGet(rctx, g.ARTICLE_LIKE_COUNT, strconv.Itoa(id)).Val())
	article.LikeCount = int64(likeCount)

	// Comment count
	article.CommentCount, err = model.GetArticleCommentCount(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, article)
}

// Get article archives
func (*Front) GetArchiveList(c *gin.Context) {
	var query FArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	list, total, err := model.GetBlogArticleList(GetDB(c), query.Page, query.Size, query.CategoryId, query.TagId)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	archives := make([]ArchiveVO, 0)
	for _, article := range list {
		archives = append(archives, ArchiveVO{
			ID:         article.ID,
			Title:      article.Title,
			Created_at: article.CreatedAt,
		})
	}

	ReturnSuccess(c, PageResult[ArchiveVO]{
		Total: total,
		List:  archives,
		Page:  query.Page,
		Size:  query.Size,
	})
}

// Like an article
// Record user liked a specific article to prevent duplicate likes
func (*Front) LikeArticle(c *gin.Context) {
	auth, _ := CurrentUserAuth(c)

	articleId, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rdb := GetRDB(c)

	// Record that a user liked a specific article
	articleLikeUserKey := g.ARTICLE_USER_LIKE_SET + strconv.Itoa(auth.ID)
	// If already recorded, like toggles to unlike
	if rdb.SIsMember(rctx, articleLikeUserKey, articleId).Val() {
		rdb.SRem(rctx, articleLikeUserKey, articleId)
		rdb.HIncrBy(rctx, g.ARTICLE_LIKE_COUNT, strconv.Itoa(articleId), -1)
	} else { // Not recorded, this is a new like
		rdb.SAdd(rctx, articleLikeUserKey, articleId)
		rdb.HIncrBy(rctx, g.ARTICLE_LIKE_COUNT, strconv.Itoa(articleId), 1)
	}

	ReturnSuccess(c, nil)
}

// Article search
func (*Front) SearchArticle(c *gin.Context) {
	result := make([]ArticleSearchVO, 0)

	keyword := c.Query("keyword")
	if keyword == "" {
		ReturnSuccess(c, result)
		return
	}

	db := GetDB(c)

	articleList, err := model.List(db, []model.Article{}, "*", "",
		"is_delete = 0 AND status = 1 AND (title LIKE ? OR content LIKE ?)",
		"%"+keyword+"%", "%"+keyword+"%")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	for _, article := range articleList {
		// Highlight keywords in title
		title := strings.ReplaceAll(article.Title, keyword,
			"<span style='color:#f47466'>"+keyword+"</span>")

		content := article.Content
		// Keyword start index in content
		keywordStartIndex := unicodeIndex(content, keyword)
		if keywordStartIndex != -1 { // Keyword exists in content
			preIndex, afterIndex := 0, 0
			if keywordStartIndex > 25 {
				preIndex = keywordStartIndex - 25
			}
			// Avoid Chinese substring garbling (use rune to handle multibyte)
			preText := substring(content, preIndex, keywordStartIndex)
			// string([]rune(content[preIndex:keywordStartIndex]))

			// Keyword end index in content
			keywordEndIndex := keywordStartIndex + unicodeLen(keyword)
			afterLength := len(content) - keywordEndIndex
			if afterLength > 175 {
				afterIndex = keywordEndIndex + 175
			} else {
				afterIndex = keywordEndIndex + afterLength
			}
			// afterText := string([]rune(content)[keywordStartIndex:afterIndex])
			afterText := substring(content, keywordStartIndex, afterIndex)
			// Highlight keywords in content
			content = strings.ReplaceAll(preText+afterText, keyword,
				"<span style='color:#f47466'>"+keyword+"</span>")
		}

		result = append(result, ArticleSearchVO{
			ID:      article.ID,
			Title:   title,
			Content: content,
		})
	}

	ReturnSuccess(c, result)
}

// Get substring index in Chinese string (rune-based, not bytes)
func unicodeIndex(str, substr string) int {
	// Byte index of substring in string
	result := strings.Index(str, substr)
	if result > 0 {
		prefix := []byte(str)[0:result]
		rs := []rune(string(prefix))
		result = len(rs)
	}
	return result
}

// Get real length of Chinese string (rune-based, not bytes)
func unicodeLen(str string) int {
	var r = []rune(str)
	return len(r)
}

// Fix incorrect position for Chinese characters
func substring(source string, start int, end int) string {
	var unicodeStr = []rune(source)
	length := len(unicodeStr)
	if start >= end {
		return ""
	}
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start <= 0 && end >= length {
		return source
	}
	var substring = ""
	for i := start; i < end; i++ {
		substring += string(unicodeStr[i])
	}
	return substring
}
