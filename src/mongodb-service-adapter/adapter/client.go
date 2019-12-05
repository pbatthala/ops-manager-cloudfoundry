package adapter

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/opsmanager"
	"github.com/tidwall/gjson"

	"mongodb-service-adapter/digest"
)

type OMClient struct {
	URL      string
	Username string
	APIKey   string
}

type Automation struct {
	MongoDbVersions []MongoDbVersionsType
}

type MongoDbVersionsType struct {
	Name string
}

type GroupCreateRequest struct {
	Name  string   `json:"name"`
	OrgID string   `json:"orgId,omitempty"`
	Tags  []string `json:"tags"`
}

type GroupUpdateRequest struct {
	Tags []string `json:"tags"`
}

type GroupHosts struct {
	TotalCount int `json:"totalCount"`
}

var versionsManifest = []string{"/var/vcap/packages/versions/versions.json", "../../mongodb_versions/versions.json"}

func (oc *OMClient) Client() opsmanager.Client {
	r := httpclient.NewURLResolverWithPrefix(oc.URL, "api/public/v1.0")
	return opsmanager.NewClientWithDigestAuth(r, oc.Username, oc.APIKey)
}

func (oc *OMClient) CreateGroup(id string, request GroupCreateRequest) (opsmanager.ProjectResponse, error) {
	log.Println(fmt.Sprintf("CreateGroup in id : %s ,request : %+v", id, request))

	if request.Name == "" {
		request.Name = fmt.Sprintf("PCF_%s", strings.ToUpper(id))
	}

	safeclient := opsmanager.NewClient(
		opsmanager.WithResolver(httpclient.NewURLResolverWithPrefix(oc.URL, "api/public/v1.0")),
		opsmanager.WithHTTPClient(httpclient.NewClient(
			httpclient.WithDigestAuthentication(oc.Username, oc.APIKey),
			httpclient.WithAcceptedStatusCodes([]int{http.StatusOK, http.StatusCreated, http.StatusNotFound}),
		)),
	)

	group, err := safeclient.GetProjectByName(request.Name)
	if err != nil {
		log.Printf("CreateGroup GetGroupByName with request.Name: %s, error: %v", request.Name, err)
		return group, err
	}

	if group.Name == request.Name {
		log.Printf("Continue with existing group %q", group.ID)
		apiKey, err := safeclient.CreateAgentAPIKEY(group.ID, "MongoDB On-Demand broker generated Agent API Key")
		if err != nil {
			log.Printf("CreateGroup CreateGroupAPIKey group.ID: %s, error: %v", group.ID, err)
			return group, err
		}
		group.AgentAPIKey = apiKey.Key
		return group, nil
	}

	resp, err := safeclient.CreateOneProject(request.Name, request.OrgID)
	if err != nil {
		log.Printf("CreateGroup CreateOneProject, request: %+v, error: %v", request, err)
		return group, err
	}

	return resp, nil
}

func (oc *OMClient) GetLatestVersion(groupID string) (string, error) {
	cfg, err := oc.Client().GetAutomationConfig(groupID)
	if err != nil || len(cfg.MongoDBVersions) == 0 {
		return "", fmt.Errorf("unable to find the latest MongoDB version from the MongoDB Ops Manager API. Please contact your system administrator to ensure versions are available in the Version Manager for group '%q' in MongoDB Ops Manager. If your MongoDB Ops Manager is running in Local Mode, then after validating versions are available, please indicate a specific MongoDB version using 'version’ paramater when calling 'create-service'", groupID)
	}

	versions := make([]string, len(cfg.MongoDBVersions))
	n := 0
	for _, i := range cfg.MongoDBVersions {
		if !strings.HasSuffix(i.Name, "ent") {
			versions[n] = i.Name
			n++
		}
	}
	versions = versions[:n]
	latestVersion := versions[len(versions)-1]

	return latestVersion, nil
}

func (oc *OMClient) ValidateVersion(groupID string, version string) (string, error) {
	cfg, err := oc.Client().GetAutomationConfig(groupID)
	if err != nil {
		return "", err
	}

	for _, v := range cfg.MongoDBVersions {
		if v.Name == version {
			return version, nil
		}
	}

	return "", errors.New("failed to find expected version, got " + version)
}

func (oc *OMClient) ValidateVersionManifest(version string) (string, error) {
	b, err := ioutil.ReadFile(versionsManifest[0])
	if err != nil {
		b, err = ioutil.ReadFile(versionsManifest[1])
		if err != nil {
			return "", err
		}
	}

	v := gjson.GetBytes(b, fmt.Sprintf(`versions.#[name="%s"].name`, version))
	log.Printf("Using %q version of MongoDB", v.String())
	if v.String() == "" {
		return "", errors.New("failed to find expected version, continue with provided versions ")
	}

	return version, nil
}

func (oc *OMClient) HasBackupAgent(groupID string) (bool, error) {
	u := fmt.Sprintf("/api/public/v1.0/groups/%s/agents/BACKUP", groupID)
	b, err := oc.doRequest("GET", u, nil)
	if err != nil {
		return false, err
	}
	state := gjson.GetBytes(b, "results.#.stateName").String()
	if strings.Contains(state, "ACTIVE") {
		return true, err
	}
	return false, err
}

func (oc *OMClient) doRequest(method string, path string, body io.Reader) ([]byte, error) {
	uri := fmt.Sprintf("%s%s", oc.URL, path)
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if err = digest.ApplyDigestAuth(oc.Username, oc.APIKey, uri, req); err != nil {
		return nil, err
	}
	log.Printf("API Call: %s%s", oc.URL, path)

	// dump, err := httputil.DumpRequestOut(req, true)
	// if err != nil {
	// 	return nil, err
	// }

	// // log.Printf("API Request: %q", dump)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// dump, err = httputil.DumpResponse(res, true)
	// if err != nil {
	// 	return nil, err
	// }
	// // log.Printf("API Response: %q", dump)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// oc.GetGroupByName return 404 if group not found
	if res.StatusCode == 404 {
		log.Printf("Received %d status code for %s path", res.StatusCode, path)
		return b, nil
	} else if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("%s %s request error: code=%d body=%q", method, path, res.StatusCode, b)
	}
	return b, nil
}