
# Phone Review API

This API provides endpoints to manage phone reviews, comments, and likes. It includes user authentication and role-based access control.

## Overview

### 1. Database Schema:

![SchemaERD](/screenshots/Golang%20Phone%20Review%20ERD.png)

### 2. Project Dependencies:

- Gin HTTP Framework
- GORM Library
- Swaggo Documentation Framework
- MySQL
- JWT Authentication
- godotenv

## Getting Started

To run this application locally:

1. Clone the repository:
	```bash
	git clone https://github.com/diazlp/golang-phone-review-api.git
	```
2. Install dependencies
	```bash
	go mod tidy
	```
3. Create a `.env` file with neccessary environment variables
	```bash
	API_SECRET=<your_secret>
	TOKEN_HOUR_LIFESPAN=<token_lifespan>
	```
4. Change `database configuration` on `/configs/database.go`
	```bash
	username := <your_db_username>
	password := <your_db_password>
	host   	 := <your_db_host>
	database := <project_db_name>
	```
5. Run the application:
	```bash
	go run main.go
	```
Note: Kindly run `swag init` after you make codebase changes to update `swagger API documentation`


## API Endpoints

### Authentication

- `POST /register`: Register a new user.
- `POST /login`: Log in with existing credentials.
- `POST /change-password`: Change user password.

### Phone Endpoints

- `GET /phones`: Get all phones.
- `GET /phones/:phone_id`: Get details of a specific phone.
- `GET /phones/:phone_id/reviews`: Get reviews for a specific phone.
- `POST /phones/:phone_id/reviews`: Create a review for a phone.
- `POST /phones`: Create a new phone.

### Review Endpoints

- `GET /reviews`: Get all reviews.
- `GET /reviews/:review_id/comments`: Get comments for a review.
- `GET /reviews/:review_id/likes`: Get likes for a review.
- `PUT /reviews/:review_id`: Update a review.
- `DELETE /reviews/:review_id`: Delete a review.
- `POST /reviews/:review_id/comments`: Create a comment for a review.
- `POST /reviews/:review_id/likes`: Like a review.

### Comment Endpoints

- `GET /comments`: Get all comments.
- `GET /comments/:comment_id/likes`: Get likes for a comment.
- `PUT /comments/:comment_id`: Update a comment.
- `DELETE /comments/:comment_id`: Delete a comment.
- `POST /comments/:comment_id/likes`: Like a comment.

### Like Endpoints

- `GET /likes`: Get all likes.
- `DELETE /likes/:like_id`: Delete a like.

## Authentication & Authorization

- User registration is open to everyone.
- User authentication is required for some endpoints.
- All user can use `GET` endpoints.
- Role-based access control is implemented:
  - "Guest" role: Can like reviews and comments.
  - "Writer" role: Can create and edit reviews and comments.
  - "Admin" role: Can create phone, create review and delete review.
