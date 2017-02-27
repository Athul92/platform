// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/gorilla/mux"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
	"net/http"
)

func InitReaction() {
	l4g.Debug(utils.T("api.reaction.init.debug"))

	BaseRoutes.NeedPost.Handle("/reactions/save", ApiUserRequired(saveReaction)).Methods("POST")
	BaseRoutes.NeedPost.Handle("/reactions/delete", ApiUserRequired(deleteReaction)).Methods("POST")
	BaseRoutes.NeedPost.Handle("/reactions", ApiUserRequired(listReactions)).Methods("GET")
}

func saveReaction(c *Context, w http.ResponseWriter, r *http.Request) {
	reaction := model.ReactionFromJson(r.Body)
	if reaction == nil {
		c.SetInvalidParam("saveReaction", "reaction")
		return
	}

	if reaction.UserId != c.Session.UserId {
		c.Err = model.NewLocAppError("saveReaction", "api.reaction.save_reaction.user_id.app_error", nil, "")
		c.Err.StatusCode = http.StatusForbidden
		return
	}

	params := mux.Vars(r)

	channelId := params["channel_id"]
	if len(channelId) != 26 {
		c.SetInvalidParam("saveReaction", "channelId")
		return
	}

	if !HasPermissionToChannelContext(c, channelId, model.PERMISSION_READ_CHANNEL) {
		return
	}

	postId := params["post_id"]
	if len(postId) != 26 || postId != reaction.PostId {
		c.SetInvalidParam("saveReaction", "postId")
		return
	}

	var post *model.Post

	if result := <-Srv.Store.Post().Get(reaction.PostId); result.Err != nil {
		c.Err = result.Err
		return
	} else if post = result.Data.(*model.PostList).Posts[postId]; post.ChannelId != channelId {
		c.Err = model.NewLocAppError("saveReaction", "api.reaction.save_reaction.mismatched_channel_id.app_error",
			nil, "channelId="+channelId+", post.ChannelId="+post.ChannelId+", postId="+postId)
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	if result := <-Srv.Store.Reaction().Save(reaction); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		go sendReactionEvent(model.WEBSOCKET_EVENT_REACTION_ADDED, channelId, reaction, post)
		reaction := result.Data.(*model.Reaction)
		w.Write([]byte(reaction.ToJson()))
	}
}

func deleteReaction(c *Context, w http.ResponseWriter, r *http.Request) {
	reaction := model.ReactionFromJson(r.Body)
	if reaction == nil {
		c.SetInvalidParam("deleteReaction", "reaction")
		return
	}

	if reaction.UserId != c.Session.UserId {
		c.Err = model.NewLocAppError("deleteReaction", "api.reaction.delete_reaction.user_id.app_error", nil, "")
		c.Err.StatusCode = http.StatusForbidden
		return
	}

	params := mux.Vars(r)

	channelId := params["channel_id"]
	if len(channelId) != 26 {
		c.SetInvalidParam("deleteReaction", "channelId")
		return
	}

	if !HasPermissionToChannelContext(c, channelId, model.PERMISSION_READ_CHANNEL) {
		return
	}

	postId := params["post_id"]
	if len(postId) != 26 || postId != reaction.PostId {
		c.SetInvalidParam("deleteReaction", "postId")
		return
	}

	var post *model.Post

	if result := <-Srv.Store.Post().Get(reaction.PostId); result.Err != nil {
		c.Err = result.Err
		return
	} else if post = result.Data.(*model.PostList).Posts[postId]; post.ChannelId != channelId {
		c.Err = model.NewLocAppError("deleteReaction", "api.reaction.delete_reaction.mismatched_channel_id.app_error",
			nil, "channelId="+channelId+", post.ChannelId="+post.ChannelId+", postId="+postId)
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	if result := <-Srv.Store.Reaction().Delete(reaction); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		go sendReactionEvent(model.WEBSOCKET_EVENT_REACTION_REMOVED, channelId, reaction, post)
		ReturnStatusOK(w)
	}
}

func sendReactionEvent(event string, channelId string, reaction *model.Reaction, post *model.Post) {
	// send out that a reaction has been added/removed

	message := model.NewWebSocketEvent(event, "", channelId, "", nil)
	message.Add("reaction", reaction.ToJson())
	Publish(message)

	// THe post is always modified since the UpdateAt always changes
	InvalidateCacheForChannelPosts(post.ChannelId)
	post.HasReactions = true
	post.UpdateAt = model.GetMillis()
	umessage := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_POST_EDITED, "", channelId, "", nil)
	umessage.Add("post", post.ToJson())
	Publish(umessage)

}

func listReactions(c *Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	channelId := params["channel_id"]
	if len(channelId) != 26 {
		c.SetInvalidParam("deletePost", "channelId")
		return
	}

	postId := params["post_id"]
	if len(postId) != 26 {
		c.SetInvalidParam("listReactions", "postId")
		return
	}

	pchan := Srv.Store.Post().Get(postId)

	if !HasPermissionToChannelContext(c, channelId, model.PERMISSION_READ_CHANNEL) {
		return
	}

	if result := <-pchan; result.Err != nil {
		c.Err = result.Err
		return
	} else if post := result.Data.(*model.PostList).Posts[postId]; post.ChannelId != channelId {
		c.Err = model.NewLocAppError("listReactions", "api.reaction.list_reactions.mismatched_channel_id.app_error",
			nil, "channelId="+channelId+", post.ChannelId="+post.ChannelId+", postId="+postId)
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	if result := <-Srv.Store.Reaction().GetForPost(postId); result.Err != nil {
		c.Err = result.Err
		return
	} else {
		reactions := result.Data.([]*model.Reaction)

		w.Write([]byte(model.ReactionsToJson(reactions)))
	}
}