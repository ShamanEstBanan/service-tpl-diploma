package workers

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"log"
	"net/http"
	"service-tpl-diploma/internal/domain"
	"time"
)

type updateOrderStatusJob struct {
	orderID              string
	client               *http.Client
	accrualSystemAddress string
	st                   storage
}

func (j *updateOrderStatusJob) Run(ctx context.Context) error {
	fmt.Println("OrderId in job:", j.orderID)

	//тут юзкейс похода в сервис чужой
	//получаем данные по orderID
	request, err := http.NewRequest(http.MethodGet, j.accrualSystemAddress, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	response, err := j.client.Do(request)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if response.StatusCode != http.StatusOK {
		respMessage := fmt.Sprintf("Responce:\nCode %v\n Message:%v", response.StatusCode, response.Body)
		return errors.New(respMessage)
	}
	// обновляем значение в БД если статус INVALID или PROCESSED
	err = j.st.UpdateOrder(ctx, 0, "", "", decimal.New(0, 0))
	return nil
}

type storage interface {
	GetOrdersForProcessing(ctx context.Context) ([]string, error)
	UpdateOrder(ctx context.Context, orderID int, userID string, status string, accural decimal.Decimal) error
}

type StatusUpdater struct {
	storage              storage
	jobsCh               chan domain.Job
	lg                   *zap.Logger
	accrualSystemAddress string
}

func NewStatusUpdater(st storage, jobCh chan domain.Job, lg *zap.Logger, accrualSystemAddress string) *StatusUpdater {
	return &StatusUpdater{
		storage:              st,
		jobsCh:               jobCh,
		lg:                   lg,
		accrualSystemAddress: accrualSystemAddress,
	}
}

func (s *StatusUpdater) Start() {
	client := NewClient()
	go func() {
		for {
			orders, err := s.storage.GetOrdersForProcessing(context.Background())
			if err != nil {
				s.lg.Error("Err while take all order to processing: ", zap.Error(err))
			}
			for _, order := range orders {
				s.jobsCh <- &updateOrderStatusJob{
					orderID:              order,
					client:               client,
					accrualSystemAddress: s.accrualSystemAddress,
				}
			}
			log.Println("Array of orders for processing: ", orders)

			time.Sleep(2 * time.Second)
		}
	}()
}

func NewClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}
