package task

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	js "github.com/bitly/go-simplejson"
	"github.com/juju/errors"
)

type Job struct {
	Id            int64  `db:"id" json:"id"`
	Name          string `db:"name" json:"name"`                     // 512, unique
	Executor      string `db:"executor" json:"executor"`             // 4096
	ExecutorFlags string `db:"executor_flags" json:"executor_flags"` // 4096
	Retries       int    `db:"retries" json:"retries"`
	Owner         string `db:"owner" json:"owner"`
	SuccessCnt    int    `db:"success_cnt" json:"success_cnt"`
	ErrCnt        int    `db:"error_cnt" json:"error_cnt"`
	CreateTs      int64  `db:"create_ts" json:"create_ts"`
	LastTaskId    string `db:"last_task_id" json:"last_task_id"`
	LastSuccessTs int64  `db:"last_success_ts" json:"last_success_ts"`
	LastErrTs     int64  `db:"last_error_ts" json:"last_error_ts"`
	LastStatus    string `db:"last_status" json:"last_status"`
	Cpus          int    `db:"cpus" json:"cpus"`
	Mem           int    `db:"mem" json:"mem"`
	Disk          int64  `db:"disk" json:"disk"`
	Disabled      bool   `db:"disabled" json:"disabled"`
	Uris          string `db:"uris" json:"uris"` // 2048, using comma to split
	Schedule      string `db:"schedule" json:"schedule"`
	WebHookUrl    string `db:"hook" json:"hook"`
}

type JobHelper struct {
	Server       string
	t            http.Transport //reuse connection
	ExecutorUrls string
}

func (jh *JobHelper) getCreateJobUrl() string {
	return jh.Server + "/job"
}

func (jh *JobHelper) getRunJobUrl(j *Job) string {
	return jh.Server + "/job/run/" + strconv.Itoa(int(j.Id))
}

func (jh *JobHelper) BuildRepoJob(repo string) *Job {
	return &Job{
		Executor:      "./example_executor",
		ExecutorFlags: "./startdrone.sh " + repo,
		Owner:         "CI ROBOT",
		Name:          repo,
		Uris:          jh.ExecutorUrls,
	}
}

func (jh *JobHelper) CreateJob(j *Job) error {
	//todo: post to create job
	c := http.Client{Transport: &jh.t}
	buf, err := json.Marshal(j)
	if err != nil {
		return errors.Trace(err)
	}

	resp, err := c.Post(jh.getCreateJobUrl(), "text/json", bytes.NewReader(buf))
	if err != nil {
		return errors.Trace(err)
	}

	defer resp.Body.Close()

	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}

	obj, err := js.NewJson(buf)
	if err != nil {
		return errors.Trace(err)
	}

	id, ok := obj.Get("data").CheckGet("id")
	if ok {
		println(id)
	}

	j.Id, err = id.Int64()
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (jh *JobHelper) RunJob(j *Job) error {
	c := http.Client{Transport: &jh.t}
	buf := []byte("{}") //empty json for now
	resp, err := c.Post(jh.getRunJobUrl(j), "text/json", bytes.NewReader(buf))
	if err != nil {
		return errors.Trace(err)
	}

	defer resp.Body.Close()

	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}
