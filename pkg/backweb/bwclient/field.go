package bwclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (c *client) SetUserFieldResponsibles(companyID int64, userID int64, fieldGUIDs []string) (int, error) {
	__url := fmt.Sprintf("/companies/%d/field_responsibles", companyID)

	fieldResponsible := struct {
		FieldResponsibles []struct {
			UserID     int64    `json:"user_id"`
			FieldGUIDs []string `json:"field_guids"`
		} `json:"field_responsibles"`
	}{
		FieldResponsibles: []struct {
			UserID     int64    `json:"user_id"`
			FieldGUIDs []string `json:"field_guids"`
		}{
			{
				UserID:     userID,
				FieldGUIDs: fieldGUIDs,
			},
		},
	}

	bs, _ := json.Marshal(fieldResponsible)

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
			return resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
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
	return resp.StatusCode, nil
}
