# data-platform-api-payment-terms-exconf-rmq-kube
data-platform-api-payment-terms-exconf-rmq-kube は、データ連携基盤において、API で 支払条件の存在性チェックを行うためのマイクロサービスです。

## 動作環境
・ OS: LinuxOS  
・ CPU: ARM/AMD/Intel  

## 存在確認先テーブル名
以下のsqlファイルに対して、支払条件の存在確認が行われます。

* data-platform-payment-terms-data.sql（データ連携基盤 支払条件 - 支払条件データ）

## existence_check.go による存在性確認
Input で取得されたファイルに基づいて、existence_check.go で、 API がコールされます。
caller.go の 以下の箇所が、指定された API をコールするソースコードです。

```
func (e *ExistenceConf) Conf(data rabbitmq.RabbitmqMessage) map[string]interface{} {
	existData := map[string]interface{}{
		"ExistenceConf": false,
	}
	input := dpfm_api_input_reader.SDC{}
	err := json.Unmarshal(data.Raw(), &input)
	if err != nil {
		return existData
	}

	conf := "PaymentTerms"
	paymentTerms := *input.PaymentTerms.PaymentTerms
	notKeyExistence := make([]string, 0, 1)
	KeyExistence := make([]string, 0, 1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	existData[conf] = paymentTerms
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

	existData["ExistenceConf"] = true
	return existData
}

```

## Input
data-platform-api-payment-terms-exconf-rmq-kube では、以下のInputファイルをRabbitMQからJSON形式で受け取ります。  

```
{
	"connection_key": "request",
	"result": true,
	"redis_key": "abcdefg",
	"api_status_code": 200,
	"runtime_session_id": "boi9ar543dg91ipdnspi099u231280ab0v8af0ew",
	"business_partner": 201,
	"filepath": "/var/lib/aion/Data/rededge_sdc/abcdef.json",
	"service_label": "PAYMENT_TERMS",
	"PaymentTerms": {
		"PaymentTerms": "JPY"
	},
	"api_schema": "DPFMPaymentTernsReads",
	"accepter": ["PaymentTerms"],
	"order_id": 1,
	"deleted": false
}

```

## Output
data-platform-api-payment-terms-exconf-rmq-kube では、[golang-logging-library-for-data-platform](https://github.com/latonaio/golang-logging-library-for-data-platform) により、Output として、RabbitMQ へのメッセージを JSON 形式で出力します。品目マスタの対象値が存在する場合 true、存在しない場合 false、を返します。"cursor" ～ "time"は、golang-logging-library-for-data-platform による 定型フォーマットの出力結果です。

```
{
	"cursor": "/Users/latona2/bitbucket/data-platform-api-payment-terms-exconf-rmq-kube/main.go#L69",
	"function": "main.dataCallProcess",
	"level": "INFO",
	"message": {
		"PaymentTerms": {
			"PaymentTerms": "0001",
			"ExistenceConf": true
		}
	},
	"runtime_session_id": "boi9ar543dg91ipdnspi099u231280ab0v8af0ew",
	"time": "2022-11-14T23:18:48+09:00"
}
```