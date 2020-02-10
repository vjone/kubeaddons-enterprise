package spark_test

import (
	"github.com/mesosphere/kubeaddons-enterprise/test/utils"
	"testing"
)

func TestSparkGroup(t *testing.T) {
	if err := utils.GroupTest(t, "spark", "../../addons"); err != nil {
		t.Fatal(err)
	}
}
