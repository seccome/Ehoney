package harbor

import (
    "decept-defense/models"
    "decept-defense/pkg/configs"
    "decept-defense/pkg/util"
    "encoding/json"
    "fmt"
    "go.uber.org/zap"
    "net/url"
    "strings"
    "time"
)

type Repositories struct {
    ArtifactCount       int           `json:"artifact_count"`
    CreationTime        string        `json:"creation_time"`
    ID                  int           `json:"id"`
    Name                string        `json:"name"`
    ProjectID           int           `json:"project_id"`
    UpdateTime          string        `json:"update_time"`
}

type Artifacts struct {
    Digest              string        `json:"digest"`
    Icon                string        `json:"icon"`
    Id                  int           `json:"id"`
    Labels              string        `json:"labels"`
    ManifestMediaType   string        `json:"manifest_media_type"`
    MediaType           string        `json:"media_type"`
    ProjectID           int           `json:"project_id"`
    PullTime            string        `json:"pull_time"`
    PushTime            string        `json:"push_time"`
    References          string        `json:"references"`
    RepositoryID        int           `json:"repository_id"`
    Size                int           `json:"size"`
    Tags                []Tags        `json:"tags"`
    Type                string        `json:"type"`
}

type Tags struct {
    ArtifactID    int    `json:"artifact_id"`
    ID            int    `json:"id"`
    Immutable     bool   `json:"immutable"`
    Name          string `json:"name"`
    PullTime      string `json:"pull_time"`
    PushTime      string `json:"push_time"`
    RepositoryID  int    `json:"repository_id"`
    Signed        bool   `json:"signed"`
}


func Setup(){
    ticker := time.NewTicker(time.Hour * 3)
    done := make(chan bool)
    RefreshImages()
    go func() {
        for {
            select {
            case <-done:
                return
            case t := <-ticker.C:
                fmt.Println("Tick at", t)
                RefreshImages()
            }
        }
    }()
}

func RefreshImages() error{
    var repositories []Repositories
    var images []models.Images
    parser, err  := url.Parse(configs.GetSetting().Harbor.HarborURL)
    if err != nil{
        return err
    }
    harborHost := parser.Host
    header := map[string]string{
        "authorization" : "Basic " + configs.GetSetting().Harbor.User + ":" + configs.GetSetting().Harbor.Password,
    }
    requestURI := strings.Join([]string{configs.GetSetting().Harbor.HarborURL, "api", "v" + configs.GetSetting().Harbor.APIVersion, "projects", configs.GetSetting().Harbor.HarborProject, "repositories"}, "/")
    rsp, err := util.SendGETRequest(header, requestURI)

    if err != nil{
        zap.L().Info("SendGETRequest error " + err.Error())
        return err
    }
    err = json.NewDecoder(strings.NewReader(string(rsp))).Decode(&repositories)
    if err != nil{
        zap.L().Info("SendGETRequest error " + err.Error())
        return err
    }
    for _, r := range repositories{
        var artifacts []Artifacts
        requestURI = strings.Join([]string{configs.GetSetting().Harbor.HarborURL, "api", "v" + configs.GetSetting().Harbor.APIVersion, "projects", configs.GetSetting().Harbor.HarborProject, "repositories", strings.Split(r.Name, "/")[1], "artifacts"}, "/")
        rsp, err:= util.SendGETRequest(header, requestURI)
        if err != nil{
            zap.L().Info("SendGETRequest error " + err.Error())
            return err
        }
        err = json.NewDecoder(strings.NewReader(string(rsp))).Decode(&artifacts)
        if err != nil{
            zap.L().Info("SendGETRequest error " + err.Error())
            return err
        }
        for _, s := range artifacts{
            for _, tag := range s.Tags {
                var image models.Images
                image.ImageName = r.Name
                image.ImageAddress = harborHost + "/" + r.Name + ":" + tag.Name
                images = append(images, image)
            }
        }
    }
    for _, i := range images{
        var p models.Images
        p.ImageAddress = i.ImageAddress
        p.ImageName = i.ImageName
        p.ImageType = i.ImageType
        p.ImagePort = i.ImagePort
        err =  p.CreateImage()
        if err != nil{
            return err
        }
    }
    return nil
}
