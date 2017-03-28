package dashboard

import (
	"encoding/json"
	"fmt"
	client "github.com/camunda-ci/camunda-ci-dashboard/http"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	jsonAPI = "/api/json"
)

var (
	// Debug set to true enable request debugging
	Debug = false
)

// Jenkins is high-level API for accessing the underlying Jenkins instance.
type Jenkins interface {
	GetQueue() (*Queue, error)
	GetJobsFromView(viewName string) ([]Job, error)
	GetJobsFromViewWithTree(viewName string, tree string) ([]Job, error)
	GetOverallLoad() (*OverallLoad, error)
	GetExecutors() (*Executors, error)
	GetBusyExecutors() (int, error)
}

// JenkinsClient implements the Jenkins interface and holds the client connected to the underlying Jenkins instance.
type JenkinsClient struct {
	client *client.HTTPClient
}

// Queue represents the Jenkins Build queue.
type Queue struct {
	Items []struct {
		Actions []struct {
			Causes []struct {
				ShortDescription string `json:"shortDescription"`
				UpstreamBuild    int    `json:"upstreamBuild"`
				UpstreamProject  string `json:"upstreamProject"`
				UpstreamURL      string `json:"upstreamUrl"`
			} `json:"causes"`
		} `json:"actions"`
		Blocked      bool   `json:"blocked"`
		Buildable    bool   `json:"buildable"`
		ID           int    `json:"id"`
		InQueueSince int64  `json:"inQueueSince"`
		Params       string `json:"params"`
		Stuck        bool   `json:"stuck"`
		Task         struct {
			Name  string `json:"name"`
			URL   string `json:"url"`
			Color string `json:"color"`
		} `json:"task"`
		URL                        string `json:"url"`
		Why                        string `json:"why"`
		BuildableStartMilliseconds int64  `json:"buildableStartMilliseconds"`
		Pending                    bool   `json:"pending"`
	} `json:"items"`
}

func (q *Queue) String() string {
	return fmt.Sprintf("%#v", q)
}

// Job represents a Jenkins job
type Job struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Color     string `json:"color"`
	LastBuild struct {
		Actions []struct {
			FailCount          int           `json:"failCount,omitempty"`
			SkipCount          int           `json:"skipCount,omitempty"`
			TotalCount         int           `json:"totalCount,omitempty"`
			FoundFailureCauses []interface{} `json:"foundFailureCauses,omitempty"`
		} `json:"actions"`
	} `json:"lastBuild"`
}

func (j *Job) String() string {
	return fmt.Sprintf("%#v", j)
}

// View represents a view inside Jenkins including all jobs on it.
type View struct {
	Jobs []Job `json:"jobs"`
}

func (v *View) String() string {
	return fmt.Sprintf("%#v", v)
}

// Executors represents the configured executors of the underlying Jenkins instance.
type Executors struct {
	BusyExecutors int `json:"busyExecutors"`
	Computer      []struct {
		Actions []struct {
		} `json:"actions"`
		DisplayName string `json:"displayName"`
		Executors   []struct {
		} `json:"executors"`
		Icon            string `json:"icon"`
		IconClassName   string `json:"iconClassName"`
		Idle            bool   `json:"idle"`
		JnlpAgent       bool   `json:"jnlpAgent"`
		LaunchSupported bool   `json:"launchSupported"`
		LoadStatistics  struct {
		} `json:"loadStatistics"`
		ManualLaunchAllowed bool `json:"manualLaunchAllowed"`
		MonitorData         struct {
			HudsonNodeMonitorsSwapSpaceMonitor struct {
				AvailablePhysicalMemory int64 `json:"availablePhysicalMemory"`
				AvailableSwapSpace      int64 `json:"availableSwapSpace"`
				TotalPhysicalMemory     int64 `json:"totalPhysicalMemory"`
				TotalSwapSpace          int64 `json:"totalSwapSpace"`
			} `json:"hudson.node_monitors.SwapSpaceMonitor"`
			HudsonNodeMonitorsArchitectureMonitor string `json:"hudson.node_monitors.ArchitectureMonitor"`
			HudsonNodeMonitorsResponseTimeMonitor struct {
				Average int `json:"average"`
			} `json:"hudson.node_monitors.ResponseTimeMonitor"`
			HudsonNodeMonitorsTemporarySpaceMonitor struct {
				Path string `json:"path"`
				Size int64  `json:"size"`
			} `json:"hudson.node_monitors.TemporarySpaceMonitor"`
			HudsonNodeMonitorsDiskSpaceMonitor struct {
				Path string `json:"path"`
				Size int64  `json:"size"`
			} `json:"hudson.node_monitors.DiskSpaceMonitor"`
			HudsonNodeMonitorsClockMonitor struct {
				Diff int `json:"diff"`
			} `json:"hudson.node_monitors.ClockMonitor"`
		} `json:"monitorData"`
		NumExecutors       int           `json:"numExecutors"`
		Offline            bool          `json:"offline"`
		OfflineCause       interface{}   `json:"offlineCause"`
		OfflineCauseReason string        `json:"offlineCauseReason"`
		OneOffExecutors    []interface{} `json:"oneOffExecutors"`
		TemporarilyOffline bool          `json:"temporarilyOffline"`
	} `json:"computer"`
	DisplayName    string `json:"displayName"`
	TotalExecutors int    `json:"totalExecutors"`
}

func (e *Executors) String() string {
	return fmt.Sprintf("%#v", e)
}

