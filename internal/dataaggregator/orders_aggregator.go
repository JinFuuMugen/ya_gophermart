package dataaggregator

import (
	"encoding/json"
	"fmt"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/database"
	"github.com/JinFuuMugen/ya_gophermart.git/internal/models"
	"io"
	"net/http"
)

func GetOrders(user string, addr string) ([]models.Order, error) {
	orders, err := database.GetOrdersDB(user)
	if err != nil {
		return nil, fmt.Errorf("error getting orders from database: %w", err)
	}

	client := http.Client{}

	accrualOrder := make([]models.Order, 0)
	for _, o := range orders {
		resp, err := client.Get(addr + "/api/orders/" + o.Number)
		if err != nil {
			return nil, fmt.Errorf("error executing accural request: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		switch resp.StatusCode {
		case http.StatusOK:
			var orderData models.Order
			if err := json.Unmarshal(body, &orderData); err != nil {
				return nil, fmt.Errorf("error parsing JSON: %w", err)
			}

			orderData.Dateadd = o.Dateadd
			accrualOrder = append(accrualOrder, orderData)

		case http.StatusNoContent:
		default:
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	}
	return accrualOrder, nil
}
