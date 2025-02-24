package workerpool

func (w *WorkerPool) Push(action Action) error {
	if w.closed.Load() {
		return errPoolClosed
	}

	w.input <- action
	return nil
}
