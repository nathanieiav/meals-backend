package cronservice

import (
	"fmt"
	"project-skbackend/configs"
	"project-skbackend/internal/repositories/orderrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type (
	CronService struct {
		cfg  *configs.Config
		rodr orderrepo.IOrderRepository
	}

	ICronService interface {
		Init() (gocron.Scheduler, error)
	}
)

func NewCronService(
	cfg *configs.Config,
	rodr orderrepo.IOrderRepository,
) *CronService {
	return &CronService{
		cfg:  cfg,
		rodr: rodr,
	}
}

func (s *CronService) Init() (gocron.Scheduler, error) {
	// * get time location from env
	tz, err := time.LoadLocation(s.cfg.API.Timezone)
	if err != nil {
		return nil, err
	}

	// * define a new scheduler instance
	gsch, err := gocron.NewScheduler(
		gocron.WithLocation(tz),
	)

	if err != nil {
		return nil, err
	}

	// * add a order job
	s.orderSchedule(gsch)

	// * start the scheduler
	gsch.Start()

	return gsch, nil
}

func (s *CronService) orderSchedule(gsch gocron.Scheduler) {
	var (
		errs []error
	)

	err := s.scheduleOrderCancelled(gsch)
	if err != nil {
		errs = append(errs, err)
	}

	err = s.scheduleOrderPickedUp(gsch)
	if err != nil {
		errs = append(errs, err)
	}

	err = s.scheduleOrderOutForDelivery(gsch)
	if err != nil {
		errs = append(errs, err)
	}

	err = s.scheduleOrderDelivered(gsch)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		utlogger.Error(err)
	}
}

func (s *CronService) scheduleOrderCancelled(gsch gocron.Scheduler) error {
	_, err := gsch.NewJob(
		gocron.DurationJob(
			// TODO: use the env variable
			time.Duration(1)*time.Minute,
		),
		gocron.NewTask(
			func() error {
				err := s.rodr.UpdateAutomaticallyStatus(consttypes.OS_CANCELLED, s.cfg.OrderBuffer.AutomaticallyCancelled, []consttypes.OrderStatus{consttypes.OS_PLACED})
				if err != nil {
					utlogger.Error(err)
					return err
				}

				return nil
			},
		),
	)

	if err != nil {
		utlogger.Error(err)
		return err
	}

	utlogger.Info(fmt.Sprintf("Service for Cron %s Running!", "Update Order Expired"))

	return nil
}

func (s *CronService) scheduleOrderPickedUp(gsch gocron.Scheduler) error {
	_, err := gsch.NewJob(
		gocron.DurationJob(
			// TODO: use the env variable
			time.Duration(1)*time.Minute,
		),
		gocron.NewTask(
			func() error {
				err := s.rodr.UpdateAutomaticallyStatus(consttypes.OS_PICKED_UP, s.cfg.OrderBuffer.AutomaticallyBeingPickedUp, []consttypes.OrderStatus{consttypes.OS_PREPARED})
				if err != nil {
					utlogger.Error(err)
					return err
				}

				return nil
			},
		),
	)

	if err != nil {
		utlogger.Error(err)
		return err
	}

	utlogger.Info(fmt.Sprintf("Service for Cron %s Running!", "Update Order Picked Up"))

	return nil
}

func (s *CronService) scheduleOrderOutForDelivery(gsch gocron.Scheduler) error {
	_, err := gsch.NewJob(
		gocron.DurationJob(
			// TODO: use the env variable
			time.Duration(1)*time.Minute,
		),
		gocron.NewTask(
			func() error {
				err := s.rodr.UpdateAutomaticallyStatus(consttypes.OS_OUT_FOR_DELIVERY, s.cfg.OrderBuffer.AutomaticallyOutForDelivery, []consttypes.OrderStatus{consttypes.OS_PICKED_UP})
				if err != nil {
					utlogger.Error(err)
					return err
				}

				return nil
			},
		),
	)

	if err != nil {
		utlogger.Error(err)
		return err
	}

	utlogger.Info(fmt.Sprintf("Service for Cron %s Running!", "Update Order Out for Delivery"))

	return nil
}

func (s *CronService) scheduleOrderDelivered(gsch gocron.Scheduler) error {
	_, err := gsch.NewJob(
		gocron.DurationJob(
			// TODO: use the env variable
			time.Duration(1)*time.Minute,
		),
		gocron.NewTask(
			func() error {
				err := s.rodr.UpdateAutomaticallyStatus(consttypes.OS_DELIVERED, s.cfg.OrderBuffer.AutomaticallyDelivered, []consttypes.OrderStatus{consttypes.OS_OUT_FOR_DELIVERY})
				if err != nil {
					utlogger.Error(err)
					return err
				}

				return nil
			},
		),
	)

	if err != nil {
		utlogger.Error(err)
		return err
	}

	utlogger.Info(fmt.Sprintf("Service for Cron %s Running!", "Update Order Delivered"))

	return nil
}
