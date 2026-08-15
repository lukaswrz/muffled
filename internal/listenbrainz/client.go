package listenbrainz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
}

func (c *Client) endpoint(parts ...string) (*url.URL, error) {
	var encodedParts []string
	for _, part := range parts {
		encodedParts = append(encodedParts, url.PathEscape(part))
	}

	pathSuffix := strings.Join(encodedParts, "/")
	ref, err := url.Parse(pathSuffix)
	if err != nil {
		return nil, err
	}

	// FIXME: Use https://pkg.go.dev/net/new@master#URL.Clone
	new, err := url.Parse(c.BaseURL.String())
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(new.Path, "/") {
		new.Path += "/"
	}

	return new.ResolveReference(ref), nil
}

func NewClient(base string) (*Client, error) {
	p, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(p.Path, "/") {
		p.Path += "/"
	}

	return &Client{
		BaseURL:    p,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

type PlayingNowResponse struct {
	Payload struct {
		Count      int      `json:"count"`
		Listens    []Listen `json:"listens"`
		PlayingNow bool     `json:"playing_now"`
		UserID     string   `json:"user_id"`
	} `json:"payload"`
}

type Listen struct {
	PlayingNow    bool          `json:"playing_now"`
	TrackMetadata TrackMetadata `json:"track_metadata"`
}

type TrackMetadata struct {
	AdditionalInfo AdditionalInfo `json:"additional_info"`
	ArtistName     string         `json:"artist_name"`
	ReleaseName    string         `json:"release_name"`
	TrackName      string         `json:"track_name"`
}

type AdditionalInfo struct {
	ArtistMBIDs             []string `json:"artist_mbids"`
	ArtistNames             []string `json:"artist_names"`
	DurationMS              int      `json:"duration_ms"`
	RecordingMBID           string   `json:"recording_mbid"`
	ReleaseGroupMBID        string   `json:"release_group_mbid"`
	ReleaseMBID             string   `json:"release_mbid"`
	SubmissionClient        string   `json:"submission_client"`
	SubmissionClientVersion string   `json:"submission_client_version"`
	TrackNumber             int      `json:"tracknumber"`
}

func (c *Client) GetPlayingNow(username string) (PlayingNowResponse, error) {
	url, err := c.endpoint("user", username, "playing-now")
	if err != nil {
		return PlayingNowResponse{}, err
	}

	raw, err := c.HTTPClient.Get(url.String())
	if err != nil {
		return PlayingNowResponse{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		cerr := raw.Body.Close()
		if cerr != nil && err == nil {
			err = cerr
		}
	}()

	if raw.StatusCode != http.StatusOK {
		return PlayingNowResponse{}, fmt.Errorf("unexpected status code after making a GET request to %q: %d", url.String(), raw.StatusCode)
	}

	response := PlayingNowResponse{}
	if err := json.NewDecoder(raw.Body).Decode(&response); err != nil {
		return PlayingNowResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}
