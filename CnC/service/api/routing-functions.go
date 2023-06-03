package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) Welcome(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	zombies, err := rt.getZombies()
	if err != nil {
		http.Error(w, err.Error(), err.Type())
		return
	}
	if errJson := json.NewEncoder(w).Encode(zombies); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}

}

func (rt *_router) Zombie(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var zombie Zombie
	errr := json.NewDecoder(r.Body).Decode(&zombie)
	if errr != nil {
		Println("E ANCHE QUI")
		http.Error(w, errr.Error(), http.StatusBadRequest)
		return
	}
	id, err := rt.createNewZombie(zombie)
	if err != nil {
		http.Error(w, err.Error(), err.Type())
		return
	}

	if errJson := json.NewEncoder(w).Encode(id); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) Action(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var action Action
	errr := json.NewDecoder(r.Body).Decode(&action)
	if errr != nil {
		http.Error(w, errr.Error(), http.StatusBadRequest)
		return
	}
	var data map[string]interface{}

	data, err := rt.zombieAction(action)
	if err != nil {
		http.Error(w, err.Error(), err.Type())
		return
	}

	if errJson := json.NewEncoder(w).Encode(data); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	id := ps.ByName("id")

	err := rt.zombieDelete(id)
	if err != nil {
		http.Error(w, err.Error(), err.Type())
		return
	}

	if errJson := json.NewEncoder(w).Encode("OK"); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}
