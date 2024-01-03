package repositories

import (
	"fmt"
	"gilangrizaltin/Backend_Golang/internal/models"
	"log"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	*sqlx.DB
}

type IOrderRepository interface {
	RepositoryGetOrderByID(id int, page string) ([]models.GetUserOrderHistory, error)
	RepositoryCreateOrder(client *sqlx.Tx, Order_Id string, User_Id int, dataOrder *models.OrderDetailModel, paymentUrl string) error
	RepositoryOrderSuccess(orderId string) (int64, error)
	RepositoryOrderFailed(orderId string) (int64, error)
	RepositoryUpdateSeatSchedule(client *sqlx.Tx, dataOrder *models.OrderDetailModel) (int64, error)
	RepositoryGetScheduleDetail(scheduleId int) ([]models.ScheduleDetail, error)
	RepositoryCountAllOrder(id int) ([]int, error)
	Begin() (*sqlx.Tx, error)
}

func InitializeOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) RepositoryGetOrderByID(id int, page string) ([]models.GetUserOrderHistory, error) {
	data := []models.GetUserOrderHistory{}
	offset := 1
	if page != "" {
		pageData, _ := strconv.Atoi(page)
		offset = pageData
	}
	values := []any{
		id,
		(offset - 1) * 4,
	}
	query := `SELECT o.id as "no",
    m.movie_name as "movie",
	c.cinema_name as "cinema",
	COALESCE(u.first_name, '') || ' ' || COALESCE(u.last_name, '') AS "user_name",
    s.schedule_time as "time",
	to_char(s.schedule_date::timestamp at time zone 'UTC', 'YYYY-MM-DD') as "date",
    o.seats as "seats",
	o.total_ticket as "total_ticket",
	o.total_purchase as "total_purchase",
    o.activate_until as "active_until",
    o.payment_link as "va_number",
    o.status as "status"
	FROM
    order_transaction o 
	JOIN
	users u on o.user_id = u.id
	JOIN
	schedules s on o.schedules_id = s.id
	JOIN
	cinemas c on s.cinema_id = c.id
	join
	movies m on s.movie_id = m.id
	WHERE 
	u.id = $1
	AND
	o.deleted_at is null
	LIMIT 4 OFFSET $2`
	fmt.Println(offset)
	err := r.Select(&data, query, values...)
	if err != nil {
		log.Println("Error executing query:", err.Error())
		return nil, err
	}
	return data, nil
}

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

func (r *OrderRepository) RepositoryOrderSuccess(orderId string) (int64, error) {
	query := `update order_transaction set status = 'Done' where id = $1 returning id`
	result, err := r.Exec(query, orderId)
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return rows, nil
}

func (r *OrderRepository) RepositoryOrderFailed(orderId string) (int64, error) {
	query := `update order_transaction set status = 'Cancelled' where id = $1 returning id`
	result, err := r.Exec(query, orderId)
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return rows, nil
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
	c.cinema_logo as "cinema_logo",
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
		s.id, m.movie_name, m.small_photo_movie, s.price_per_ticket, s.schedule_date, s.schedule_time, c.cinema_name, s.seat_booked, c.cinema_logo
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

func (r *OrderRepository) RepositoryCountAllOrder(id int) ([]int, error) {
	var totalData = []int{}
	values := []any{
		id,
	}
	query := `
		SELECT
			COUNT(*) AS "Total_Order"
		FROM
			order_transaction o 
		WHERE 
			o.user_id = $1
		AND
			o.deleted_at is null`
	err := r.Select(&totalData, query, values...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return totalData, nil
}

func (r *OrderRepository) Begin() (*sqlx.Tx, error) {
	tx, errTx := r.Beginx()
	if errTx != nil {
		return nil, errTx
	}
	return tx, nil
}
