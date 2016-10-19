package iframe

import (
	"github.com/evergreen-ci/evergreen/plugin"
	"github.com/gorilla/mux"
	//"github.com/mitchellh/mapstructure"
	"html/template"
	"net/http"
)

func init() {
	plugin.Publish(&IframePlugin{})
}

const (
	PluginName = "iframe"
)

type IframePlugin struct {
	//TODO
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
				Page:     plugin.TaskPage,
				Position: plugin.PageLeft,
				PanelHTML: template.HTML(`
				 <div ng-controller="IFrameCtrl" class="mci-pod ifplug"><div class="row"><div class="col-lg-12">
				 <iframe ng-src="[[exampleIFrameUrl]]" onload='javascript:(function(o){o.style.height=o.contentWindow.document.body.scrollHeight+"px";}(this));' style="height:400px;width:100%;border:none;overflow:hidden;"></iframe>
			     </div></div></div>`),
				Includes: []template.HTML{
					template.HTML(`<script>
mciModule.controller('IFrameCtrl',
function ($http, $scope, $window, $sce) {  
	$scope.exampleIFrameUrl = $sce.trustAsResourceUrl("http://example.com?" + $window.plugins.iframe.task);
});
					</script>`),
				},
				DataFunc: func(context plugin.UIContext) (interface{}, error) {
					return struct {
						Task string `json:"task"`
					}{context.Task.Id}, nil
				},
			},
		},
	}, nil
}
