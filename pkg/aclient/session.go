package aclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
)

// Type session inherits client's data and add user token
type Session struct {
	client *client
	token  string
}

// Get user data from server by user token
func (s *Session) GetUser() (model.User, error) {
	user := model.User{}
	getUserURL, _ := url.JoinPath(s.client.url, "current_user")
	req, err := http.NewRequest(http.MethodGet, getUserURL, nil)
	if err != nil {
		return user, err
	}
	req.Header.Set("X-Token", s.token)
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return model.User{}, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.User{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return model.User{}, fmt.Errorf("error getting user data: %s", string(bs))
	}

	err = json.Unmarshal(bs, &user)
	if err != nil {
		return user, fmt.Errorf("error unmarshalling user data: %s", err)
	}
	return user, err
}

// Returns company positions
func (s *Session) GetPositions(companyID int64) (map[int64]string, error) {

	getCompanyPositionsURL, _ := url.JoinPath(s.client.url, fmt.Sprintf("companies/%d/positions", companyID))
	req, err := http.NewRequest(http.MethodGet, getCompanyPositionsURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Token", s.token)
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting company positions: %s", string(bs))
	}

	var positions []model.CompanyPosition
	err = json.Unmarshal(bs, &positions)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling company positions: %s", err)
	}

	res := make(map[int64]string)
	for _, p := range positions {
		res[p.ID] = p.Name
	}

	return res, nil
}
