// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package store

import (
	"time"

	l4g "github.com/alecthomas/log4go"
	"github.com/mattermost/platform/model"
)

type StoreResult struct {
	Data interface{}
	Err  *model.AppError
}

type StoreChannel chan StoreResult

func Must(sc StoreChannel) interface{} {
	r := <-sc
	if r.Err != nil {
		l4g.Close()
		time.Sleep(time.Second)
		panic(r.Err)
	}

	return r.Data
}

type Store interface {
	Team() TeamStore
	Channel() ChannelStore
	Post() PostStore
	User() UserStore
	Audit() AuditStore
	Compliance() ComplianceStore
	Session() SessionStore
	OAuth() OAuthStore
	System() SystemStore
	Webhook() WebhookStore
	Command() CommandStore
	Preference() PreferenceStore
	License() LicenseStore
	PasswordRecovery() PasswordRecoveryStore
	Emoji() EmojiStore
	Status() StatusStore
	FileInfo() FileInfoStore
	Reaction() ReactionStore
	MarkSystemRanUnitTests()
	Close()
	DropAllTables()
	TotalMasterDbConnections() int
	TotalReadDbConnections() int
}

type TeamStore interface {
	Save(team *model.Team) StoreChannel
	Update(team *model.Team) StoreChannel
	UpdateDisplayName(name string, teamId string) StoreChannel
	Get(id string) StoreChannel
	GetByName(name string) StoreChannel
	SearchByName(name string) StoreChannel
	GetAll() StoreChannel
	GetAllPage(offset int, limit int) StoreChannel
	GetAllTeamListing() StoreChannel
	GetAllTeamPageListing(offset int, limit int) StoreChannel
	GetTeamsByUserId(userId string) StoreChannel
	GetByInviteId(inviteId string) StoreChannel
	PermanentDelete(teamId string) StoreChannel
	AnalyticsTeamCount() StoreChannel
	SaveMember(member *model.TeamMember) StoreChannel
	UpdateMember(member *model.TeamMember) StoreChannel
	GetMember(teamId string, userId string) StoreChannel
	GetMembers(teamId string, offset int, limit int) StoreChannel
	GetMembersByIds(teamId string, userIds []string) StoreChannel
	GetTotalMemberCount(teamId string) StoreChannel
	GetActiveMemberCount(teamId string) StoreChannel
	GetTeamsForUser(userId string) StoreChannel
	GetChannelUnreadsForAllTeams(excludeTeamId, userId string) StoreChannel
	GetChannelUnreadsForTeam(teamId, userId string) StoreChannel
	RemoveMember(teamId string, userId string) StoreChannel
	RemoveAllMembersByTeam(teamId string) StoreChannel
	RemoveAllMembersByUser(userId string) StoreChannel
}

