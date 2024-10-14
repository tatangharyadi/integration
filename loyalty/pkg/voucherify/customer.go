package voucherify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/tatangharyadi/integration/loyalty/models"
	voucherify "github.com/tatangharyadi/integration/loyalty/pkg/voucherify/models"
)

type Task struct {
	handler  Handler
	customer models.Customer
}

func (t *Task) Process() {
	customer, err := importCustomer(t.handler, t.customer)
	if err != nil {
		t.handler.Logger.Error().Err(err).Msg("failed to import customer")
		return
	}
	t.handler.Logger.Info().Str("customer_id", customer.Id).Msg("customer imported")
	voucher, err := createVoucher(t.handler, customer)
	if err != nil {
		t.handler.Logger.Error().Err(err).Msg("failed to create voucher")
		return
	}
	t.handler.Logger.Info().Str("voucher_id", voucher.Id).Msg("voucher created")
	err = publishVoucher(t.handler, voucher.Code, customer.Id)
	if err != nil {
		t.handler.Logger.Error().Err(err).Msg("failed to publish voucher")
		return
	}
}

type WorkerPool struct {
	Tasks       []Task
	concurrency int
	tasksChan   chan Task
	wg          sync.WaitGroup
}

func (wp *WorkerPool) worker() {
	for task := range wp.tasksChan {
		task.Process()
		wp.wg.Done()
	}
}

func (wp *WorkerPool) Run() {
	wp.tasksChan = make(chan Task, len(wp.Tasks))

	for i := 0; i < wp.concurrency; i++ {
		go wp.worker()
	}

	wp.wg.Add(len(wp.Tasks))
	for _, task := range wp.Tasks {
		wp.tasksChan <- task
	}
	close(wp.tasksChan)

	wp.wg.Wait()
}

func importCustomer(h Handler, customer models.Customer) (models.Customer, error) {
	url := fmt.Sprintf("%s/customers", h.Env.LoyaltyUrl)
	client := &http.Client{}
	customer.MealBenefit.TransactionTimestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	customer.CreditBenefit.TransactionTimestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	customer.PersonalCredit.TransactionTimestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	voucherifyCustomer := mapVoucherify(customer)

	customerBytes, err := json.Marshal(voucherifyCustomer)
	if err != nil {
		return models.Customer{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(customerBytes))
	if err != nil {
		return models.Customer{}, err
	}

	req.Header.Set("X-App-Id", h.Env.VoucherifyId)
	req.Header.Set("X-App-Token", h.Env.VoucherifySecretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return models.Customer{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Customer{}, err
	}
	var VoucherifyCustomer voucherify.Customer
	if err := json.Unmarshal(body, &VoucherifyCustomer); err != nil {
		return models.Customer{}, err
	}

	customer = mapCustomer(VoucherifyCustomer)

	return customer, nil
}

func createVoucher(h Handler, customer models.Customer) (voucherify.Voucher, error) {
	voucherCode := fmt.Sprintf("%s-MEAL", customer.Id)
	url := fmt.Sprintf("%s/vouchers/%s", h.Env.LoyaltyUrl, voucherCode)
	client := &http.Client{}
	voucherifyVoucher := voucherify.Voucher{
		Category: "MEAL_BENEFIT",
		Type:     "DISCOUNT_VOUCHER",
		Discount: voucherify.VoucherDiscount{
			Type:       "PERCENT",
			PercentOff: 100,
		},
		ValidationRules: []string{
			"val_5GTal44soXkP",
		},
	}

	voucherBytes, err := json.Marshal(voucherifyVoucher)
	if err != nil {
		return voucherify.Voucher{}, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(voucherBytes))
	if err != nil {
		return voucherify.Voucher{}, err
	}

	req.Header.Set("X-App-Id", h.Env.VoucherifyId)
	req.Header.Set("X-App-Token", h.Env.VoucherifySecretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return voucherify.Voucher{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return voucherify.Voucher{}, err
	}

	voucherifyVoucher = voucherify.Voucher{}
	if err := json.Unmarshal(body, &voucherifyVoucher); err != nil {
		return voucherify.Voucher{}, err
	}

	return voucherifyVoucher, nil
}

func publishVoucher(h Handler, voucherCode, customerId string) error {
	url := fmt.Sprintf("%s/publications", h.Env.LoyaltyUrl)
	client := &http.Client{}

	voucherifyPublication := voucherify.VoucherPublication{
		Voucher: voucherCode,
		Customer: voucherify.VoucherCustomer{
			SourceId: customerId,
		},
	}

	publicationBytes, err := json.Marshal(voucherifyPublication)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(publicationBytes))
	if err != nil {
		return err
	}

	req.Header.Set("X-App-Id", h.Env.VoucherifyId)
	req.Header.Set("X-App-Token", h.Env.VoucherifySecretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
