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
	Text      string    `json:"text" bson:"text"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	CreatedBy int       `json:"createdBy" bson:"created_by"`
}

func NewPost(parentId string, parentType value.PostParentType, text string, createdBy int) *Post {
	return &Post{
		ID:         utils.GenerateRandomID(),
		ParentId:   parentId,
		ParentType: parentType,
		Text:       text,
		CreatedAt:  time.Now(),
		CreatedBy:  createdBy,
	}
}

func NewReply(replyTo *Post, text string, createdBy int) *Post {
	return &Post{
		ID:         utils.GenerateRandomID(),
		ParentId:   replyTo.ID,
		ParentType: value.ParentTypePost, // since it's a reply, the parent type is always Post
		Text:       text,
		CreatedAt:  time.Now(),
		CreatedBy:  createdBy,
	}
}
