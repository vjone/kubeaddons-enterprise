package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/mesosphere/kubeaddons-enterprise/test/utils"
	"github.com/mesosphere/kubeaddons/pkg/api/v1beta1"
	"github.com/mesosphere/kubeaddons/pkg/repositories/local"
)

var addonsDir = "../addons"

func init() {
	b, err := ioutil.ReadFile("groups.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(b, utils.AddonTestingGroups); err != nil {
		panic(err)
	}
}

func TestValidateUnhandledAddons(t *testing.T) {
	unhandled, err := findUnhandled()
	if err != nil {
		t.Fatal(err)
	}

	if len(unhandled) != 0 {
		names := make([]string, len(unhandled))
		for _, addon := range unhandled {
			names = append(names, addon.GetName())
		}
		t.Fatal(fmt.Errorf("the following addons are not handled as part of a testing group: %+v", names))
	}
}

func TestGeneralGroup(t *testing.T) {
	if err := utils.GroupTest(t, "general", addonsDir); err != nil {
		t.Fatal(err)
	}
}

func TestKafkaGroup(t *testing.T) {
	if err := utils.GroupTest(t, "kafka", addonsDir); err != nil {
		t.Fatal(err)
	}
}

func TestCassandraGroup(t *testing.T) {
	if err := utils.GroupTest(t, "cassandra", addonsDir); err != nil {
		t.Fatal(err)
	}
}

func findUnhandled() ([]v1beta1.AddonInterface, error) {
	var unhandled []v1beta1.AddonInterface
	repo, err := local.NewRepository("base", "../addons")
	if err != nil {
		return unhandled, err
	}
	addons, err := repo.ListAddons()
	if err != nil {
		return unhandled, err
	}

	for _, revisions := range addons {
		addon := revisions[0]
		found := false
		for _, v := range utils.AddonTestingGroups {
			for _, name := range v {
				if name == addon.GetName() {
					found = true
				}
			}
		}
		if !found {
			unhandled = append(unhandled, addon)
		}
	}

	return unhandled, nil
}

func kubectl(args ...string) error {
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
