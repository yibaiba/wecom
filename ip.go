package wecom

import "context"

// APIDomainIPs returns WeCom API server IPs.
func (c *Client) APIDomainIPs(ctx context.Context) ([]string, error) {
	return c.ipList(ctx, "/cgi-bin/get_api_domain_ip")
}

// CallbackIPs returns WeCom callback server IPs.
func (c *Client) CallbackIPs(ctx context.Context) ([]string, error) {
	return c.ipList(ctx, "/cgi-bin/getcallbackip")
}

func (c *Client) ipList(ctx context.Context, path string) ([]string, error) {
	var out struct {
		IPList []string `json:"ip_list"`
	}
	if err := c.GetJSON(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return out.IPList, nil
}
