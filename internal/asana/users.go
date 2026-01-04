package asana

import (
	"context"
	"net/http"
)

func (c *Client) GetMe(ctx context.Context) (*User, error) {
	var user User
	if err := c.do(ctx, http.MethodGet, "/users/me", nil, nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
