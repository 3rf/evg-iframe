package iframe

import (
	"fmt"
	"github.com/evergreen-ci/evergreen/plugin"
	"github.com/mitchellh/mapstructure"
	"html/template"
	"net/http"
	"strings"
)

/* EXAMPLE CONFIG
plugins:
    iframe:
        kitty:
            title: "Skunkworks Plugin"
            height: 150
            route: "http://localhost:7777/diff?task=${task}"
        json:
            title: "JSON View"
            height: 600
            route: "http://localhost:7777/show?task=${task}"
*/

func init() {
	plugin.Publish(&IframePlugin{})
}

const (
	PluginName = "iframe"
)

type IframePlugin struct {
	PanelCache []plugin.UIPanel
}

type ifConf struct {
	Title  string
	Route  string
	Height int
}

func (ifp *IframePlugin) Name() string {
	return PluginName
}

func (ifp *IframePlugin) GetUIHandler() http.Handler {
	return nil
}

func (ifp *IframePlugin) Configure(conf map[string]interface{}) error {
	ifs := map[string]ifConf{}
	err := mapstructure.Decode(conf, &ifs)
	if err != nil {
		return err
	}
	count := 0
	for id, vals := range ifs {
		panel := plugin.UIPanel{
			Page:      plugin.TaskPage,
			Position:  plugin.PageLeft,
			PanelHTML: template.HTML(buildTemplate(id, vals.Title, vals.Height)),
			Includes: []template.HTML{
				template.HTML(buildController(id, vals.Route)),
			},
		}
		if count == 0 {
			panel.DataFunc = func(context plugin.UIContext) (interface{}, error) {
				return struct {
					Task string `json:"task"`
				}{context.Task.Id}, nil
			}
		}
		count++
		ifp.PanelCache = append(ifp.PanelCache, panel)
	}
	return nil
}

func (ifp *IframePlugin) GetPanelConfig() (*plugin.PanelConfig, error) {
	return &plugin.PanelConfig{
		Panels: ifp.PanelCache,
	}, nil
}

func buildTemplate(name, title string, height int) string {
	return fmt.Sprintf(`<h3 class="section-heading"> <i class="fa fa-external-link-square"></i> %s</h3> 
<div ng-controller="%sIFrameCtrl" class="mci-pod ifplug"><div class="row"><div class="col-lg-12">
<iframe ng-src="[[iframeURL]]" onload='javascript:(function(o){o.style.height=o.contentWindow.document.body.scrollHeight+"px";}(this));' style="height: %vpx; width:100%%; border:none; overflow:hidden;"></iframe>
</div></div></div>`, title, name, height)
}

func buildController(name, route string) string {
	routePieces := strings.Split(route, "${task}")
	for i := range routePieces {
		routePieces[i] = fmt.Sprintf(`"%s"`, routePieces[i])
	}
	routeFormula := routePieces[0] + " + $window.plugins.iframe.task + " + routePieces[1]
	return fmt.Sprintf(`<script>
mciModule.controller('%sIFrameCtrl',
function ($http, $scope, $window, $sce) {  
	$scope.iframeURL = $sce.trustAsResourceUrl(%s);
});
					</script>`, name, routeFormula)

}
