package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"service-tpl-diploma/internal/domain"
	"time"
)

type storage interface {
	GetOrdersForProcessing(ctx context.Context) ([]string, error)
	UpdateOrder(ctx context.Context, orderID string, status string, accrual float32) error
	UpdateAccountBalance(ctx context.Context, orderID string, accrual float32) error
}

type updateOrderStatusJob struct {
	orderID              string
	accrualSystemAddress string
	st                   storage
	lg                   *zap.Logger
}

func (j *updateOrderStatusJob) Run(ctx context.Context) error {
	j.lg.Info("OrderId in job:", zap.String("orderId:", j.orderID))
	//тут юзкейс похода в сервис чужой
	//получаем данные по orderID
	url := fmt.Sprintf("%s/api/orders/%s", j.accrualSystemAddress, j.orderID)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if response.StatusCode != http.StatusOK {
		respMessage := fmt.Sprintf("Responce:\nCode %v\n Message:%v", response.StatusCode, response.Body)
		log.Println(respMessage)
		return nil
	}
	orderInfo := domain.AccrualServiceResponse{}
	err = json.NewDecoder(response.Body).Decode(&orderInfo)
	if err != nil {
		log.Println(err)
		return nil
	}

	// обновляем значение в БД если статус INVALID или PROCESSED
	if orderInfo.Status == domain.OrderAccrualStatusINVALID || orderInfo.Status == domain.OrderAccrualStatusPROCESSED {
		err = j.st.UpdateOrder(ctx, orderInfo.OrderId, orderInfo.Status, orderInfo.Accrual)
		if err != nil {
			j.lg.Error("err while update status", zap.Error(err))
		}
		err = j.st.UpdateAccountBalance(ctx, orderInfo.OrderId, orderInfo.Accrual)
	}
	return nil
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
	go func() {
		for {
			orders, err := s.storage.GetOrdersForProcessing(context.Background())
			if err != nil {
				s.lg.Error("Err while take all order to processing: ", zap.Error(err))
			}
			for _, order := range orders {
				s.jobsCh <- &updateOrderStatusJob{
					orderID:              order,
					accrualSystemAddress: s.accrualSystemAddress,
					st:                   s.storage,
					lg:                   s.lg,
				}
			}
			log.Println("Array of orders for processing: ", orders)

			time.Sleep(20 * time.Second)
		}
	}()
}
