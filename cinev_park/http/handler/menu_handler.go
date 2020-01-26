package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/GoGroup/Movie-and-events/comment"
	"github.com/GoGroup/Movie-and-events/controller"
	"github.com/GoGroup/Movie-and-events/event"
	"github.com/GoGroup/Movie-and-events/movie"

	"github.com/GoGroup/Movie-and-events/cinema"
	"github.com/GoGroup/Movie-and-events/hall"
	"github.com/GoGroup/Movie-and-events/model"
	"github.com/GoGroup/Movie-and-events/schedule"
)

type MenuHandler struct {
	tmpl   *template.Template
	csrv   cinema.CinemaService
	hsrv   hall.HallService
	ssrv   schedule.ScheduleService
	msrv   movie.MovieService
	comsrv comment.CommentService
	evsrv  event.EventService
}

func NewMenuHandler(t *template.Template, cs cinema.CinemaService, hs hall.HallService, ss schedule.ScheduleService, ms movie.MovieService, comser comment.CommentService, evs event.EventService) *MenuHandler {

	return &MenuHandler{tmpl: t, csrv: cs, hsrv: hs, ssrv: ss, msrv: ms, comsrv: comser, evsrv: evs}

}

func (m *MenuHandler) Index(w http.ResponseWriter, r *http.Request) {

	fmt.Println(m.tmpl.ExecuteTemplate(w, "index.layout", nil))

}

func (m *MenuHandler) EventList(w http.ResponseWriter, r *http.Request) {
	events, errs := m.evsrv.Events()
	if len(errs) > 0 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

	}
	tempo := struct{ Events []model.Event }{Events: events}
	fmt.Println(m.tmpl.ExecuteTemplate(w, "eventList.layout", tempo))

}
func getCode(r *http.Request, defaultCode int) (int, string) {
	p := strings.Split(r.URL.Path, "/")
	if len(p) == 1 {
		fmt.Println("in first if")
		return defaultCode, p[0]
	} else if len(p) > 1 {
		fmt.Println("..in first if")
		code, err := strconv.Atoi(p[2])
		fmt.Println(err)
		fmt.Println(p)
		fmt.Println(code)
		if err == nil {
			fmt.Println(".....in first if")
			return code, p[2]
		} else {
			fmt.Println("...........in first if")
			fmt.Println(p)
			return defaultCode, p[1]
		}
	} else {
		fmt.Println("...........in not if")
		return defaultCode, ""
	}

}

func (m *MenuHandler) EachMovieHandler(w http.ResponseWriter, r *http.Request) {

	d, stringid := getCode(r, 0)
	fmt.Println(d)
	fmt.Println(stringid)
	trailerKey := controller.GetTrailer(stringid)
	id, _ := strconv.Atoi(stringid)

	details, _, _ := controller.GetMovieDetails(id)
	details.Trailer = trailerKey

	fmt.Println(m.tmpl.ExecuteTemplate(w, "EachMovie.layout", details))
	//	fmt.Println(m.tmpl.ExecuteTemplate(w, "index.layout", nil))

}

