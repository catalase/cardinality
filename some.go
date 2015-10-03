package main

type Some []Code

func (some Some) Meet(n int) error {
	errc, returnc := make(chan error), make(chan bool)
	q := make(chan int)

	for i := 0; i < n; i++ {
		go func() {
			var err error
			for {
				select {
				case i := <-q:
					some[i], err = One()
					errc <- err
				case <-returnc:
					return
				}
			}
		}()
	}

	go func() {
		for i := 0; i < len(some); i++ {
			select {
			case q <- i:
			case <-returnc:
				return
			}
		}
	}()

	defer close(returnc)

	for i := 0; i < len(some); i++ {
		if err := <-errc; err != nil {
			return err
		}
	}

	return nil
}
