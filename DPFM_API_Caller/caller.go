package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-payment-terms-exconf-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-payment-terms-exconf-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-payment-terms-exconf-rmq-kube/database"
	"sync"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
)

type ExistenceConf struct {
	ctx context.Context
	db  *database.Mysql
	l   *logger.Logger
}

func NewExistenceConf(ctx context.Context, db *database.Mysql, l *logger.Logger) *ExistenceConf {
	return &ExistenceConf{
		ctx: ctx,
		db:  db,
		l:   l,
	}
}

func (e *ExistenceConf) Conf(input *dpfm_api_input_reader.SDC) *dpfm_api_output_formatter.PaymentTerms {
	paymentTerms := *input.PaymentTerms.PaymentTerms
	notKeyExistence := make([]string, 0, 1)
	KeyExistence := make([]string, 0, 1)

	existData := &dpfm_api_output_formatter.PaymentTerms{
		PaymentTerms:      paymentTerms,
		ExistenceConf: false,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if !e.confPaymentTerms(paymentTerms) {
			notKeyExistence = append(notKeyExistence, paymentTerms)
			return
		}
		KeyExistence = append(KeyExistence, paymentTerms)
	}()

	wg.Wait()

	if len(KeyExistence) == 0 {
		return existData
	}
	if len(notKeyExistence) > 0 {
		return existData
	}

	existData.ExistenceConf = true
	return existData
}

func (e *ExistenceConf) confPaymentTerms(val string) bool {
	rows, err := e.db.Query(
		`SELECT PaymentTerms 
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_payment_terms_payment_terms_data 
		WHERE PaymentTerms = ?;`, val,
	)
	if err != nil {
		e.l.Error(err)
		return false
	}

	for rows.Next() {
		var paymentTerms string
		err := rows.Scan(&paymentTerms)
		if err != nil {
			e.l.Error(err)
			continue
		}
		if paymentTerms == val {
			return true
		}
	}
	return false
}
