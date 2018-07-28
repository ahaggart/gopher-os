package multiprocessing

const (
	// ThreadReady is the status code for a thread that is ready to execute
	ThreadReady = iota

	// ThreadActive is the status code for a thread that is currently executing
	ThreadActive

	// ThreadWaiting is the status code for a thread that is waiting for another
	//  action to complete
	ThreadWaiting
)
