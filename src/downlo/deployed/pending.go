package deployed

import (
  "log"
  "path/filepath"
  "io/ioutil"
  "encoding/json"
  "downlo"
)

func PendingUpdates(workspace string) (candidates []downlo.Release) {
  dirs := scanDirForManifest(workspace)

  candidates = loadManifests(dirs)

  return
}

func scanDirForManifest(path string) (manifests []string) {
  dirpath, _ := filepath.Abs(path)

  dirpath = filepath.Join(dirpath, "manifests")

  matches, err := ioutil.ReadDir(dirpath)
  if err != nil {
    log.Printf("Error reading path")
  }

  for _, f := range matches {
    if f.IsDir() == false {
      fullFilePath := filepath.Join(dirpath, f.Name())
      manifests = append(manifests, fullFilePath)
    }
  }

  return
}

func loadManifests(paths []string) (candidates []downlo.Release) {
  for _, p := range paths {
    candidates = append(candidates, loadManifest(p))
  }

  return
}

func loadManifest(path string) (candidate downlo.Release) {
  contents, err := ioutil.ReadFile(path)

  if err != nil {
    log.Printf("Error reading file: %s \n", path)
  }

  parseError := json.Unmarshal(contents, &candidate)

  if parseError != nil {
    log.Printf("Error reading deploy: %s \n", path)
    log.Println(parseError)
  }
  return
}