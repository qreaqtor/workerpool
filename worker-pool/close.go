package workerpool

func (w *WorkerPool) Close() error {
	if w.closed.Swap(true) {
		return errPoolClosed
	}

	w.cancel()
	w.wg.Wait()
	close(w.input)

	return nil
}
