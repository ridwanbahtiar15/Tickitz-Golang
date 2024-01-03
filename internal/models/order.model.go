package models

type GetUserOrderHistory struct {
	No     string `db:"no" form:"no" json:"no" valid:"-"`
	User   string `db:"user_name" form:"user_name" json:"user_name" valid:"-"`
	Movie  string `db:"movie" form:"movie" json:"movie" valid:"_"`
	Cinema string `db:"cinema" form:"cinema" json:"cinema" valid:"_"`
	Time   string `db:"time" form:"time" json:"time" valid:"_"`
	Date   string `db:"date" form:"date" json:"date" valid:"_"`
	// Schedules    int    `db:"schedules" form:"schedules" json:"schedules" valid:"_"`
	Seats        string `db:"seats" form:"seats" json:"seats" valid:"-"`
	Ticket       int    `db:"total_ticket" form:"total_ticket" json:"total_ticket" valid:"-"`
	Price_Amount int    `db:"total_purchase" form:"total_purchase" json:"total_purchase" valid:"-"`
	Activate     string `db:"active_until" form:"active_until" json:"active_until" valid:"-"`
	VANumber     string `db:"va_number" form:"va_number" json:"va_number" valid:"-"`
	Status       string `db:"status" form:"status" json:"status" valid:"-"`
}

type OrderDetailModel struct {
	Schedules    int    `db:"schedules" form:"schedules" json:"schedules" valid:"numeric"`
	Seats        string `db:"seats" form:"seats" json:"seats" valid:"required"`
	Ticket       int    `db:"total_ticket" form:"total_ticket" json:"total_ticket" valid:"required"`
	Price_Amount int    `db:"total_purchase" form:"total_purchase" json:"total_purchase" valid:"required"`
	Activate     string `db:"active_until" form:"active_until" json:"active_until" valid:"required"`
	Payment      string `db:"payment" form:"payment" json:"payment" valid:"required"`
}

type MidtransTokenRequest struct {
	PaymentType        string `json:"payment_type"`
	TransactionDetails struct {
		OrderID  string `json:"order_id"`
		GrossAmt int64  `json:"gross_amount"`
	} `json:"transaction_details"`
}

type MidtransTokenResponse struct {
	Token string `json:"token"`
}

type MidtransPaymentLinkResponse struct {
	PaymentURL string `json:"redirect_url"`
}

type ScheduleDetail struct {
	ID           int    `db:"no" form:"no" json:"no" valid:"-"`
	Movie_Name   string `db:"movie_name" form:"movie_name" json:"movie_name" valid:"-"`
	Movie_Photo  string `db:"movie_photo" form:"movie_photo" json:"movie_photo" valid:"-"`
	Ticket_Price int    `db:"price" form:"price" json:"price" valid:"-"`
	Date         string `db:"date" form:"date" json:"date" valid:"-"`
	Time         string `db:"time" form:"time" json:"time" valid:"-"`
	Cinema       string `db:"cinema" form:"cinema" json:"cinema" valid:"-"`
	Seat         string `db:"seats" form:"seats" json:"seats" valid:"-"`
	Cinema_Logo  string `db:"cinema_logo" form:"cinema_logo" json:"cinema_logo" valid:"-"`
}
