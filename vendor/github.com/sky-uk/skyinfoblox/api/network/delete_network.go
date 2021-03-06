package network

import (
	"fmt"
	"github.com/sky-uk/skyinfoblox/api"
	"net/http"
)

// DeleteNetAPI base object.
type DeleteNetAPI struct {
	*api.BaseAPI
}

// NewDeleteNetwork returns a new object of type DeleteNetworkAPI.
func NewDeleteNetwork(objRef string) *DeleteNetAPI {
	this := new(DeleteNetAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, fmt.Sprintf("%s/%s", wapiVersion, objRef), nil, new(string))
	return this
}

// GetResponse casts the response object and
// returns the single network object
func (gn DeleteNetAPI) GetResponse() string {
	return *gn.ResponseObject().(*string)
}
