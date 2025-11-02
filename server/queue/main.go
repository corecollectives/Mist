package queue

func InitQueue() *Queue {
	q := NewQueue(5)
	return q
}

// this function is just for testing purposes
// func TestJobs(q *Queue) {
// 	for i := 1; i <= 5; i++ {
// 		err := q.AddJob(Job{
// 			ID:        i,
// 			CreatedAt: time.Now(),
// 		})
// 		if err != nil {
// 			fmt.Println("error adding job:", err)
// 		}
// 	}

// 	time.Sleep(6 * time.Second)

// 	q.Close()

// }
