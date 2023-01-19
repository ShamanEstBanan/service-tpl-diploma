package app

//
//type Job interface {
//	Run(ctx context.Context) error
//}
//
//type MyBestSuperJob struct {
//	function func(interface{}) error
//}
//
//func (j *MyBestSuperJob) Run(ctx context.Context) error {
//	err := j.function
//	if err != nil {
//		return err
//	}
//	rand.Seed(time.Now().UnixNano())
//	delay := time.Second * time.Duration(rand.Intn(10))
//	time.Sleep(delay)
//	fmt.Println(delay)
//	if delay == 9*time.Second {
//		return errs.New("test error")
//	}
//	return nil
//}
//
//func RunPool(ctx context.Context, size int, jobs chan Job) error {
//	gr, ctx := errgroup.WithContext(ctx)
//	for i := 0; i < size; i++ {
//		gr.Go(func() error {
//			for {
//				select {
//				case job := <-jobs:
//					err := job.Run(ctx)
//					if err != nil {
//						fmt.Printf("Job error: %s \n", err)
//						return err
//					}
//				case <-ctx.Done():
//					fmt.Println("Context canceled")
//					return ctx.Err()
//				}
//			}
//		})
//	}
//
//	return gr.Wait()
//}
