package models

type MovieModel struct {
	Id              int         `db:"no" valid:"-"`
	Movie_Photo     interface{} `db:"movie_photo" json:"movie_photo"`
	Big_Movie_Photo interface{} `db:"big_movie_photo" json:"big_movie_photo"`
	Movie_Name      string      `db:"movie_name" form:"movie_name" json:"movie_name" valid:"-"`
	Genre           string      `db:"genre" form:"genre" json:"genre" valid:"matches(^[a-zA-Z ]+$), optional"`
	Release_Date    string      `db:"release_date" form:"release_date" json:"release_date" valid:"matches(^[a-zA-Z ]+$), optional"`
	Duration        string      `db:"duration" form:"duration" json:"duration" valid:"optional"`
	Director_Name   string      `db:"director" form:"director" json:"director" valid:"optional"`
	Cast            string      `db:"cast" form:"cast" json:"cast" valid:"in(Admin|Normal User),optional"`
	Category        string      `db:"category" form:"category" json:"category" valid:"optional"`
	Sinopsis        string      `db:"sinopsis" form:"sinopsis" json:"sinopsis" valid:"-"`
}

type QueryParamGetMovie struct {
	Movie_Id    int    `db:"movie_id" form:"movie_id" json:"movie_id" valid:"optional"`
	Movie_Name  string `db:"movie_name" form:"movie_name" json:"movie_name" valid:"optional"`
	Movie_Genre string `db:"movie_genre" form:"movie_genre" json:"movie_genre" valid:"optional"`
	Sort        string `db:"sort" form:"sort" json:"sort" valid:"optional"`
	Page        int    `db:"page" form:"page" json:"page" valid:"optional"`
}

type NewMovieModel struct {
	Id              int         `db:"no" valid:"-"`
	Movie_Photo     interface{} `db:"movie_photo" json:"movie_photo"`
	Big_Movie_Photo interface{} `db:"big_movie_photo" json:"big_movie_photo"`
	Movie_Name      string      `db:"movie_name" form:"movie_name" json:"movie_name" valid:"required"`
	Genre           string      `db:"genre" form:"genre" json:"genre" valid:"in(Thriller|Adventure|Horror|Romantic|Sci fi), required"`
	Release_Date    string      `db:"release_date" form:"release_date" json:"release_date" valid:"required"`
	Duration        string      `db:"duration" form:"duration" json:"duration" valid:"required"`
	Director_Name   string      `db:"director" form:"director" json:"director" valid:"required"`
	Cast            string      `db:"cast" form:"cast" json:"cast" valid:"required"`
	Category        string      `db:"category" form:"category" json:"category" valid:"in(G|PG|PG-13|R|NC-17)"`
	Sinopsis        string      `db:"sinopsis" form:"sinopsis" json:"sinopsis" valid:"required"`
	Schedules       string      `form:"schedules" json:"schedules" valid:"required"`
}

type UpdateMovieModel struct {
	Id              int         `db:"no" valid:"-"`
	Movie_Photo     interface{} `db:"movie_photo" json:"movie_photo" valid:"optional"`
	Big_Movie_Photo interface{} `db:"big_movie_photo" json:"big_movie_photo" valid:"optional"`
	Movie_Name      string      `db:"movie_name" form:"movie_name" json:"movie_name" valid:"optional"`
	Genre           string      `db:"genre" form:"genre" json:"genre" valid:"in(Thriller|Adventure|Horror|Romantic|Sci fi), optional"`
	Release_Date    string      `db:"release_date" form:"release_date" json:"release_date" valid:"optional"`
	Duration        string      `db:"duration" form:"duration" json:"duration" valid:"optional"`
	Director_Name   string      `db:"director" form:"director" json:"director" valid:"optional"`
	Cast            string      `db:"cast" form:"cast" json:"cast" valid:"optional"`
	Category        string      `db:"category" form:"category" json:"category" valid:"in(G|PG|PG-13|R|NC-17), optional"`
	Sinopsis        string      `db:"sinopsis" form:"sinopsis" json:"sinopsis" valid:"optional"`
	Schedules       string      `form:"schedules" json:"schedules" valid:"optional"`
}

type NewMovieSchedule struct {
	Date         string `db:"date" form:"date" json:"date" binding:"required" valid:"required"`
	Ticket_Price int    `db:"ticket_price" form:"ticket_price" json:"ticket_price" valid:"numeric, required"`
	Cinema       int    `db:"cinema" form:"cinema" json:"cinema" valid:"numeric, required"`
	Time         string `db:"time" form:"time" json:"time" valid:"numeric, required"`
}

type Schedule struct {
	Date         string  `db:"date" form:"date" json:"date" valid:"-"`
	Ticket_Price int     `db:"ticket_price" form:"ticket_price" json:"ticket_price" valid:"-"`
	Cinema       string  `db:"cinema" form:"cinema" json:"cinema" valid:"-"`
	Seat         *string `db:"seat" form:"seat" json:"seat" valid:"-"`
	Time         string  `db:"time" form:"time" json:"time" valid:"numeric, optional"`
}
