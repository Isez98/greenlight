package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"greenlight.isez.dev/internal/data"
	"greenlight.isez.dev/internal/validator"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

type MovieForm struct {
	Title       *string
	Year        *int32
	Runtime     *data.Runtime
	Genres      *[]string
	Description *string
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown hey %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *application) readMovieForm(r *http.Request) (MovieForm, error) {
	var form MovieForm
	/// title
	if title := r.FormValue("title"); title != "" {
		form.Title = &title
	}

	if yearStr := r.FormValue("year"); yearStr != "" {
		year, err := strconv.ParseInt(yearStr, 10, 32)
		if err != nil {
			return form, fmt.Errorf("year must be an integer value")
		}
		y := int32(year)
		form.Year = &y
	}

	/// runtime
	if runtimeStr := r.FormValue("runtime"); runtimeStr != "" {
		parts := strings.Split(runtimeStr, " ")
		if len(parts) != 2 || parts[1] != "mins" {
			return form, fmt.Errorf("runtime must be in the format '<number> mins'")
		}
		i, err := strconv.ParseInt(parts[0], 10, 32)
		if err != nil {
			return form, fmt.Errorf("runtime must be an integer value")
		}
		r := data.Runtime(i)
		form.Runtime = &r
	}
	/// genres
	if genresStr := r.FormValue("genres"); genresStr != "" {
		genres := app.strArrToArr(genresStr)
		form.Genres = &genres
	}
	/// description
	if description := r.FormValue("description"); description != "" {
		form.Description = &description
	}
	return form, nil
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an interger value")
		return defaultValue
	}

	return i
}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		fn()
	}()
}

func (app *application) strArrToArr(str string) []string {
	chars := []string{"]", "^", "\\\\", "[", "\"", "(", ")", "-"}
	r := strings.Join(chars, "")
	re := regexp.MustCompile("[" + r + "]+")
	str = re.ReplaceAllString(str, "")
	return strings.Split(str, ",")
}
