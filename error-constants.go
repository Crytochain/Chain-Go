package customerror

import "errors"

var (
	// EMPTYRESPONSE - Server response is empty
	EMPTYRESPONSE = errors.New("Empty response")
	// UNPARSEABLEINTERFACE - the conversion failed
	UNPARSEABLEINTERFACE = errors.New("Unparseable Interface")
	// WEBSOCKETNOTDENIFIED - Websocket connection dont exist
	WEBSOCKETNOTDENIFIED = errors.New("Websocket connection dont exist")
)
