package repositories

import (
	"fmt"
	"gilangrizaltin/Backend_Golang/internal/models"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MovieRepository struct {
	*sqlx.DB
}

type IMovieRepository interface {
	RepositoryGetAllMovie(body *models.QueryParamGetMovie) ([]models.MovieModel, error)
	RepositoryGetMovie(movieID int) ([]models.MovieModel, error)
	RepositoryGetSchedule(movieID int) ([]models.Schedule, error)
	RepositoryGetCinema(body *models.QuerySchedule) ([]models.Cinema, error)
	RepositoryAddMovie(body *models.NewMovieModel, url string, client *sqlx.Tx) (string, error)
	RepositoryAddMovieSchedule(body []models.NewMovieSchedule, client *sqlx.Tx, movieID string) error
	RepositoryEditMovie(body *models.UpdateMovieModel, movieID int, url string, client *sqlx.Tx) (int64, error)
	RepositoryDeleteMovie(movieID int) (int64, error)
	RepositoryCountAllMovie(body *models.QueryParamGetMovie) ([]int, error)
	Begin() (*sqlx.Tx, error)
}

func InitializeMovieRepository(db *sqlx.DB) *MovieRepository {
	return &MovieRepository{db}
}

func (r *MovieRepository) RepositoryGetAllMovie(body *models.QueryParamGetMovie) ([]models.MovieModel, error) {
	data := []models.MovieModel{}
	query := `SELECT m.id as "no",
    m.small_photo_movie as "movie_photo",
    m.movie_name as "movie_name",
    to_char(m.release_date::timestamp at time zone 'UTC', 'YYYY-MM-DD') as "release_date",
    m.directed_by as "director",
    m.duration as "duration",
    m.sinopsis as "sinopsis",
    m.movie_cast as "cast",
    m.genre as "genre",
    m.categories as "category"
	FROM
    movies m 
	where m.deleted_at is null `
	var conditional []string
	values := []any{}
	if body.Movie_Id != 0 {
		conditional = append(conditional, fmt.Sprintf("m.id = $%d", len(values)+1))
		values = append(values, body.Movie_Id)
	}
	if body.Movie_Name != "" {
		conditional = append(conditional, fmt.Sprintf("m.movie_name ilike $%d", len(values)+1))
		values = append(values, "%"+body.Movie_Name+"%")
	}
	if body.Movie_Genre != "" {
		conditional = append(conditional, fmt.Sprintf("m.genre ilike $%d", len(values)+1))
		values = append(values, "%"+body.Movie_Genre+"%")
	}
	if len(conditional) == 1 {
		query += fmt.Sprintf(" AND %s", conditional[0])
	}
	if len(conditional) > 1 {
		query += fmt.Sprintf(" AND %s", strings.Join(conditional, " AND "))
	}
	if body.Sort == "" {
		query += " ORDER BY m.created_at desc"
	}
	if body.Sort != "" {
		query += " ORDER BY "
		if body.Sort == "Newest" {
			query += " m.created_at desc"
		}
		if body.Sort == "Oldest" {
			query += " m.created_at asc"
		}
		if body.Sort == "A - Z" {
			query += " m.movie_name asc"
		}
		if body.Sort == "Z - A" {
			query += " m.movie_name asc"
		}
	}
	var page = body.Page
	if body.Page == 0 {
		page = 1
	}
	query += " LIMIT 8 OFFSET " + strconv.Itoa((page-1)*8)
	err := r.Select(&data, query, values...)
	if err != nil {
		log.Println("Error executing query:", err.Error())
		return nil, err
	}
	return data, nil
}

func (r *MovieRepository) RepositoryGetMovie(movieID int) ([]models.MovieModel, error) {
	data := []models.MovieModel{}
	query := `SELECT m.id as "no",
    m.small_photo_movie as "movie_photo",
    m.movie_name as "movie_name",
    to_char(m.release_date::timestamp at time zone 'UTC', 'YYYY-MM-DD') as "release_date",
    m.directed_by as "director",
    m.duration as "duration",
    m.sinopsis as "sinopsis",
    m.movie_cast as "cast",
    m.genre as "genre",
    m.categories as "category"
	FROM
    movies m
	where m.id = $1`
	err := r.Select(&data, query, movieID)
	if err != nil {
		log.Println("Error executing query:", err.Error())
		return nil, err
	}
	return data, nil
}

func (r *MovieRepository) RepositoryGetSchedule(movieID int) ([]models.Schedule, error) {
	data := []models.Schedule{}
	query := `SELECT
		s.id as "no",
	    s.price_per_ticket as "ticket_price",
	    to_char(s.schedule_date::timestamp at time zone 'UTC', 'YYYY-MM-DD') as "date",
	    s.schedule_time as "time",
	    c.cinema_name as "cinema",
		c.cinema_logo as "cinema_logo",
	    s.seat_booked as "seat"
	FROM
	    schedules s
	JOIN
	    cinemas c ON s.cinema_id = c.id
	where s.movie_id = $1`

	// data := []models.MovieDetailsSchedule{}
	// query := `SELECT
	// 	to_char(s.schedule_date::timestamp at time zone 'UTC', 'YYYY-MM-DD') as "date"
	// FROM
	// 	schedules s
	// WHERE
	// 	s.movie_id = $1
	// GROUP BY
	// 	s.schedule_date`
	err := r.Select(&data, query, movieID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *MovieRepository) RepositoryGetCinema(body *models.QuerySchedule) ([]models.Cinema, error) {
	data := []models.Cinema{}
	query := `SELECT
	s.id as "no",
    c.cinema_name as "cinema",
    s.seat_booked as "seat"
FROM
    schedules s
JOIN
    cinemas c ON s.cinema_id = c.id
where 
	s.schedule_date = $1
and 
	s.schedule_time = $2`
	values := []any{
		body.Date, body.Time,
	}
	var page = body.Page
	if body.Page == 0 {
		page = 1
	}
	query += " LIMIT 6 OFFSET " + strconv.Itoa((page-1)*6)
	err := r.Select(&data, query, values...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *MovieRepository) RepositoryAddMovie(body *models.NewMovieModel, url string, client *sqlx.Tx) (string, error) {
	query := `insert into movies (small_photo_movie, movie_name, release_date, directed_by, duration, sinopsis, movie_cast ,genre, categories) values
	(:small_photo_movie, :movie_name, :release_date, :directed_by, :duration, :sinopsis, :movie_cast ,:genre, :categories) returning id`
	params := make(map[string]interface{})
	params["small_photo_movie"] = url
	params["movie_name"] = body.Movie_Name
	params["release_date"] = body.Release_Date
	params["directed_by"] = body.Director_Name
	params["duration"] = body.Duration
	params["sinopsis"] = body.Sinopsis
	params["movie_cast"] = body.Cast
	params["genre"] = body.Genre
	params["categories"] = body.Category
	rows, err := client.NamedQuery(query, params)
	if err != nil {
		return "", err
	}
	var Id_Movie string
	var movieID string
	for rows.Next() {
		err = rows.Scan(&Id_Movie)
		if err != nil {
			return "", nil
		}
		if Id_Movie != "" {
			movieID = Id_Movie
			break
		}
	}
	defer rows.Close()
	return movieID, nil
}

func (r *MovieRepository) RepositoryAddMovieSchedule(body []models.NewMovieSchedule, client *sqlx.Tx, movieID string) error {
	query := `insert into schedules (movie_id, price_per_ticket, schedule_date, schedule_time, cinema_id)
	values `
	var filteredBody []string
	filterBody := make(map[string]interface{})
	filterBody["movie_id"] = movieID
	j := 1
	for i := 0; i < len(body); i++ {
		filteredBody = append(filteredBody, "(:movie_id ")

		filteredBody = append(filteredBody, fmt.Sprintf(`:price_per_ticket%d `, j))
		filterBody[fmt.Sprintf("price_per_ticket%d", j)] = body[i].Ticket_Price

		filteredBody = append(filteredBody, fmt.Sprintf(`:schedule_date%d `, j))
		filterBody[fmt.Sprintf("schedule_date%d", j)] = body[i].Date

		filteredBody = append(filteredBody, fmt.Sprintf(`:schedule_time%d `, j))
		filterBody[fmt.Sprintf("schedule_time%d", j)] = body[i].Time

		filteredBody = append(filteredBody, fmt.Sprintf(`:cinema_id%d)`, j))
		filterBody[fmt.Sprintf("cinema_id%d", j)] = body[i].Cinema
		j++
	}
	if len(filteredBody) > 0 {
		query += strings.Join(filteredBody, ", ")
	}
	// log.Println(query)
	// log.Println(filterBody)
	rows, err := client.NamedQuery(query, filterBody)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (r *MovieRepository) RepositoryEditMovie(body *models.UpdateMovieModel, movieID int, url string, client *sqlx.Tx) (int64, error) {
	var conditional []string
	query := `
        UPDATE movies
        SET `
	params := make(map[string]interface{})
	params["Id"] = movieID
	if url != "" {
		conditional = append(conditional, "small_photo_movie = :Url")
		params["Url"] = url
	}
	if body.Movie_Name != "" {
		conditional = append(conditional, "movie_name = :Movie_name")
		params["Movie_name"] = body.Movie_Name
	}
	if body.Release_Date != "" {
		conditional = append(conditional, "release_date = :Release_date")
		params["Release_date"] = body.Release_Date
	}
	if body.Director_Name != "" {
		conditional = append(conditional, "directed_by = :Director")
		params["Director"] = body.Director_Name
	}
	if body.Duration != "" {
		conditional = append(conditional, "duration = :Duration")
		params["Duration"] = body.Duration
	}
	if body.Genre != "" {
		conditional = append(conditional, "genre = :Genre")
		params["Genre"] = body.Genre
	}
	if body.Cast != "" {
		conditional = append(conditional, "movie_cast = :Movie_cast")
		params["Movie_cast"] = body.Cast
	}
	if body.Category != "" {
		conditional = append(conditional, "categories = :Category")
		params["Category"] = body.Category
	}
	if body.Sinopsis != "" {
		conditional = append(conditional, "sinopsis = :Sinopsis")
		params["Sinopsis"] = body.Sinopsis
	}
	if len(conditional) == 1 {
		query += conditional[0]
	}
	if len(conditional) > 1 {
		query += strings.Join(conditional, ", ")
	}
	query += ` ,updated_at = NOW() WHERE id = :Id`
	log.Println(query)
	result, err := client.NamedExec(query, params)
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return rows, nil
}

func (r *MovieRepository) RepositoryDeleteMovie(movieID int) (int64, error) {
	query := `update movies
	set 
	deleted_at = now ()
	where id = $1
	returning full_name`
	result, err := r.Exec(query, movieID)
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return rows, nil
}

func (r *MovieRepository) RepositoryCountAllMovie(body *models.QueryParamGetMovie) ([]int, error) {
	var totalData = []int{}
	query := `
		SELECT
			COUNT(*) AS "Total_Movie"
		FROM
			movies m 
		WHERE 
			m.deleted_at is null`
	var conditional []string
	values := []any{}
	if body.Movie_Id != 0 {
		conditional = append(conditional, fmt.Sprintf("m.id = $%d", len(values)+1))
		values = append(values, body.Movie_Id)
	}
	if body.Movie_Name != "" {
		conditional = append(conditional, fmt.Sprintf("m.movie_name ilike $%d", len(values)+1))
		values = append(values, "%"+body.Movie_Name+"%")
	}
	if body.Movie_Genre != "" {
		conditional = append(conditional, fmt.Sprintf("m.genre ilike $%d", len(values)+1))
		values = append(values, "%"+body.Movie_Genre+"%")
	}
	if len(conditional) == 1 {
		query += fmt.Sprintf(" AND %s", conditional[0])
	}
	if len(conditional) > 1 {
		query += fmt.Sprintf(" AND %s", strings.Join(conditional, " AND "))
	}
	err := r.Select(&totalData, query, values...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return totalData, nil
}

func (r *MovieRepository) Begin() (*sqlx.Tx, error) {
	tx, errTx := r.Beginx()
	if errTx != nil {
		return nil, errTx
	}
	return tx, nil
}
