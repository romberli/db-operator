package params

import "github.com/romberli/go-util/constant"

type Client struct {
	Socket   string `json:"socket" config:"socket"`
	User     string `json:"user" config:"user"`
	Password string `json:"password" config:"password"`
}

// NewClient returns a new *Client
func NewClient(socket, user, password string) *Client {
	return &Client{
		Socket:   socket,
		User:     user,
		Password: password,
	}
}

// NewClientWithDefault returns a new *Client with default values
func NewClientWithDefault() *Client {
	return &Client{
		Socket:   constant.DefaultRandomString,
		User:     constant.DefaultRandomString,
		Password: constant.DefaultRandomString,
	}
}

// GetSocket returns the socket
func (c *Client) GetSocket() string {
	return c.Socket
}

// GetUser returns the user
func (c *Client) GetUser() string {
	return c.User
}

// GetPassword returns the password
func (c *Client) GetPassword() string {
	return c.Password
}
