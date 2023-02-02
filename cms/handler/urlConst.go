package handler

const (
	limitPerPage int32 = 10

	// home
	homeURL = "/"

	// error path
	ErrorPath = "/error"

	// authentication route
	registrationURL = "/registration"
	loginURL        = "/login"
	logoutPath      = "/logout"

	// users path
	dashboardPath   = "/dashboard"
	getAllUsersPath = "/users"

	// circular category
	createCircularCategoryPath       = "/circular-category/create"
	circularCategoriesPath           = "/circular-categories"
	updateCircularCategoryPath       = "/circular-category/update/{id}"
	deleteCircularCategoryPath       = "/circular-category/delete/{id}"
	updateCircularCategoryStatusPath = "/circular-category/update/status/{id}"
)
