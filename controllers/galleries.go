package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/SuWh1/WebDevGo/context"
	"github.com/SuWh1/WebDevGo/errors"
	"github.com/SuWh1/WebDevGo/models"
	"github.com/go-chi/chi/v5"
)

type Galleries struct {
	Templates struct {
		New   Template
		Edit  Template
		Index Template
	}

	GalleryService *models.GalleryService
}

// render gallery page
func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int
		Title  string
	}
	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")

	gallery, err := g.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusSeeOther)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id")) // if we get there then we will have a page with url that have id of the gallery
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}
	// if id is valid, then find gallery by id from models package
	gallery, err := g.GalleryService.ById(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return
		}
		http.Error(w, "SMTH went wrong", http.StatusInternalServerError)
		return
	}
	// now before we want to give access to edit it, we need to cheack if user owns it or not
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to edit this gallery", http.StatusForbidden)
		return
	}

	var data struct {
		ID    int
		Title string
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	g.Templates.Edit.Execute(w, r, data)
}

// handle updatage of gallery: for not just title
func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}
	gallery, err := g.GalleryService.ById(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return
		}
		http.Error(w, "SMTH went wrong", http.StatusInternalServerError)
		return
	}
	// now before we want to give access to edit it, we need to cheack if user owns it or not
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to edit this gallery", http.StatusForbidden)
		return
	}

	gallery.Title = r.FormValue("title")
	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "SMTH went wrong", http.StatusInternalServerError)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID    int
		Title string
	}
	var data struct {
		Galleries []Gallery
	}
	user := context.User(r.Context())
	galleries, err := g.GalleryService.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "SMTH went wrong", http.StatusInternalServerError)
		return
	}
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}
	g.Templates.Index.Execute(w, r, data)
}
