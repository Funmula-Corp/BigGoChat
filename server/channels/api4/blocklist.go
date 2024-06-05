package api4

import "net/http"

func (api *API) InitBlocklist() {
	api.BaseRoutes.User.Handle("/blockuser", api.APISessionRequired(getUserBlockUsers)).Methods("GET")
	api.BaseRoutes.User.Handle("/blockuser/{user_id:[A-za-z0-9]+}", api.APISessionRequired(addUserBlockUser)).Methods("PUT")
	api.BaseRoutes.User.Handle("/blockuser/{user_id:[A-za-z0-9]+}", api.APISessionRequired(deleteUserBlockUser)).Methods("DELETE")

	api.BaseRoutes.Channel.Handle("/blockuser", api.APISessionRequired(getChannelBlockUsers)).Methods("GET")
	api.BaseRoutes.Channel.Handle("/blockuser/{user_id:[A-za-z0-9]+}", api.APISessionRequired(addChannelBlockUser)).Methods("PUT")
	api.BaseRoutes.Channel.Handle("/blockuser/{user_id:[A-za-z0-9]+}", api.APISessionRequired(deleteChannelBlockUser)).Methods("DELETE")
}

func addUserBlockUser(c *Context, w http.ResponseWriter, r * http.Request){

}

func deleteUserBlockUser(c *Context, w http.ResponseWriter, r * http.Request){

}

func getUserBlockUsers(c *Context, w http.ResponseWriter, r * http.Request){
}


func addChannelBlockUser(c *Context, w http.ResponseWriter, r * http.Request){

}

func deleteChannelBlockUser(c *Context, w http.ResponseWriter, r * http.Request){

}

func getChannelBlockUsers(c *Context, w http.ResponseWriter, r * http.Request){
}

