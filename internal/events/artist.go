package events

import "hack.moontide.ink/pingfisher/muffled/internal/listenbrainz"

type Artist struct {
	Name string `json:"name,omitzero"`
	MBID string `json:"mbid,omitzero"`
}

func mapArtists(info listenbrainz.AdditionalInfo) []Artist {
	names := info.ArtistNames
	mbids := info.ArtistMBIDs

	n := min(len(names), len(mbids))

	artists := make([]Artist, 0, n)

	for i := range n {
		artists = append(artists, Artist{
			Name: names[i],
			MBID: mbids[i],
		})
	}

	return artists
}
