// Package strava provides the API client for interfacing with Stava data
package strava

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/zwinslett/speed-daemon/model"
)

type Client struct {
	httpClient   *http.Client
	clientID     string
	clientSecret string
	refreshToken string
	accessToken  string
	expiresAt    int64
	baseURL      string
	authURL      string
}

func NewClient() *Client {
	return &Client{
		httpClient:   &http.Client{},
		clientID:     os.Getenv("STRAVA_CLIENT_ID"),
		clientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
		refreshToken: os.Getenv("STRAVA_REFRESH_TOKEN"),
		baseURL:      "https://www.strava.com/api/v3",
		authURL:      "https://www.strava.com/oauth/token",
	}
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func (c *Client) SetAccessToken(ctx context.Context) error {
	data := url.Values{}
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("refresh_token", c.refreshToken)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.authURL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body bytes.Buffer
		_, _ = body.ReadFrom(resp.Body)
		return fmt.Errorf("failed: %d body=%s", resp.StatusCode, body.String())
	}

	var tr tokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tr)
	if err != nil {
		return err
	}
	c.accessToken = tr.AccessToken
	c.expiresAt = tr.ExpiresAt
	return nil
}

func (c *Client) RefreshAccessToken(ctx context.Context) error {
	if time.Now().Unix() >= c.expiresAt {
		return c.SetAccessToken(ctx)
	}
	return nil
}

func (c *Client) doGet(ctx context.Context, url string, output any) error {
	err := c.RefreshAccessToken(ctx)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body bytes.Buffer
		_, _ = body.ReadFrom(resp.Body)
		return fmt.Errorf("failed: %d body=%s", resp.StatusCode, body.String())
	}
	return json.NewDecoder(resp.Body).Decode(output)
}

func (c *Client) GetActivityByID(ctx context.Context, id int64) (model.DetailedActivity, error) {
	url := fmt.Sprintf("%s/activities/%d", c.baseURL, id)

	var activity model.DetailedActivity
	err := c.doGet(ctx, url, &activity)
	if err != nil {
		return model.DetailedActivity{}, err
	}
	return activity, nil
}

func (c *Client) GetActivitiesByRange(ctx context.Context, after int64, before int64) ([]model.Activity, error) {
	rawURL := fmt.Sprintf("%s/athlete/activities", c.baseURL)
	params := url.Values{}
	params.Set("after", fmt.Sprintf("%d", after))
	params.Set("before", fmt.Sprintf("%d", before))
	params.Set("per_page", "200")
	fullURL := rawURL + "?" + params.Encode()

	var activities []model.Activity
	err := c.doGet(ctx, fullURL, &activities)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (c *Client) GetRecentActivities(ctx context.Context, perPage int) ([]model.Activity, error) {
	rawUrl := fmt.Sprintf("%s/athlete/activities", c.baseURL)
	params := url.Values{}
	params.Set("per_page", fmt.Sprintf("%d", perPage))
	fullUrl := rawUrl + "?" + params.Encode()

	var activities []model.Activity
	err := c.doGet(ctx, fullUrl, &activities)
	if err != nil {
		return nil, err
	}

	return activities, nil
}

func (c *Client) GetActivityZones(ctx context.Context, id int64) ([]model.Zones, error) {
	url := fmt.Sprintf("%s/activities/%d/zones", c.baseURL, id)
	var zones []model.Zones
	err := c.doGet(ctx, url, &zones)
	if err != nil {
		return nil, err
	}
	return zones, nil
}
