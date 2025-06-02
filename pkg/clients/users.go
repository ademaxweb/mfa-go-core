package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ademaxweb/mfa-go-core/pkg/data"
	"net/http"
	"strings"
	"time"
)

type UsersClient interface {
	GetAllUsers() ([]data.User, error)
	GetUser(id int) (*data.User, error)
	GetUserByEmail(email string) (*data.User, error)
	CreateUser(user data.User) (int, error)
	UpdateUser(id int, user data.User) (*data.User, error)
	DeleteUser(id int) error
}

type httpUsersClient struct {
	baseURL string
	http    *http.Client
}

func NewUsersClient(baseURL string) (UsersClient, error) {
	c := &httpUsersClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		http: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    10 * time.Second,
				DisableCompression: false,
			},
		},
	}

	if err := c.HealthCheck(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *httpUsersClient) GetAllUsers() ([]data.User, error) {
	resp, err := c.http.Get(fmt.Sprintf("%s/users", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var users []data.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	return users, nil
}

func (c *httpUsersClient) GetUser(id int) (*data.User, error) {
	resp, err := c.http.Get(fmt.Sprintf("%s/users/%d", c.baseURL, id))
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user data.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	return &user, nil
}

func (c *httpUsersClient) GetUserByEmail(email string) (*data.User, error) {
	resp, err := c.http.Get(fmt.Sprintf("%s/users/%s", c.baseURL, email))
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user data.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	return &user, nil
}

func (c *httpUsersClient) CreateUser(user data.User) (int, error) {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return 0, fmt.Errorf("json marshal failed: %w", err)
	}

	resp, err := c.http.Post(
		fmt.Sprintf("%s/users", c.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return 0, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return 0, ErrInvalidData
		}
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var createdUser struct {
		Id int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
		return 0, fmt.Errorf("json decode failed: %w", err)
	}

	return createdUser.Id, nil
}

func (c *httpUsersClient) UpdateUser(id int, user data.User) (*data.User, error) {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("json marshal failed: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/users/%d", c.baseURL, id),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, ErrInvalidData
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var updatedUser data.User
	if err := json.NewDecoder(resp.Body).Decode(&updatedUser); err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	return &updatedUser, nil
}

func (c *httpUsersClient) DeleteUser(id int) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/users/%d", c.baseURL, id),
		nil,
	)
	if err != nil {
		return fmt.Errorf("request creation failed: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *httpUsersClient) HealthCheck() error {
	resp, err := c.http.Get(fmt.Sprintf("%s/health", c.baseURL))
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrServiceUnavailable
	}
	return nil
}
