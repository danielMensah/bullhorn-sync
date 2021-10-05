package worker

//func TestNew(t *testing.T) {
//	tests := []struct {
//		name     string
//		queue    *queue.Mock
//		expected *Worker
//	}{
//		{
//			name:  "creates a new worker",
//			queue: &queue.Mock{},
//			expected: &Worker{
//				queue: &queue.Mock{},
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			worker := New(tt.queue)
//			assert.Equal(t, tt.expected, worker)
//		})
//	}
//}
//
//func TestWorker_Run(t *testing.T) {
//	tests := []struct {
//		name        string
//		queue       *queue.Mock
//		timeout     time.Duration
//		records     []model.Record
//		expectMocks func(t *testing.T, q *queue.Mock)
//	}{
//		{
//			name:    "ok",
//			queue:   &queue.Mock{},
//			timeout: time.Second * 20,
//			records: []model.Record{
//				{
//					Id:   uuid.New(),
//					Name: "Company",
//					Data: `{"age":24}`,
//				},
//			},
//			expectMocks: func(t *testing.T, q *queue.Mock) {
//				q.On("Send", `{"age":24}`).Return(uuid.New().String(), nil)
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
//			defer cancel()
//
//			wkr := New(tt.queue)
//
//			records := make(chan []model.Record)
//
//			wg := &sync.WaitGroup{}
//			wg.Add(1)
//			go wkr.Run(ctx, records, wg)
//
//			records <- tt.records
//
//			wg.Wait()
//			close(records)
//
//			if tt.expectMocks != nil {
//				tt.queue.AssertExpectations(t)
//			}
//		})
//	}
//}
