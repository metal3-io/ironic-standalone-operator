package ironic

type Status struct {
	// Object is reconciled and all resources are ready.
	Ready bool
	// Fatal error, further reconciliation is not possible.
	Fatal error
	// Message explaining what is not ready.
	Message string
	// Whether a requeue will be needed.
	requeue bool
}

func (status Status) IsError() bool {
	return status.Fatal != nil
}

func (status Status) IsReady() bool {
	return status.Ready && !status.IsError()
}

func (status Status) NeedsRequeue() bool {
	if status.Fatal != nil {
		return false
	}

	return status.requeue
}

func (status Status) String() string {
	if status.Fatal != nil {
		return status.Fatal.Error()
	}

	if !status.Ready {
		if status.Message != "" {
			return status.Message
		}
		return "resources are not ready yet"
	}

	return "resources are available"
}

// Everything is done, no more reconciliation required.
func ready() (Status, error) {
	return Status{Ready: true}, nil
}

// We have updated dependent resources.
func updated() (Status, error) {
	return Status{Message: "dependent resources are being updated", requeue: true}, nil
}

// We are passively waiting for something external to happen.
func inProgress(message string) (Status, error) {
	return Status{Message: message}, nil
}

// Checking or updating status failed, we hope it's going to resolve itself
// (e.g. a glitch in access to Kube API).
func transientError(err error) (Status, error) {
	return Status{requeue: true}, err
}
