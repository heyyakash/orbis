package modals

type JobRequest struct {
	Command  string `json:"command"`
	Schedule string `json:"schedule"`
}