type ChannelStore interface {
	Save(channel *model.Channel) StoreChannel
	CreateDirectChannel(userId string, otherUserId string) StoreChannel
	SaveDirectChannel(channel *model.Channel, member1 *model.ChannelMember, member2 *model.ChannelMember) StoreChannel
	Update(channel *model.Channel) StoreChannel
	Get(id string, allowFromCache bool) StoreChannel
	InvalidateChannel(id string)
	InvalidateChannelByName(teamId, name string)
	GetFromMaster(id string) StoreChannel
	Delete(channelId string, time int64) StoreChannel
	SetDeleteAt(channelId string, deleteAt int64, updateAt int64) StoreChannel
	PermanentDeleteByTeam(teamId string) StoreChannel
	PermanentDelete(channelId string) StoreChannel
	GetByName(team_id string, name string, allowFromCache bool) StoreChannel
	GetByNameIncludeDeleted(team_id string, name string, allowFromCache bool) StoreChannel
	GetDeletedByName(team_id string, name string) StoreChannel
	GetChannels(teamId string, userId string) StoreChannel
	GetMoreChannels(teamId string, userId string, offset int, limit int) StoreChannel
	GetPublicChannelsForTeam(teamId string, offset int, limit int) StoreChannel
	GetChannelCounts(teamId string, userId string) StoreChannel
	GetTeamChannels(teamId string) StoreChannel
	GetAll(teamId string) StoreChannel
	GetForPost(postId string) StoreChannel
	SaveMember(member *model.ChannelMember) StoreChannel
	UpdateMember(member *model.ChannelMember) StoreChannel
	GetMembers(channelId string, offset, limit int) StoreChannel
	GetMember(channelId string, userId string) StoreChannel
	GetAllChannelMembersForUser(userId string, allowFromCache bool) StoreChannel
	InvalidateAllChannelMembersForUser(userId string)
	IsUserInChannelUseCache(userId string, channelId string) bool
	GetAllChannelMembersNotifyPropsForChannel(channelId string, allowFromCache bool) StoreChannel
	InvalidateCacheForChannelMembersNotifyProps(channelId string)
	GetMemberForPost(postId string, userId string) StoreChannel
	InvalidateMemberCount(channelId string)
	GetMemberCountFromCache(channelId string) int64
	GetMemberCount(channelId string, allowFromCache bool) StoreChannel
	GetPinnedPosts(channelId string) StoreChannel
	RemoveMember(channelId string, userId string) StoreChannel
	PermanentDeleteMembersByUser(userId string) StoreChannel
	PermanentDeleteMembersByChannel(channelId string) StoreChannel
	UpdateLastViewedAt(channelIds []string, userId string) StoreChannel
	SetLastViewedAt(channelId string, userId string, newLastViewedAt int64) StoreChannel
	IncrementMentionCount(channelId string, userId string) StoreChannel
	AnalyticsTypeCount(teamId string, channelType string) StoreChannel
	ExtraUpdateByUser(userId string, time int64) StoreChannel
	GetMembersForUser(teamId string, userId string) StoreChannel
	SearchInTeam(teamId string, term string) StoreChannel
	SearchMore(userId string, teamId string, term string) StoreChannel
	GetMembersByIds(channelId string, userIds []string) StoreChannel
	AnalyticsDeletedTypeCount(teamId string, channelType string) StoreChannel
}

type PostStore interface {
	Save(post *model.Post) StoreChannel
	Update(newPost *model.Post, oldPost *model.Post) StoreChannel
	Get(id string) StoreChannel
	GetSingle(id string) StoreChannel
	Delete(postId string, time int64) StoreChannel
	PermanentDeleteByUser(userId string) StoreChannel
	PermanentDeleteByChannel(channelId string) StoreChannel
	GetPosts(channelId string, offset int, limit int, allowFromCache bool) StoreChannel
	GetFlaggedPosts(userId string, offset int, limit int) StoreChannel
	GetPostsBefore(channelId string, postId string, numPosts int, offset int) StoreChannel
	GetPostsAfter(channelId string, postId string, numPosts int, offset int) StoreChannel
	GetPostsSince(channelId string, time int64, allowFromCache bool) StoreChannel
	GetEtag(channelId string, allowFromCache bool) StoreChannel
	Search(teamId string, userId string, params *model.SearchParams) StoreChannel
	AnalyticsUserCountsWithPostsByDay(teamId string) StoreChannel
	AnalyticsPostCountsByDay(teamId string) StoreChannel
	AnalyticsPostCount(teamId string, mustHaveFile bool, mustHaveHashtag bool) StoreChannel
	InvalidateLastPostTimeCache(channelId string)
	GetPostsCreatedAt(channelId string, time int64) StoreChannel
	Overwrite(post *model.Post) StoreChannel
}

