package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// TunnelRoute is the full record for a route.
type TunnelRoute struct {
	Network    string     `json:"network"`
	TunnelID   string     `json:"tunnel_id"`
	TunnelName string     `json:"tunnel_name"`
	Comment    string     `json:"comment"`
	CreatedAt  *time.Time `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type TunnelRoutesListParams struct {
	AccountID       string
	TunnelID        string
	Comment         string
	IsDeleted       *bool
	NetworkSubset   string
	NetworkSuperset string
	ExistedAt       *time.Time
	PaginationOptions
}

type TunnelRoutesCreateParams struct {
	AccountID string `json:"-"`
	Network   string `json:"-"`
	TunnelID  string `json:"tunnel_id"`
	Comment   string `json:"comment,omitempty"`
}

type TunnelRoutesUpdateParams struct {
	AccountID string `json:"-"`
	Network   string `json:"network"`
	TunnelID  string `json:"tunnel_id"`
	Comment   string `json:"comment,omitempty"`
}

type TunnelRoutesForIPParams struct {
	AccountID string `json:"-"`
	Network   string `json:"-"`
}

type TunnelRoutesDeleteParams struct {
	AccountID string `json:"-"`
	Network   string `json:"-"`
}

// tunnelRouteListResponse is the API response for listing tunnel routes.
type tunnelRouteListResponse struct {
	Response
	Result []TunnelRoute `json:"result"`
}

type tunnelRouteResponse struct {
	Response
	Result TunnelRoute `json:"result"`
}

// encode encodes non-nil fields into URL encoded form.
func (o TunnelRoutesListParams) encode() string {
	v := url.Values{}
	if o.TunnelID != "" {
		v.Set("tunnel_id", o.TunnelID)
	}
	if o.Comment != "" {
		v.Set("comment", o.Comment)
	}
	if o.IsDeleted != nil {
		v.Set("is_deleted", fmt.Sprintf("%t", *o.IsDeleted))
	}
	if o.NetworkSubset != "" {
		v.Set("network_subset", o.NetworkSubset)
	}
	if o.NetworkSuperset != "" {
		v.Set("network_superset", o.NetworkSuperset)
	}
	if o.ExistedAt != nil {
		v.Set("existed_at", (*o.ExistedAt).Format(time.RFC3339))
	}
	return v.Encode()
}

// ListTunnelRoutes lists all defined routes for tunnels in the account.
//
// See: https://api.cloudflare.com/#tunnel-route-list-tunnel-routes
func (api *API) ListTunnelRoutes(ctx context.Context, params TunnelRoutesListParams) ([]TunnelRoute, error) {
	if params.AccountID == "" {
		return []TunnelRoute{}, ErrMissingAccountID
	}

	uri := fmt.Sprintf("/%s/%s/teamnet/routes?%s", AccountRouteRoot, params.AccountID, params.encode())
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, params)

	if err != nil {
		return []TunnelRoute{}, err
	}

	var resp tunnelRouteListResponse
	err = json.Unmarshal(res, &resp)
	if err != nil {
		return []TunnelRoute{}, errors.Wrap(err, errUnmarshalError)
	}

	return resp.Result, nil
}

// GetTunnelRouteForIP finds the Tunnel Route that encompasses the given IP.
//
// See: https://api.cloudflare.com/#tunnel-route-get-tunnel-route-by-ip
func (api *API) GetTunnelRouteForIP(ctx context.Context, params TunnelRoutesForIPParams) (TunnelRoute, error) {
	uri := fmt.Sprintf("/%s/%s/teamnet/routes/ip/%s", AccountRouteRoot, params.AccountID, params.Network)

	responseBody, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return TunnelRoute{}, err
	}

	var routeResponse tunnelRouteResponse
	err = json.Unmarshal(responseBody, &routeResponse)
	if err != nil {
		return TunnelRoute{}, errors.Wrap(err, errUnmarshalError)
	}

	return routeResponse.Result, nil
}

// CreateTunnelRoute add a new route to the account routing table for the given
// tunnel.
//
// See: https://api.cloudflare.com/#tunnel-route-create-route
func (api *API) CreateTunnelRoute(ctx context.Context, params TunnelRoutesCreateParams) (TunnelRoute, error) {
	if params.AccountID == "" {
		return TunnelRoute{}, ErrMissingAccountID
	}

	uri := fmt.Sprintf("/%s/%s/teamnet/routes/network/%s", AccountRouteRoot, params.AccountID, url.PathEscape(params.Network))

	responseBody, err := api.makeRequestContext(ctx, http.MethodPost, uri, params)
	if err != nil {
		return TunnelRoute{}, err
	}

	var routeResponse tunnelRouteResponse
	err = json.Unmarshal(responseBody, &routeResponse)
	if err != nil {
		return TunnelRoute{}, errors.Wrap(err, errUnmarshalError)
	}

	return routeResponse.Result, nil
}

// DeleteTunnelRoute delete an existing route from the account routing table.
//
// See: https://api.cloudflare.com/#tunnel-route-delete-route
func (api *API) DeleteTunnelRoute(ctx context.Context, params TunnelRoutesDeleteParams) error {
	if params.AccountID == "" {
		return ErrMissingAccountID
	}

	uri := fmt.Sprintf("/%s/%s/teamnet/routes/network/%s", AccountRouteRoot, params.AccountID, url.PathEscape(params.Network))

	responseBody, err := api.makeRequestContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	var routeResponse tunnelRouteResponse
	err = json.Unmarshal(responseBody, &routeResponse)
	if err != nil {
		return errors.Wrap(err, errUnmarshalError)
	}

	return nil
}

// UpdateTunnelRoute updates an existing route in the account routing table for
// the given tunnel.
//
// See: https://api.cloudflare.com/#tunnel-route-update-route
func (api *API) UpdateTunnelRoute(ctx context.Context, params TunnelRoutesUpdateParams) (TunnelRoute, error) {
	if params.AccountID == "" {
		return TunnelRoute{}, ErrMissingAccountID
	}

	uri := fmt.Sprintf("/%s/%s/teamnet/routes/network/%s", AccountRouteRoot, params.AccountID, url.PathEscape(params.Network))

	responseBody, err := api.makeRequestContext(ctx, http.MethodPatch, uri, params)
	if err != nil {
		return TunnelRoute{}, err
	}

	var routeResponse tunnelRouteResponse
	err = json.Unmarshal(responseBody, &routeResponse)
	if err != nil {
		return TunnelRoute{}, errors.Wrap(err, errUnmarshalError)
	}

	return routeResponse.Result, nil
}