func (m *MenuHandler) EachNowShowing(w http.ResponseWriter, r *http.Request) {
	var id int

	p := strings.Split(r.URL.Path, "/")
	if len(p) == 1 {
		fmt.Println("in first if")
		//return defaultCode, p[0]
	} else if len(p) > 1 {
		fmt.Println("..in first if")
		code, err := strconv.Atoi(p[3])

		if err == nil {
			fmt.Println(".....in first if")
			id = code
		}
	}

	if r.FormValue("comment") != "" {
		c := model.Comment{UserID: 2, UserName: "Hanna", MovieID: uint(id), Message: r.FormValue("comment")}

		m.comsrv.StoreComment(&c)
	}

	trailerKey := controller.GetTrailer(strconv.Itoa(id))
	details, _, _ := controller.GetMovieDetails(id)
	details.Trailer = trailerKey
	comments, _ := m.comsrv.RetrieveComments(uint(id))

	tempo := struct {
		Comments    []model.Comment
		MovieDetail *model.MovieDetails
	}{
		Comments:    comments,
		MovieDetail: details,
	}

	fmt.Println(m.tmpl.ExecuteTemplate(w, "EachNowShowing.layout", tempo))
	//	fmt.Println(m.tmpl.ExecuteTemplate(w, "index.layout", nil))

}
func (m *MenuHandler) Movies(w http.ResponseWriter, r *http.Request) {
	var nowshowingdetails []model.MovieDetails
	upcomingmovies, err1, err2 := controller.GetUpcomingMovies()
	fmt.Println(upcomingmovies)
	nowshowingmovies, err := m.msrv.Movies()
	for _, a := range nowshowingmovies {

		movie, _, _ := controller.GetMovieDetails(a.TmdbID)

		nowshowingdetails = append(nowshowingdetails, *movie)
	}

	if len(err) > 0 || err1 != nil || err2 != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	tempo := struct {
		Upc *model.UpcomingMovies
		Nws []model.MovieDetails
	}{
		Upc: upcomingmovies,
		Nws: nowshowingdetails,
	}

	fmt.Println(m.tmpl.ExecuteTemplate(w, "Movie.layout", tempo))

}
func (m *MenuHandler) Theaters(w http.ResponseWriter, r *http.Request) {
	var errr []error
	var NewCinemaArray []model.Cinema

	c, err := m.csrv.Cinemas()
	for _, element := range c {
		element.Halls, errr = m.hsrv.CinemaHalls(element.ID)
		NewCinemaArray = append(NewCinemaArray, element)
	}
	fmt.Println(NewCinemaArray)

	if len(err) > 0 || len(errr) > 0 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	fmt.Println(m.tmpl.ExecuteTemplate(w, "theatersList.layout", NewCinemaArray))

}

func (m *MenuHandler) TheaterSchedule(w http.ResponseWriter, r *http.Request) {
	var CName string
	var CId string
	fmt.Println("In theatesrs sdlkfsdjf")
	p := strings.Split(r.URL.Path, "/")
	if len(p) == 1 {
		fmt.Println("in first if")
		//return defaultCode, p[0]
	} else if len(p) > 1 {
		fmt.Println("..in first if")
		code, err := strconv.Atoi(p[4])
		fmt.Println(err)
		fmt.Println(p)
		fmt.Println(code)
		if err == nil {
			fmt.Println(".....in first if")
			CName = p[3]
			CId = p[4]
			//return code, p[2]
		} else {
			fmt.Println("...........in first if")
			fmt.Println(p)
			//return defaultCode, p[1]
		}
	} else {
		fmt.Println("...........in not if")
		//return defaultCode, ""
	}

	//CName := r.FormValue("cName")
	//CId, _ := strconv.Atoi(r.FormValue("cId"))
	CcId, _ := strconv.Atoi(CId)
	uCId := uint(CcId)
	H := model.HallSchedule{}
	B := model.BindedSchedule{}

	Days := []string{"Monday", "Tuesday", "Wednsday", "Thursday", "Friday", "Saturday", "Sunday"}

	for _, d := range Days {
		fmt.Println(d)
		schdls, _ := m.ssrv.HallSchedules(uCId, d)
		fmt.Println(schdls)
		for _, s := range schdls {
			mo, _, _ := controller.GetMovieDetails(s.MoviemID)
			fmt.Println(mo)

			B.PosterPath = mo.PosterPath
			B.MovieName = mo.Title
			B.Runtime = mo.RunTime
			hall, _ := m.hsrv.Hall(uint(s.HallID))
			fmt.Println("hall is", hall)
			B.HallName = hall.HallName
			B.StartTime = s.StartingTime
			B.Day = d
			B.Dimension = s.Dimension

			H.All = append(H.All, B)
		}

	}
	H.CinemaName = CName
	fmt.Println("************************************")
	fmt.Println(H)
	// if len(err) > 0 {
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// }

	fmt.Println(m.tmpl.ExecuteTemplate(w, "scheduleDisplay.layout", H))

}