type UserStore interface {
	Save(user *model.User) StoreChannel
	Update(user *model.User, allowRoleUpdate bool) StoreChannel
	UpdateLastPictureUpdate(userId string) StoreChannel
	UpdateUpdateAt(userId string) StoreChannel
	UpdatePassword(userId, newPassword string) StoreChannel
	UpdateAuthData(userId string, service string, authData *string, email string, resetMfa bool) StoreChannel
	UpdateMfaSecret(userId, secret string) StoreChannel
	UpdateMfaActive(userId string, active bool) StoreChannel
	Get(id string) StoreChannel
	GetAll() StoreChannel
	InvalidateProfilesInChannelCacheByUser(userId string)
	InvalidateProfilesInChannelCache(channelId string)
	GetProfilesInChannel(channelId string, offset int, limit int) StoreChannel
	GetAllProfilesInChannel(channelId string, allowFromCache bool) StoreChannel
	GetProfilesNotInChannel(teamId string, channelId string, offset int, limit int) StoreChannel
	GetProfilesByUsernames(usernames []string, teamId string) StoreChannel
	GetAllProfiles(offset int, limit int) StoreChannel
	GetProfiles(teamId string, offset int, limit int) StoreChannel
	GetProfileByIds(userId []string, allowFromCache bool) StoreChannel
	InvalidatProfileCacheForUser(userId string)
	GetByEmail(email string) StoreChannel
	GetByAuth(authData *string, authService string) StoreChannel
	GetAllUsingAuthService(authService string) StoreChannel
	GetByUsername(username string) StoreChannel
	GetForLogin(loginId string, allowSignInWithUsername, allowSignInWithEmail, ldapEnabled bool) StoreChannel
	VerifyEmail(userId string) StoreChannel
	GetEtagForAllProfiles() StoreChannel
	GetEtagForProfiles(teamId string) StoreChannel
	UpdateFailedPasswordAttempts(userId string, attempts int) StoreChannel
	GetTotalUsersCount() StoreChannel
	GetSystemAdminProfiles() StoreChannel
	PermanentDelete(userId string) StoreChannel
	AnalyticsUniqueUserCount(teamId string) StoreChannel
	AnalyticsActiveCount(time int64) StoreChannel
	GetUnreadCount(userId string) StoreChannel
	GetUnreadCountForChannel(userId string, channelId string) StoreChannel
	GetRecentlyActiveUsersForTeam(teamId string) StoreChannel
	Search(teamId string, term string, options map[string]bool) StoreChannel
	SearchInChannel(channelId string, term string, options map[string]bool) StoreChannel
	SearchNotInChannel(teamId string, channelId string, term string, options map[string]bool) StoreChannel
	AnalyticsGetInactiveUsersCount() StoreChannel
	AnalyticsGetSystemAdminCount() StoreChannel
}

type SessionStore interface {
	Save(session *model.Session) StoreChannel
	Get(sessionIdOrToken string) StoreChannel
	GetSessions(userId string) StoreChannel
	GetSessionsWithActiveDeviceIds(userId string) StoreChannel
	Remove(sessionIdOrToken string) StoreChannel
	RemoveAllSessions() StoreChannel
	PermanentDeleteSessionsByUser(teamId string) StoreChannel
	UpdateLastActivityAt(sessionId string, time int64) StoreChannel
	UpdateRoles(userId string, roles string) StoreChannel
	UpdateDeviceId(id string, deviceId string, expiresAt int64) StoreChannel
	AnalyticsSessionCount() StoreChannel
}

type AuditStore interface {
	Save(audit *model.Audit) StoreChannel
	Get(user_id string, offset int, limit int) StoreChannel
	PermanentDeleteByUser(userId string) StoreChannel
}

type ComplianceStore interface {
	Save(compliance *model.Compliance) StoreChannel
	Update(compliance *model.Compliance) StoreChannel
	Get(id string) StoreChannel
	GetAll(offset, limit int) StoreChannel
	ComplianceExport(compliance *model.Compliance) StoreChannel
}

type OAuthStore interface {
	SaveApp(app *model.OAuthApp) StoreChannel
	UpdateApp(app *model.OAuthApp) StoreChannel
	GetApp(id string) StoreChannel
	GetAppByUser(userId string) StoreChannel
	GetApps() StoreChannel
	GetAuthorizedApps(userId string) StoreChannel
	DeleteApp(id string) StoreChannel
	SaveAuthData(authData *model.AuthData) StoreChannel
	GetAuthData(code string) StoreChannel
	RemoveAuthData(code string) StoreChannel
	PermanentDeleteAuthDataByUser(userId string) StoreChannel
	SaveAccessData(accessData *model.AccessData) StoreChannel
	UpdateAccessData(accessData *model.AccessData) StoreChannel
	GetAccessData(token string) StoreChannel
	GetAccessDataByUserForApp(userId, clientId string) StoreChannel
	GetAccessDataByRefreshToken(token string) StoreChannel
	GetPreviousAccessData(userId, clientId string) StoreChannel
	RemoveAccessData(token string) StoreChannel
}

type SystemStore interface {
	Save(system *model.System) StoreChannel
	SaveOrUpdate(system *model.System) StoreChannel
	Update(system *model.System) StoreChannel
	Get() StoreChannel
	GetByName(name string) StoreChannel
}

