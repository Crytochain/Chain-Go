// requestData project ssh-post.go
package requestData

type SSHPostParameters struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	Topics   []string `json:"topics"`
	Payload  string   `json:"payload"`
	Priority string   `json:"priority"`
	TTL      string   `json:"ttl"`
}
