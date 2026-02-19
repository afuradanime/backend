package domain

import (
	"time"

	"github.com/afuradanime/backend/internal/core/utils"
)

type ContextType string

const (
	ContextTypeProfile ContextType = "Profile"
	ContextTypeAnimeOp ContextType = "AnimeOpinion"
	ContextTypeForum   ContextType = "Forum"
)

// Thread "holder" that knows the context that the thread is related to
type ThreadContext struct {
	// for example, if context id is a userid, context type is "user"
	// and the thread is the user profile thread
	ID        string `bson:"_id"`
	ContextId int    `bson:"contextId"`
	// The owner of a thread context is an admin unless contextType is "user",
	// then the owner is the user itself, and only the owner can pin/unpin thread posts
	ContextType        ContextType `bson:"contextType"`
	PinnedThreadPostID *string     `bson:"pinnedThreadPostId,omitempty"`
}

// the actual thread post, that is related to a thread context
type ThreadPost struct {
	ID        string  `bson:"_id"`
	ContextID int     `bson:"contextId"`
	UserId    int     `bson:"userId"`
	Content   string  `bson:"content"`
	CreatedAt int64   `bson:"createdAt"`
	ReplyTo   *string `bson:"replyTo,omitempty"`
}

func NewContext(contextId int, contextType ContextType) *ThreadContext {
	return &ThreadContext{
		ID:          utils.GenerateRandomID(),
		ContextId:   contextId,
		ContextType: contextType,
	}
}

func NewThreadPost(context int, userId int, content string) *ThreadPost {
	return &ThreadPost{
		ID:        utils.GenerateRandomID(),
		ContextID: context,
		UserId:    userId,
		Content:   content,
		CreatedAt: time.Now().Unix(),
	}
}

func (t *ThreadPost) AddCreationTime(createdAt int64) {
	t.CreatedAt = createdAt
}

func (t *ThreadPost) ReplyToPost(replyTo string) {
	t.ReplyTo = &replyTo
}

func (t *ThreadPost) IsReply() bool {
	return t.ReplyTo != nil
}

func (t *ThreadPost) IsPinned(context *ThreadContext) bool {
	return context.PinnedThreadPostID != nil && *context.PinnedThreadPostID == t.ID
}

func (t *ThreadPost) Pin(context *ThreadContext) {
	context.PinnedThreadPostID = &t.ID
}

func (t *ThreadPost) Unpin(context *ThreadContext) {
	if context.PinnedThreadPostID != nil && *context.PinnedThreadPostID == t.ID {
		context.PinnedThreadPostID = nil
	}
}
