package osd

type Applies func(client *Client)

func UsePrefix(prefix string) func(*Client) {
	return func(c *Client) {
		c.prefix = prefix
	}
}
func UseDuration(duration int64) func(*Client) {
	return func(c *Client) {
		c.duration = duration
	}
}

func UseCallback(callback string) func(*Client) {
	return func(c *Client) {
		c.callback = callback
	}
}
