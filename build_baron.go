package iframe

import (
	"fmt"
	"github.com/10gen-labs/slogger/v1"
	"github.com/evergreen-ci/evergreen"
	"github.com/evergreen-ci/evergreen/db"
	"github.com/evergreen-ci/evergreen/db/bsonutil"
	"github.com/evergreen-ci/evergreen/model/task"
	"github.com/evergreen-ci/evergreen/plugin"
	"github.com/evergreen-ci/evergreen/thirdparty"
	"github.com/evergreen-ci/evergreen/util"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func init() {
	plugin.Publish(&IframePlugin{})
}

const (
	PluginName = "iframe"
)

type IframePlugin struct {
	opts *ifpluginOptions
}

func (ifp *IframePlugin) Name() string {
	return PluginName
}

// GetUIHandler adds a path for looking up build failures in JIRA.
func (ifp *IframePlugin) GetUIHandler() http.Handler {
	//if ifp.opts == nil {
	//	panic("build baron plugin missing configuration")
	//}
	r := mux.NewRouter()
	//r.Path("/note/{task_id}").Methods("GET").HandlerFunc(ifp.getNote)
	return r
}

func (ifp *IframePlugin) Configure(conf map[string]interface{}) error {
	//TODO
	return nil
}

func (ifp *IframePlugin) GetPanelConfig() (*plugin.PanelConfig, error) {
	return &plugin.PanelConfig{
		Panels: []plugin.UIPanel{
			{
				Page:      plugin.TaskPage,
				Position:  plugin.PageRight,
				PanelHTML: template.HTML(`<div><h2>HELLOW WORLD</h2></div>`),
				Includes: []template.HTML{
					template.HTML(`<script type="text/javascript" src="/plugin/buildbaron/static/js/task_build_baron.js"></script>`),
				},
				DataFunc: func(context plugin.UIContext) (interface{}, error) {
					return struct {
						Task string `json:"task"`
					}{plugin.UIContext.Task.Id}, nil
				},
			},
		},
	}, nil
}
