package client

import "github.com/cloudquery/cq-provider-sdk/provider/schema"

// MultiplexBy returns function which multiplex client into slice of clients by setting MultiplexedResourceId field
func MultiplexBy(resourcesGetter func(client *Client) []string) func(meta schema.ClientMeta) []schema.ClientMeta {
	return func(meta schema.ClientMeta) []schema.ClientMeta {
		var l = make([]schema.ClientMeta, 0)
		client := meta.(*Client)
		for _, id := range resourcesGetter(client) {
			l = append(l, client.withResource(id))
		}
		return l
	}
}

// Organizations returns organizations of client
func Organizations(client *Client) []string {
	return client.organizations
}

// Clouds returns organizations of client
func Clouds(client *Client) []string {
	return client.clouds
}

// Folders returns organizations of client
func Folders(client *Client) []string {
	return client.folders
}

// EmptyMultiplex returns slice with single client passed by param
func EmptyMultiplex(meta schema.ClientMeta) []schema.ClientMeta {
	return []schema.ClientMeta{meta}
}
