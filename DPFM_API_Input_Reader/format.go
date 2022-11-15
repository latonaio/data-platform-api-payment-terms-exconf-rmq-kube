package dpfm_api_input_reader

import (
	"data-platform-api-payment-terms-exconf-rmq-kube/DPFM_API_Caller/requests"
)

func (sdc *SDC) ConvertToPaymentTerms() *requests.PaymentTerms {
	data := sdc.PaymentTerms
	return &requests.PaymentTerms{
		PaymentTerms: data.PaymentTerms,
	}
}
