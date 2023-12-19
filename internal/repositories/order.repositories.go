package repositories

import (
	"gilangrizaltin/Backend_Golang/internal/models"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	*sqlx.DB
}

type IOrderRepository interface {
	RepositoryCreateOrder(client *sqlx.Tx, Order_Id string, User_Id int, dataOrder *models.OrderDetailModel, paymentUrl string) error
	RepositoryUpdateSeatSchedule(client *sqlx.Tx, dataOrder *models.OrderDetailModel) (int64, error)
	RepositoryGetScheduleDetail(scheduleId int) ([]models.ScheduleDetail, error)
	Begin() (*sqlx.Tx, error)
}

func InitializeOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) RepositoryGetOrderByID() {}

func (r *OrderRepository) RepositoryCreateOrder(client *sqlx.Tx, Order_Id string, User_Id int, dataOrder *models.OrderDetailModel, paymentUrl string) error {
	query := `insert into order_transaction (id, user_id, schedules_id, seats, total_ticket, total_purchase, paid, activate_until, payment_link)
	values (:id, :user_id, :schedules_id, :seats, :total_ticket, :total_purchase, false, :activate_until, :payment_link)`
	params := make(map[string]interface{})
	params["id"] = Order_Id
	params["user_id"] = User_Id
	params["schedules_id"] = dataOrder.Schedules
	params["seats"] = dataOrder.Seats
	params["total_ticket"] = dataOrder.Ticket
	params["total_purchase"] = dataOrder.Price_Amount
	params["activate_until"] = dataOrder.Activate
	params["payment_link"] = paymentUrl
	_, err := client.NamedExec(query, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) RepositoryGetScheduleDetail(scheduleId int) ([]models.ScheduleDetail, error) {
	data := []models.ScheduleDetail{}
	query := `SELECT
    s.id as "no",
    m.movie_name as "movie_name",
    m.small_photo_movie as "movie_photo",
    s.price_per_ticket as "price",
    to_char(s.schedule_date::timestamp at time zone 'UTC', 'YYYY-MM-DD') as "date",
    s.schedule_time as "time",
    c.cinema_name as "cinema",
    COALESCE(STRING_AGG(ot.seats, ', '), 'No data') as "seats"
FROM
    schedules s
JOIN
    cinemas c ON s.cinema_id = c.id
JOIN
    movies m ON s.movie_id = m.id
LEFT JOIN
    order_transaction ot ON s.id = ot.schedules_id
WHERE
    s.id = $1
GROUP BY
    s.id, m.movie_name, m.small_photo_movie, s.price_per_ticket, s.schedule_date, s.schedule_time, c.cinema_name, s.seat_booked;
`
	err := r.Select(&data, query, scheduleId)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *OrderRepository) RepositoryUpdateSeatSchedule(client *sqlx.Tx, dataOrder *models.OrderDetailModel) (int64, error) {
	query := `UPDATE schedules
	SET seat_booked = CONCAT(seat_booked, :seats)
	WHERE id = :schedules_id
	returning id`
	params := make(map[string]interface{})
	// params["seats"] = fmt.Sprintf(",%s", dataOrder.Seats)
	params["seats"] = dataOrder.Seats
	params["schedules_id"] = dataOrder.Schedules
	result, err := client.NamedExec(query, params)
	var dataNil int64 = 0
	if err != nil {
		return dataNil, err
	}
	data, _ := result.RowsAffected()
	return data, nil
}

func (r *OrderRepository) Begin() (*sqlx.Tx, error) {
	tx, errTx := r.Beginx()
	if errTx != nil {
		return nil, errTx
	}
	return tx, nil
}
