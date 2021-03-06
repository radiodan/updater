package updates

import (
	"encoding/json"
	"github.com/radiodan/updater/model"
	"github.com/radiodan/updater/utils"
	"io/ioutil"
	"log"
	"net/http"
)

func LatestReleasesByProject() (projects []Project) {
	deployUrl := "http://deploy.radiodan.net"

	body, err := fetch(deployUrl, utils.MachineHardware())

	if err != nil {
		log.Println("[!] Cannot connect to", deployUrl)
		return
	}

	var data interface{}
	parseError := json.Unmarshal(body, &data)
	if parseError != nil {
		log.Printf("JSON parse error")
	}

	projects = parseJsonToProjects(data)

	return
}

// Fetch body from a URL
func fetch(url string, machineHardware string) (body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Machine-Hardware-Name", machineHardware)

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("HTTP Request Error")
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	return
}

func parseJsonToProjects(data interface{}) (projects []Project) {
	parsed := data.(map[string]interface{})

	for projectName, contents := range parsed {

		refs := contents.(map[string]interface{})

		current := Project{}
		current.Name = projectName

		current.Refs = map[string]model.Release{}

		for refName, ref := range refs {
			r := ref.(map[string]interface{})

			release := model.Release{}
			release.Project = projectName
			release.Ref = refName
			release.Source = r["file"].(string)
			release.Hash = r["sha1"].(string)
			release.Commit = r["commit"].(string)

			current.Refs[refName] = release
		}

		projects = append(projects, current)
	}

	return
}
