bitly API golang library
========================

## Run tests

Your username is the lowercase name shown when you login to bitly, your access token can be fetched using the following (http://dev.bitly.com/authentication.html ):

    curl -u "username:password" -X POST "https://api-ssl.bitly.com/oauth/access_token"

To run the tests either export the environment variable or set it up inline before calling `nosetests`:

    BITLY_ACCESS_TOKEN=<accesstoken> go test

## API Documentation

http://dev.bitly.com/

Methods:
Shorten
Expand
Clicks
ClicksByDay
ClicksByMinute
Referrers
Info
LinkEncodersCount
UserLinkLookup
UserLinkEdit
UserLinkSave
UserLinkHistory
LinkClicks
LinkReferrersByDomain
LinkReferrers
LinkShares
LinkCountries
LinkInfo
LinkContent
LinkCategory
LinkSocial
LinkLocation

To Implement:
UserClicks
UserCountries
UserPopularLinks
UserReferrers
UserReferringDomains
UserShareCounts
UserShareCountsByShareType
UserShotenCounts
UserTrackingDomainList
UserTrakcingDomainClicks
UserTrackingDomainShortenCounts
UserInfo
UserNetworkHistory
UserBundleHistory

ProDomain

BundleArchive
BundleBundlesByUser
BundleClone
BundleCollaboratorAdd
BundleCollaboratorRemove
BundleContents
BundleCreate
BundleEdit
BundleLinkAdd
BundleLinkCommentAdd
BundleLinkCommentEdit
BundleLinkCommentRemove
BundleLinkEdit
BundleLinkRemove
BundleLinkReorder
BundleViewCount

RealtimeBurstingPhrases
RealtimeHotPhrases
RealtimeClickrate
HighValue
Search
