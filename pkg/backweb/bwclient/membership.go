package bwclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (c client) SetUserMemberships(companyID int64, userID int64, structureIDs []int64) (int, error) {
	__url := fmt.Sprintf("/companies/%d/memberships", companyID)

	ms := struct {
		Memberships []struct {
			UserID     int64   `json:"user_id"`
			Structures []int64 `json:"structures"`
		} `json:"memberships"`
	}{
		Memberships: []struct {
			UserID     int64   `json:"user_id"`
			Structures []int64 `json:"structures"`
		}{
			{
				UserID:     userID,
				Structures: structureIDs,
			},
		},
	}

	bs, _ := json.Marshal(ms)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s%s", c.url, __url), bytes.NewBuffer(bs))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	xServiceToken := os.Getenv("SERVICE_SECRET_KEY")
	if xServiceToken == "" {
		return http.StatusBadRequest, fmt.Errorf("SERVICE_SECRET_KEY is empty")
	}

	req.Header.Set("X-Service-Token", xServiceToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		resBS, err := io.ReadAll(resp.Body)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		switch resp.StatusCode {
		case http.StatusForbidden:
			return http.StatusForbidden, fmt.Errorf(string(resBS))
		case http.StatusUnprocessableEntity:
			return http.StatusUnprocessableEntity, fmt.Errorf(string(resBS))
		default:
			return http.StatusInternalServerError, fmt.Errorf("error from %s\n%d: %s\n%s", c.url, resp.StatusCode, string(resBS), string(bs))
		}
	}
	return http.StatusOK, err
}
