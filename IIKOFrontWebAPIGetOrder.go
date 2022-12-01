package IIKOFrontWebAPIGetOrder

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

const HOST = "http://127.0.0.1:9042"
const TIMEOUT = 9 * time.Second

type IIKOOrderInformationSmall struct {
	ID                 string      `json:"Id"`
	Amount             float64     `json:"Amount"`
	Price              float64     `json:"Price"`
	Cost               float64     `json:"Cost"`
	Deleted            bool        `json:"Deleted"`
	PrintTime          time.Time   `json:"PrintTime"`
	CookingStartTime   time.Time   `json:"CookingStartTime"`
	CookingFinishTime  time.Time   `json:"CookingFinishTime"`
	CookingTime        string      `json:"CookingTime"`
	Size               string      `json:"Size"`
	ServeTime          time.Time   `json:"ServeTime"`
	Name               string      `json:"Name"`
	Product            string      `json:"Product"`
	Comment            interface{} `json:"Comment"`
	Status             int         `json:"Status"`
	Course             int         `json:"Course"`
	Modifiers          []string    `json:"Modifiers"`
	IsCompound         bool        `json:"IsCompound"`
	PrimaryComponent   string      `json:"PrimaryComponent"`
	SecondaryComponent string      `json:"SecondaryComponent"`
	Template           string      `json:"Template"`
}
type IIKOOrderInformationFull struct {
	ID                   string      `json:"Id"`
	Number               int         `json:"Number"`
	Status               int         `json:"Status"`
	FullSum              float64     `json:"FullSum"`
	ResultSum            float64     `json:"ResultSum"`
	OriginName           interface{} `json:"OriginName"`
	IsBanquetOrder       bool        `json:"IsBanquetOrder"`
	OpenTime             time.Time   `json:"OpenTime"`
	BillTime             time.Time   `json:"BillTime"`
	ProcessedPaymentsSum float64     `json:"ProcessedPaymentsSum"`
	WaiterName           string      `json:"WaiterName"`
	CashierName          string      `json:"CashierName"`
	TableNum             int         `json:"TableNum"`
	Waiter               string      `json:"Waiter"`
	Cashier              string      `json:"Cashier"`
	Table                string      `json:"Table"`
	Guests               []struct {
		ID    string                      `json:"Id"`
		Rank  int                         `json:"Rank"`
		Name  string                      `json:"Name"`
		Items []IIKOOrderInformationSmall `json:"Items"`
	} `json:"Guests"`
	IsDeliveryOrder  bool          `json:"IsDeliveryOrder"`
	Customers        []interface{} `json:"Customers"`
	Delivery         interface{}   `json:"Delivery"`
	OrderType        string        `json:"OrderType"`
	OrderServiceType int           `json:"OrderServiceType"`
	URL              string        `json:"Url"`
}

func ConvertFullIIKOOrderInfoToSmall(input string) (itemsResult []IIKOOrderInformationSmall) {

	var items []IIKOOrderInformationSmall
	var IIKOOrderInfo []IIKOOrderInformationFull
	err := json.Unmarshal([]byte(input), &IIKOOrderInfo)
	if err != nil {
		return nil
	}
	if len(IIKOOrderInfo) > 0 {
		for i, _ := range IIKOOrderInfo {
			for ii, _ := range IIKOOrderInfo[i].Guests {
				for _, ccc := range IIKOOrderInfo[i].Guests[ii].Items {
					ccc.Product = ""
					items = append(items, ccc)
				}
			}
		}
	}
	itemsResult = items
	return itemsResult

}

