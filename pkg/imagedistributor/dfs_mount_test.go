package imagedistributor

import (
	"github.com/sealerio/sealer/pkg/define/options"
	"github.com/sealerio/sealer/pkg/imageengine"
	"testing"
)

func TestPullImage(t *testing.T) {
	imageEngine, err := imageengine.NewImageEngine(options.EngineGlobalConfigurations{})
	if err != nil {
		t.Log(err)
		return
	}
	imageSpec, err := imageEngine.Inspect(&options.InspectOptions{
		ImageNameOrID: "mysql",
	})
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(imageSpec)
}
