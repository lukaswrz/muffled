package events

import (
	"fmt"

	"hack.moontide.ink/pingfisher/muffled/internal/listenbrainz"
)

type PlayingNowEvent struct {
	Playing                 bool     `json:"playing"`
	UserID                  string   `json:"userId,omitzero"`
	Artists                 []Artist `json:"artists,omitzero"`
	Artist                  string   `json:"artist,omitzero"`
	Title                   string   `json:"title,omitzero"`
	Release                 string   `json:"release,omitzero"`
	TrackNumber             int      `json:"trackNumber,omitzero"`
	DurationMS              int      `json:"duration,omitzero"`
	RecordingMBID           string   `json:"recordingMbid,omitzero"`
	ReleaseMBID             string   `json:"releaseMbid,omitzero"`
	ReleaseGroupMBID        string   `json:"releaseGroupMbid,omitzero"`
	SubmissionClient        string   `json:"submissionClient,omitzero"`
	SubmissionClientVersion string   `json:"submissionClientVersion,omitzero"`
}

func MapPlayingNowEvent(response listenbrainz.PlayingNowResponse) (event PlayingNowEvent, err error) {
	payload := response.Payload

	nlistens := len(payload.Listens)

	if nlistens == 0 {
		return PlayingNowEvent{
			Playing: false,
		}, nil
	}

	if nlistens != 1 {
		return PlayingNowEvent{}, fmt.Errorf(
			"expected at most 1 listen, got %d",
			nlistens,
		)
	}

	listen := payload.Listens[0]
	meta := listen.TrackMetadata
	info := meta.AdditionalInfo

	return PlayingNowEvent{
		Playing: payload.PlayingNow,
		UserID:  payload.UserID,

		Artists:     mapArtists(info),
		Artist:      meta.ArtistName,
		Title:       meta.TrackName,
		Release:     meta.ReleaseName,
		TrackNumber: info.TrackNumber,
		DurationMS:  info.DurationMS,

		RecordingMBID:    info.RecordingMBID,
		ReleaseMBID:      info.ReleaseMBID,
		ReleaseGroupMBID: info.ReleaseGroupMBID,

		SubmissionClient:        info.SubmissionClient,
		SubmissionClientVersion: info.SubmissionClientVersion,
	}, nil
}