type WebhookStore interface {
	SaveIncoming(webhook *model.IncomingWebhook) StoreChannel
	GetIncoming(id string, allowFromCache bool) StoreChannel
	GetIncomingList(offset, limit int) StoreChannel
	GetIncomingByTeam(teamId string, offset, limit int) StoreChannel
	UpdateIncoming(webhook *model.IncomingWebhook) StoreChannel
	GetIncomingByChannel(channelId string) StoreChannel
	DeleteIncoming(webhookId string, time int64) StoreChannel
	PermanentDeleteIncomingByUser(userId string) StoreChannel

	SaveOutgoing(webhook *model.OutgoingWebhook) StoreChannel
	GetOutgoing(id string) StoreChannel
	GetOutgoingList(offset, limit int) StoreChannel
	GetOutgoingByChannel(channelId string, offset, limit int) StoreChannel
	GetOutgoingByTeam(teamId string, offset, limit int) StoreChannel
	DeleteOutgoing(webhookId string, time int64) StoreChannel
	PermanentDeleteOutgoingByUser(userId string) StoreChannel
	UpdateOutgoing(hook *model.OutgoingWebhook) StoreChannel

	AnalyticsIncomingCount(teamId string) StoreChannel
	AnalyticsOutgoingCount(teamId string) StoreChannel
	InvalidateWebhookCache(webhook string)
}

type CommandStore interface {
	Save(webhook *model.Command) StoreChannel
	Get(id string) StoreChannel
	GetByTeam(teamId string) StoreChannel
	Delete(commandId string, time int64) StoreChannel
	PermanentDeleteByUser(userId string) StoreChannel
	Update(hook *model.Command) StoreChannel
	AnalyticsCommandCount(teamId string) StoreChannel
}

type PreferenceStore interface {
	Save(preferences *model.Preferences) StoreChannel
	Get(userId string, category string, name string) StoreChannel
	GetCategory(userId string, category string) StoreChannel
	GetAll(userId string) StoreChannel
	Delete(userId, category, name string) StoreChannel
	DeleteCategory(userId string, category string) StoreChannel
	DeleteCategoryAndName(category string, name string) StoreChannel
	PermanentDeleteByUser(userId string) StoreChannel
	IsFeatureEnabled(feature, userId string) StoreChannel
}

type LicenseStore interface {
	Save(license *model.LicenseRecord) StoreChannel
	Get(id string) StoreChannel
}

type PasswordRecoveryStore interface {
	SaveOrUpdate(recovery *model.PasswordRecovery) StoreChannel
	Delete(userId string) StoreChannel
	Get(userId string) StoreChannel
	GetByCode(code string) StoreChannel
}

type EmojiStore interface {
	Save(emoji *model.Emoji) StoreChannel
	Get(id string, allowFromCache bool) StoreChannel
	GetByName(name string) StoreChannel
	GetAll() StoreChannel
	Delete(id string, time int64) StoreChannel
}

type StatusStore interface {
	SaveOrUpdate(status *model.Status) StoreChannel
	Get(userId string) StoreChannel
	GetByIds(userIds []string) StoreChannel
	GetOnlineAway() StoreChannel
	GetOnline() StoreChannel
	GetAllFromTeam(teamId string) StoreChannel
	ResetAll() StoreChannel
	GetTotalActiveUsersCount() StoreChannel
	UpdateLastActivityAt(userId string, lastActivityAt int64) StoreChannel
}

type FileInfoStore interface {
	Save(info *model.FileInfo) StoreChannel
	Get(id string) StoreChannel
	GetByPath(path string) StoreChannel
	GetForPost(postId string, readFromMaster bool, allowFromCache bool) StoreChannel
	InvalidateFileInfosForPostCache(postId string)
	AttachToPost(fileId string, postId string) StoreChannel
	DeleteForPost(postId string) StoreChannel
}

type ReactionStore interface {
	Save(reaction *model.Reaction) StoreChannel
	Delete(reaction *model.Reaction) StoreChannel
	InvalidateCacheForPost(postId string)
	InvalidateCache()
	GetForPost(postId string, allowFromCache bool) StoreChannel
	DeleteAllWithEmojiName(emojiName string) StoreChannel
}
