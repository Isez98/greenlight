For the live site of Greenlight click :point_right: [here](https://greenlight.isez.dev)

For the frontend repository of this project click :point_right: [here](https://github.com/Isez98/greenlight-ui)

For access to an account to try out the project, please reach out to me through LinkedIn :point_right: [here](https://www.linkedin.com/in/isacchm)

## Description

Greenlight is a web app that serves as a hobby project in which I practice implementing technologies to sharpen my skills. I have the project open source to allow other developers and recruiters (or anyone else) to view my work.

Users can register, activate their account via email, and manage a personal movie registry. Each movie entry includes a title, year of release, runtime, genres, description, and a poster image hosted on a CDN.

The backend is a RESTful API written in Go. The frontend is built with React and TypeScript.

## Features

- User registration with email-based account activation
- Token-based authentication
- Role-based permissions (`movies:read`, `movies:write`)
- Full movie CRUD with multipart form support
- Poster image upload and management via Cloudinary CDN
- Per-IP rate limiting
- CORS support

## API Endpoints

| Method | Path | Description | Auth Required |
|---|---|---|---|
| GET | `/v1/healthcheck` | Server health | No |
| GET | `/v1/movies` | List movies | No |
| POST | `/v1/movies` | Create a movie | Yes (`movies:write`) |
| GET | `/v1/movies/:id` | Get a movie | No |
| PATCH | `/v1/movies/:id` | Update a movie | Yes (`movies:write`) |
| DELETE | `/v1/movies/:id` | Delete a movie | Yes (`movies:write`) |
| POST | `/v1/users` | Register a user | No |
| PUT | `/v1/users/activated` | Activate account | No |
| GET | `/v1/users` | Get current user info | Yes |
| POST | `/v1/tokens/authentication` | Login | No |
| GET | `/v1/tokens/verify` | Verify token | No |

## Tech Stack

#### Backend

- **Go 1.23**
  - [httprouter](https://github.com/julienschmidt/httprouter) — HTTP routing
  - [lib/pq](https://github.com/lib/pq) — PostgreSQL driver
  - [cloudinary-go](https://github.com/cloudinary/cloudinary-go) — CDN image upload
  - [go-mail](https://github.com/go-mail/mail) — Transactional email
  - [realip](https://github.com/tomasen/realip) — Real IP extraction
  - [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) — Password hashing
  - [golang.org/x/time](https://pkg.go.dev/golang.org/x/time) — Rate limiting
- **PostgreSQL** — Primary database
- **Caddy** — Reverse proxy (production)

#### Frontend

- React + TypeScript
- Tailwind CSS
- Formik
- React Router
- React Query

## Prerequisites

The following environment variables must be set (e.g. in `~/.profile` or a `.envrc` file):

```shell
GREENLIGHT_DB_DSN=<your_postgresql_dsn>
CLOUDINARY_URL=<your_cloudinary_url>
SMTP_USERNAME=<your_smtp_username>
SMTP_PASSWORD=<your_smtp_password>
```
