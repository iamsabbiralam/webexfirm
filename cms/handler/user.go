package handler

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"practice/webex/cms/paginator"
	user "practice/webex/gunk/v1/user"
)

type User struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	Password    string
	Status      int
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
	DeletedAt   time.Time
	DeletedBy   string
}

type UserTempData struct {
	CSRFField        template.HTML
	Data             []User
	FormAction       string
	FormErrors       map[string]string
	FormMessage      map[string]string
	SearchTerm       string
	PaginationData   paginator.Paginator
	PresetPermission map[string]map[string]bool
	URLs             map[string]string
	GlobalURLs       map[string]string
}

func (s *Server) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	template := s.lookupTemplate("user-list.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusSeeOther)
		return
	}

	queryString := GetQueryStringData(r, []string{}, false)
	users, err := s.user.GetAllUsers(r.Context(), &user.GetAllUserRequest{
		SearchTerm: queryString.SearchTerm,
		Limit:      limitPerPage,
		Offset:     queryString.Offset,
	})
	if err != nil {
		log.Println("unable to get list: ", err)
		http.Error(w, err.Error(), http.StatusSeeOther)
	}

	userList := []User{}
	var totalUser int32
	if users != nil {
		totalUser = users.Total
		for _, item := range users.GetUser() {
			userAppendData := User{
				ID:        item.ID,
				FirstName: item.FirstName,
				LastName:  item.LastName,
				Email:     item.Email,
				Status:    int(item.Status),
				CreatedAt: time.Time{},
			}
			userList = append(userList, userAppendData)
		}
	}

	var formMessage map[string]string
	// search message conditions
	if queryString.SearchTerm != "" && len(users.GetUser()) > 0 {
		formMessage = map[string]string{"FoundMessage": "Data Found"}
	} else if queryString.SearchTerm != "" && len(users.GetUser()) == 0 {
		formMessage = map[string]string{"NotFoundMessage": "Data Not Found"}
	}

	data := UserTempData{
		Data:        userList,
		FormMessage: formMessage,
		SearchTerm:  queryString.SearchTerm,
		GlobalURLs:  adminViewURLs(),
	}
	if len(userList) > 0 {
		data.PaginationData = paginator.NewPaginator(int32(queryString.CurrentPage), limitPerPage, totalUser, r)
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, err.Error(), http.StatusSeeOther)
	}
}
