package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"greenlight.isez.dev/internal/data"
	image_uploader "greenlight.isez.dev/internal/uploader"
	"greenlight.isez.dev/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, _, err := r.FormFile("poster")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	// 1. Validate MIME type BEFORE uploading
	buf := make([]byte, 512)
	file.Read(buf)
	mimeType := http.DetectContentType(buf)
	if !strings.HasPrefix(mimeType, "image/") {
		app.badRequestResponse(w, r, errors.New("poster must be an image"))
		return
	}
	file.Seek(0, io.SeekStart) // reset after reading

	// 2. Validate text fields first
	year, err := strconv.ParseInt(r.FormValue("year"), 10, 32)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	parts := strings.Split(r.FormValue("runtime"), " ")
	if len(parts) != 2 || parts[1] != "mins" {
		app.serverErrorResponse(w, r, err)
		return
	}
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	genres := app.strArrToArr(r.FormValue("genres"))

	movie := &data.Movie{
		Title:       r.FormValue("title"),
		Year:        int32(year),
		Runtime:     data.Runtime(i),
		Genres:      genres,
		Description: r.FormValue("description"),
		Poster:      "",
	}

	// Validate struct
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// 3. Only now upload, with a server-generated public ID
	publicID := fmt.Sprintf("posters/%d_%s", time.Now().UnixNano(), uuid.New().String())
	cld, ctx := image_uploader.Credentials()
	savedFile, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       publicID, // not handler.Filename
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(false),
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	movie.Poster = savedFile.SecureURL
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	/// Parse the form
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	file, _, err := r.FormFile("poster")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	var input struct {
		Title       *string
		Year        *int32
		Runtime     *data.Runtime
		Genres      *[]string
		Description *string
		Poster      *string
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}

	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		movie.Genres = *input.Genres
	}

	if input.Description != nil {
		movie.Description = *input.Description
	}

	if *input.Poster != movie.Poster {
		/// delete old poster and replace with new poster functionality goes here
		movie.Poster = *input.Poster
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, metadata, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
