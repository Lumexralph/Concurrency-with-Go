package main

func orDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				// if channel has been closed
				if ok == false {
					return
				}
				// continue reading value from the channel
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func tee(done, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)

		for val := range orDone(done, in) {
			// create a local copy of the channels
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				// after writing to the channel,we set it's local
				// or shadowed copy to nil so that further writes
				// will block and the other channel can continue
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}

			}

		}
	}()
	return out1, out2
}