// OverallLoad represents the overall load of the underlying Jenkins instance.
type OverallLoad struct {
	AvailableExecutors struct {
		Hour struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"hour"`
		Min struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"min"`
		Sec10 struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"sec10"`
	} `json:"availableExecutors"`
	BusyExecutors struct {
		Hour struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"hour"`
		Min struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"min"`
		Sec10 struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"sec10"`
	} `json:"busyExecutors"`
	ConnectingExecutors struct {
		Hour struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"hour"`
		Min struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"min"`
		Sec10 struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"sec10"`
	} `json:"connectingExecutors"`
	DefinedExecutors struct {
		Hour struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"hour"`
		Min struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"min"`
		Sec10 struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"sec10"`
	} `json:"definedExecutors"`
	IdleExecutors struct {
		Hour struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"hour"`
		Min struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"min"`
		Sec10 struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"sec10"`
	} `json:"idleExecutors"`
	OnlineExecutors struct {
		Hour struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"hour"`
		Min struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"min"`
		Sec10 struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"sec10"`
	} `json:"onlineExecutors"`
	QueueLength struct {
		Hour struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"hour"`
		Min struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"min"`
		Sec10 struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"sec10"`
	} `json:"queueLength"`
	TotalExecutors struct {
		Hour struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"hour"`
		Min struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"min"`
		Sec10 struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"sec10"`
	} `json:"totalExecutors"`
	TotalQueueLength struct {
		Hour struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"hour"`
		Min struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"min"`
		Sec10 struct {
			History []float64 `json:"history"`
			Latest  float64   `json:"latest"`
		} `json:"sec10"`
	} `json:"totalQueueLength"`
}

func (o *OverallLoad) String() string {
	return fmt.Sprintf("%#v", o)
}

// GetQueue retrieves the Queue of the underlying Jenkins instance.
// It will return an error, if the connection or the JSON un-marshalling breaks.
func (j *JenkinsClient) GetQueue() (*Queue, error) {
	response, _ := j.client.GetFrom("queue" + jsonAPI)

	queue := &Queue{}
	if error := j.processQueueResponse(response, queue); error != nil {
		return nil, error
	}

	return queue, nil
}

func (j *JenkinsClient) processQueueResponse(resp *http.Response, queue *Queue) error {
	return processResponse(resp, queue, "Queue")
}

// GetJobsFromView returns a slice with all Jobs from a given View name from the underlying Jenkins instance.
// It will return an error, if the connection or the JSON un-marshalling breaks.
func (j *JenkinsClient) GetJobsFromView(viewName string) ([]Job, error) {
	response, _ := j.client.GetFrom("view/" + viewName + jsonAPI)

	view := &View{}
	if error := j.processViewResponse(response, view); error != nil {
		return nil, error
	}

	return view.Jobs, nil
}

// GetJobsFromViewWithTree returns a slice with all Jobs from a given View name, restricting the returned attributes by the given tree string.
// It will return an error, if the connection or the JSON un-marshalling breaks.
func (j *JenkinsClient) GetJobsFromViewWithTree(viewName string, tree string) ([]Job, error) {
	response, _ := j.client.GetFrom("view/" + viewName + jsonAPI + "?tree=" + tree)

	view := &View{}
	if error := j.processViewResponse(response, view); error != nil {
		return nil, error
	}

	return view.Jobs, nil
}

func (j *JenkinsClient) processViewResponse(resp *http.Response, view *View) error {
	return processResponse(resp, view, "View")
}

// GetOverallLoad returns the OverallLoad of the underlying Jenkins instance.
// It will return an error, if the connection or the JSON un-marshalling breaks.
func (j *JenkinsClient) GetOverallLoad() (*OverallLoad, error) {
	// TODO: implement
	return &OverallLoad{}, nil
}

// GetExecutors returns the currently configured Executors of the underlying Jenkins instance.
// It will return an error, if the connection or the JSON un-marshalling breaks.
func (j *JenkinsClient) GetExecutors() (*Executors, error) {
	response, _ := j.client.GetFrom("computer" + jsonAPI)

	executors := &Executors{}
	if error := j.processExecutorsResponse(response, executors); error != nil {
		return nil, error
	}

	return executors, nil
}

// GetBusyExecutors returns the number of currently occupied Executors of the underlying Jenkins instance.
// It will return an error, if the connection or the JSON un-marshalling breaks.
func (j *JenkinsClient) GetBusyExecutors() (int, error) {
	response, _ := j.client.GetFrom("computer" + jsonAPI + "?tree=busyExecutors")

	executors := &Executors{}
	if error := j.processExecutorsResponse(response, executors); error != nil {
		return -1, error
	}

	return executors.BusyExecutors, nil
}

func (j *JenkinsClient) processExecutorsResponse(resp *http.Response, executors *Executors) error {
	return processResponse(resp, executors, "Executors")
}

// NewJenkinsClient returns a new Jenkins instance with the given url, username and password.
func NewJenkinsClient(url string, username string, password string) Jenkins {
	config := client.NewHTTPConfig(url, username, password, "application/json")
	client := client.NewHTTPClient(config)

	return &JenkinsClient{client: client}
}

// Process given resp and un-marshall it to the given v.
// Throws either an error if the resp.Body couldn't be read or the un-marshalling failed.
func processResponse(resp *http.Response, v interface{}, component string) error {
	if Debug {
		dumpResponse, error := httputil.DumpResponse(resp, true)
		if error != nil {
			log.Printf("[DEBUG] %s: %s", component, error)
		}
		log.Printf("[DEBUG][REQ]: %s: %s", component, resp.Request.URL)
		log.Printf("[DEBUG][RESP]: %s: %s", component, string(dumpResponse))
	}

	if resp.Body != nil {
		defer resp.Body.Close()

		body, error := ioutil.ReadAll(resp.Body)
		if error != nil {
			return error
		}

		return json.Unmarshal(body, v)
	}

	return nil
}