func UnlockLicense(key string) {
	url := HOST + "/api/logout/" + key + ""
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		panic(errors.New("НЕ УДАЛОСЬ СОЗДАТЬ ОБЪЕКТ ЗАПРОСА ПРИ ЗАКРЫТИИ КЛЮЧА ДЛЯ IIKO FRONT ПРИ ВОЗНИКНОВЕНИИ ДАЛЬНЕЙШИХ ОШИБОК ПЕРЕЗАПУСТИТЕ КАССУ"))
	}
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		panic(errors.New("НЕ УДАЛОСЬ СДЕЛАТЬ ЗАПРОС ПРИ ЗАКРЫТИИ КЛЮЧА ДЛЯ IIKO FRONT ПРИ ВОЗНИКНОВЕНИИ ДАЛЬНЕЙШИХ ОШИБОК ПЕРЕЗАПУСТИТЕ КАССУ"))
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil || string(body) != "true" {
		panic(errors.New("НЕ УДАЛОСЬ ЗАКРЫТЬ КЛЮЧ ДЛЯ IIKO FRONT ПРИ ВОЗНИКНОВЕНИИ ДАЛЬНЕЙШИХ ОШИБОК ПЕРЕЗАПУСТИТЕ КАССУ"))
	}
}
func GetOrderInfo(showNewOrder bool) (orderInfo []IIKOOrderInformationSmall) {
	url := HOST + "/api/login/2050"
	method := "GET"

	client := &http.Client{
		Timeout: TIMEOUT,
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		panic(errors.New("НЕ УДАЛОСЬ СОЗДАТЬ ОБЪЕКТ ЗАПРОСА, ПОПРОБУЙТЕ ПРОСКАНИРОВАТЬ ЕЩЕ РАЗ"))
	}
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		panic(errors.New("НЕ УДАЛОСЬ ВЫПОЛНИТЬ ЗАПРОС НА ПОЛУЧЕНИЕ КЛЮЧА В IIKO, ПОПРОБУЙТЕ ПРОСКАНИРОВАТЬ ЕЩЕ РАЗ"))
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		panic(errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ КЛЮЧ ДЛЯ ПЛАГИНА В IIKO, ПОПРОБУЙТЕ ПРОСКАНИРОВАТЬ ЕЩЕ РАЗ"))
	}
	secretKey := strings.Replace(string(body), "\"", "", 2)
	if len(secretKey) > 0 {
		defer UnlockLicense(secretKey)
		url := ""
		if showNewOrder {
			url = HOST + "/api/orders?key=" + secretKey + "&$top=1&$orderby=Number%20desc&$filter=Status%20has%20Resto.Front.Api.V5.Data.Orders.OrderStatus%27New%27"
		} else {
			url = HOST + "/api/orders?key=" + secretKey + "&$top=1&$orderby=Number%20desc&$filter=Status%20has%20Resto.Front.Api.V5.Data.Orders.OrderStatus%27Closed%27"
		}
		method := "GET"

		client := &http.Client{
			Timeout: TIMEOUT,
		}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			panic(errors.New("НЕ УДАЛОСЬ СОЗДАТЬ ОБЪЕКТ ЗАПРОСА НА ПОЛУЧЕНИЕ ЗАКАЗА В IIKO FRONT, ПОПРОБУЙТЕ ПРОСКАНИРОВАТЬ ЕЩЕ РАЗ"))
		}
		req.Header.Add("Accept", "application/json")

		res, err := client.Do(req)
		if err != nil {
			panic(errors.New("НЕ УДАЛОСЬ ВЫПОЛНИТЬ ЗАПРОС НА ПОЛУЧЕНИЕ ЗАКАЗА В IIKO FRONT, ПОПРОБУЙТЕ ПРОСКАНИРОВАТЬ ЕЩЕ РАЗ"))
		}

		body, err := io.ReadAll(res.Body)
		_ = res.Body.Close()
		if err != nil {
			panic(errors.New("НЕ УДАЛОСЬ СЧИТАТЬ ЗАКАЗ ИЗ IIKO FRONT, ПОПРОБУЙТЕ ПРОСКАНИРОВАТЬ ЕЩЕ РАЗ"))
		}
		return ConvertFullIIKOOrderInfoToSmall(string(body))

	} else {
		panic(errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ КЛЮЧ ДЛЯ IIKO FRONT, ПОПРОБУЙТЕ ПРОСКАНИРОВАТЬ ЕЩЕ РАЗ"))
	}
}
