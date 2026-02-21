package domain

import (
	"time"

	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/utils"
)

// A Post represents well... a post.
// According to it's context (parentType + parentID) its can be
// used for profile posts, anime specific posts, etc...
type Post struct {
	ID string `json:"id" bson:"_id"`
	// context
	ParentId   string               `json:"parentId" bson:"parent_id"`
	ParentType value.PostParentType `json:"parentType" bson:"parent_type"`
	// content and metadata
	Text          *string   `json:"text,omitempty" bson:"text,omitempty"`
	TopMostPostId *string   `json:"topMostPostId,omitempty" bson:"top_most_post_id,omitempty"` // ID of the top most post in the thread, used for some validations
	Posts         []string  `json:"posts,omitempty" bson:"posts,omitempty"`                    // List of reply ids, used to easily fetch all replies to a post
	CreatedAt     time.Time `json:"createdAt" bson:"created_at"`
	CreatedBy     *int      `json:"createdBy,omitempty" bson:"created_by,omitempty"`
}

func NewPost(parentId string, parentType value.PostParentType, text string, createdBy int) *Post {
	var newPost Post
	newPost.ID = utils.GenerateRandomID()
	newPost.ParentId = parentId
	newPost.ParentType = parentType
	newPost.Text = &text
	newPost.TopMostPostId = nil
	newPost.CreatedAt = time.Now()
	newPost.CreatedBy = &createdBy
	return &newPost
}

func NewReply(replyTo *Post, text string, createdBy int) *Post {
	var newPost Post
	newPost.ID = utils.GenerateRandomID()
	newPost.ParentId = replyTo.ID
	newPost.ParentType = value.ParentTypePost // since it's a reply, the parent type is always Post
	newPost.Text = &text
	newPost.TopMostPostId = replyTo.TopMostPostId // the top most post id is inherited from the post being replied to
	newPost.CreatedAt = time.Now()
	newPost.CreatedBy = &createdBy
	return &newPost
}

func (p *Post) Delete() {
	p.Text = nil
	p.CreatedBy = nil
}

func (p *Post) IsReply() bool {
	return p.ParentType == value.ParentTypePost && p.TopMostPostId == nil
}

func (p *Post) IsDeleted() bool {
	return p.Text == nil && p.CreatedBy == nil
}
