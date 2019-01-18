// netServe project netServe.go
package netServe

type ProviderInterface interface {
	SendRequest(relust interface{}, method string, params interface{}, netNumber int) error
	Close() error
}
