package dataaggregator

import (
	"encoding/json"
	"fmt"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/database"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/models"
	"github.com/go-resty/resty/v2"
)

func GetOrders(user string, addr string) ([]models.Order, error) {
	orders, err := database.GetOrdersDB(user)
	if err != nil {
		return nil, fmt.Errorf("error getting orders from database: %w", err)
	}
	r := resty.New()
	accrualOrder := make([]models.Order, 0)
	for _, o := range orders {
		resp, err := r.R().Get(addr + "/api/orders/" + o.Number)
		if err != nil {
			return nil, fmt.Errorf("error executing accural request: %w", err)
		}
		var orderData models.Order

		fmt.Println(string(resp.Body())) //todo: remove

		if err := json.Unmarshal(resp.Body(), &orderData); err != nil {
			return nil, fmt.Errorf("error parsing JSON: %w", err)
		}

		orderData.Dateadd = o.Dateadd
		accrualOrder = append(accrualOrder, orderData)

	}
	return accrualOrder, nil
}
