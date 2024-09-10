package main

import (
	"database/sql"
	"fmt"

	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"tashfash.com/main/controllers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

// Schedule represents a single schedule entry
type Schedule struct {
	Time        string `json:"time"`
	Description string `json:"description"`
	ScheduledBy string `json:"scheduledBy"`
	Email       string `json:"email"`

	PhoneNumber string `json:"phoneNumber"`
}
type Folder struct {
	ID           int    `json:"id"`
	FolderName   string `json:"folder_name"`
	CollectionID int    `json:"collection_id"`
}

// Response represents the API response
type Response struct {
	HasSchedules bool       `json:"hasSchedules"`
	Schedules    []Schedule `json:"schedules,omitempty"`
}

type Promotion struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Image       string `json:"image_url"`
}

// Application represents an application entry
type Application struct {
	ID          int    `json:"id"`
	Date        string `json:"date"`
	Email       string `json:"email"`
	Phone       string `json:"phone_number"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Status      string `json:"status"`
}

type Image struct {
	ID           int    `json:"id"`
	ImageLink    string `json:"image_link"`
	CollectionID int    `json:"collection_id"`
}

type StatusCount struct {
	Pending      int           `json:"pending"`
	Accepted     int           `json:"approved"`
	Disapproved  int           `json:"disapproved"`
	TodayCount   int           `json:"today_count"`
	Applications []Application `json:"applications"`
}

// Global variable for the DB connection
var db *sql.DB

// JWT secret
var hmacSampleSecret = controllers.GetSecretKey()

func generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": username,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func verifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["email"].(string), nil
	}
	return "", err
}

func main() {
	// Database connection details
	host := "localhost"
	port := 5500
	user := "postgres"
	password := "1590"
	dbname := "TashFash"

	// Create connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection to the database
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
	}
	defer db.Close()

	// Ping the database to ensure the connection is valid
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the DB: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	r := gin.Default()

	// Serve static files
	r.Static("/static", "./frontend/static")

	// Load HTML files
	r.LoadHTMLFiles(
		"./frontend/pages/index.html",
		"./frontend/pages/apply.html",
		"./frontend/pages/step2.html",
		"./frontend/pages/upload.html",
		"./frontend/pages/collections.html",
		"./frontend/pages/admin/admin-index.html",
		"./frontend/pages/admin/dashboard/login.html",
		"./frontend/pages/admin/dashboard/dashboard.html",
		"./frontend/pages/admin/dashboard/view_application.html",
		"./frontend/pages/admin/dashboard/schedules.html",
		"./frontend/pages/admin/dashboard/picture_collection.html",
		"./frontend/pages/admin/dashboard/folder.html",
	)

	// GET Requests

	r.GET("/", LoggingMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "TashFash"})
	})
	r.GET("/apply", func(c *gin.Context) {
		c.HTML(http.StatusOK, "apply.html", gin.H{"title": "TashFash"})
	})
	r.GET("/schedule", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "schedules.html", gin.H{"title": "TashFash"})
	})
	r.GET("/dashboard", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "dashboard.html", gin.H{"title": "TashFash"})
	})
	r.GET("/view_application", func(ctx *gin.Context) {

		// Extract query parameters
		id := ctx.Query("id")
		email := ctx.Query("email")
		date := ctx.Query("date")
		phoneNumber := ctx.Query("phone_number")
		name := ctx.Query("name")
		description := ctx.Query("description")
		status := ctx.Query("status")

		// Render the template and pass the data
		ctx.HTML(http.StatusOK, "view_application.html", gin.H{
			"ID":          id,
			"Email":       email,
			"Date":        date,
			"PhoneNumber": phoneNumber,
			"Name":        name,
			"Description": description,
			"Status":      status,
		})

	})
	r.GET("/dashboardData", handleGetDashboardData)
	r.GET("/collections", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "collections.html", gin.H{"title": "TashFash"})
	})
	r.GET("/picture_collection", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "picture_collection.html", gin.H{"title": "TashFash"})
	})
	r.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{"title": "TashFash"})
	})
	r.GET("/applications", checkApplicationHandler)
	r.GET("/check-month", checkMonthHandler)
	r.GET("/check-date", checkDateHandler)
	r.GET("/get-promotions", getPromotions)
	r.GET("/get-folders", getFolders)
	r.GET("/upload-test", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "upload.html", gin.H{"title": "TashFash"})
	})
	r.GET("/folder", folderHandler)
	r.GET("/step2", step2Handler)

	// POST Requests
	r.POST("/submit-appointment", submitAppointment)
	r.POST("/upload-image", uploadImage)
	r.POST("/get_photos", getPhotos)
	r.POST("/finalize", handleFinalize)
	r.POST("/login", handleLogin)
	r.POST("/verify", handleVerify)
	r.POST("/apply-now", func(ctx *gin.Context) {
		ctx.Request.ParseForm()
		application := Application{
			Date:  ctx.Request.FormValue("date"),
			Email: ctx.Request.FormValue("email"),
			Phone: ctx.Request.FormValue("phone"),
			Name:  ctx.Request.FormValue("name"),
		}
		fmt.Println(application)

		_, err := db.Exec("INSERT INTO applications (date, email, phone_number, name) VALUES ($1, $2, $3, $4)",
			application.Date, application.Email, application.Phone, application.Name)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert application"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Application submitted successfully"})
	})
	r.POST("/create-folder", createFolder)
	r.POST("/delete-photo", deletePhoto)
	r.POST("/edit-promotion", editPromotion)
	r.POST("/upload-promotion", func(ctx *gin.Context) {
		price := ctx.PostForm("price")
		description := ctx.PostForm("description")
		image, err := ctx.FormFile("image")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}

		fmt.Println(price)
		fmt.Println(description)
		fmt.Println(image)

		rows, err := db.Query("INSERT INTO Promotion (description, price, image) VALUES ($1, $2, $3)", ctx.PostForm("price"), ctx.PostForm("description"), image)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert promotion"})
			return
		}
		defer rows.Close()

		ctx.JSON(http.StatusOK, gin.H{"message": "Promotion uploaded successfully", "price": price, "description": description})
	})

	// Start the server
	r.Run(":4000")
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

func handleVerify(c *gin.Context) {

	authToken := c.GetHeader("Authorization")

	fmt.Println(authToken)
}

func getPromotions(c *gin.Context) {

	var promoton Promotion

	rows, err := db.Query(`SELECT * FROM "Promotion" WHERE id = $1`, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query promotions"})
		return
	}
	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(&promoton.ID, &promoton.Description, &promoton.Price, &promoton.Image); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}

	}
	c.JSON(http.StatusOK, gin.H{"data": promoton})

}

func handleLogin(c *gin.Context) {

	email := c.PostForm("email")
	password := c.PostForm("password")
	fmt.Println(email, password)
	result, err := db.Query(`SELECT * FROM "USER" where email = $1 and password = $2`, email, password)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query users"})
		return
	}

	defer result.Close()
	if result.Next() {

		token, err := generateToken(c.PostForm("email"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
		c.SetCookie("token", token, 3600, "/", "", false, true)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}

}

func step2Handler(c *gin.Context) {

	date := c.Query("date")

	c.HTML(http.StatusOK, "step2.html", gin.H{"date": date})

}
func checkDateHandler(c *gin.Context) {
	// Get the date from query parameters
	date := c.Query("date")
	log.Println("Received date:", date)

	// Convert the date string to time.Time if necessary
	parsedDate, err := time.Parse("2006-1-2", date) // Expecting format YYYY-M-D
	if err != nil {
		log.Println("Error parsing date:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Query the database for schedules on the given date
	schedules, err := getSchedules(parsedDate)
	if err != nil {
		log.Println("Error retrieving schedules:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve schedules"})
		return
	}

	// Prepare the response
	resp := Response{
		HasSchedules: len(schedules) > 0,
		Schedules:    schedules,
	}

	// Send the response as JSON
	c.JSON(http.StatusOK, resp)
}

func handleFinalize(c *gin.Context) {
	c.Request.ParseForm()
	application := Application{
		Date:        c.Request.FormValue("date"),
		Email:       c.Request.FormValue("email"),
		Phone:       c.Request.FormValue("phone"),
		Description: c.Request.FormValue("description"),
		Name:        c.Request.FormValue("first_name") + " " + c.Request.FormValue("last_name"),
	}
	fmt.Println(application)
}
func checkApplicationHandler(c *gin.Context) {
	rows, err := db.Query("SELECT id, date, email, phone_number, name FROM applications")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query applications"})
		return
	}
	defer rows.Close()

	var applications []Application
	for rows.Next() {
		var app Application
		if err := rows.Scan(&app.ID, &app.Date, &app.Email, &app.Phone, &app.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		applications = append(applications, app)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred during iteration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": applications})

}

func getSchedules(date time.Time) ([]Schedule, error) {
	// Define the query and log it
	query := `SELECT time, description FROM schedules WHERE date = $1`
	log.Println("Executing query:", query, "with date:", date.Format("2006-01-02"))

	// Execute the query
	rows, err := db.Query(query, date.Format("2006-01-02"))
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	// Collect the results
	var schedules []Schedule
	for rows.Next() {
		var schedule Schedule
		if err := rows.Scan(&schedule.Time, &schedule.Description); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	// Check for errors after loop
	if err := rows.Err(); err != nil {
		log.Println("Row iteration error:", err)
		return nil, err
	}

	return schedules, nil
}

// New handler to get all schedules for a given month
func checkMonthHandler(c *gin.Context) {
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	// Convert month and year from string to integer
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	// Fetch schedules for the given month and year
	schedules, err := getSchedulesForMonth(month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve schedules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

// Function to fetch schedules for a specific month and year
func getSchedulesForMonth(month, year int) ([]Schedule, error) {
	// Calculate the first and last day of the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	// Prepare the SQL query
	query := `SELECT time, description, scheduled_by, email, phone_number FROM schedules WHERE date BETWEEN $1 AND $2`
	rows, err := db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect the results
	var schedules []Schedule
	for rows.Next() {
		var schedule Schedule
		if err := rows.Scan(&schedule.Time, &schedule.Description, &schedule.ScheduledBy, &schedule.Email, &schedule.PhoneNumber); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	// Check for errors after loop
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func editPromotion(c *gin.Context) {

	price := c.PostForm("price")
	description := c.PostForm("description")
	file, err := c.FormFile("image")

	fmt.Println(price, description, file.Filename)

	if err != nil {
		fmt.Errorf("Failed to process image")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
		return
	}

	source, err := file.Open()
	if err != nil {
		fmt.Errorf("Failed to open image")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image", "file": file.Filename, "db_returns": err})
		return
	}

	defer source.Close()

	//GETTING THE CURRENT WORKING DIRECTORY

	myDIr, err2 := os.Getwd()

	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println(myDIr)

	destinationDirectory := filepath.Join(myDIr, "frontend/static/images")

	if err := os.MkdirAll(destinationDirectory, os.ModePerm); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
	}

	destination := filepath.Join(destinationDirectory, file.Filename)

	if err := c.SaveUploadedFile(file, destination); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	db_image := "/static/images/" + file.Filename
	_, storage_error := db.Exec(
		`UPDATE "Promotion" 
		 SET price = $1, description = $2, image_url = $3 
		 WHERE id = 1`,
		price, description, db_image)
	if storage_error != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image", "file": file.Filename, "db_returns": storage_error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully", "file": file.Filename, "image_url": destination})

	fmt.Println(file)

}

func getPhotos(c *gin.Context) {

	collection_id := c.PostForm("collection_id")

	rows, err := db.Query(`SELECT * FROM images WHERE collection_id = $1`, collection_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch photos",
		})
		return
	}
	defer rows.Close()

	var photos []Image
	for rows.Next() {
		var photo Image
		err := rows.Scan(&photo.ID, &photo.ImageLink, &photo.CollectionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to scan rows",
			})
			return
		}
		photos = append(photos, photo)
	}

	c.JSON(http.StatusOK, gin.H{"photos": photos})

}

func getFolders(c *gin.Context) {
	rows, err := db.Query(`SELECT * FROM "Folders"`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch folders",
		})
		return
	}
	defer rows.Close()

	var folders []Folder
	for rows.Next() {
		var folder Folder
		err := rows.Scan(&folder.ID, &folder.FolderName, &folder.CollectionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to scan rows",
			})
			return
		}
		folders = append(folders, folder)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error occurred during rows iteration",
		})
		return
	}

	c.JSON(http.StatusOK, folders)
}

func createFolder(c *gin.Context) {

	folder_name := c.PostForm("collection_name")
	fmt.Println(folder_name)

	_, err := db.Exec(`INSERT INTO "Folders" (folder_name, collection_id) VALUES ($1 , $2)`, folder_name, 2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create folder",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "folder created successfully",
	})

}

// folderHandler handles requests to the /folder route
func folderHandler(c *gin.Context) {
	// Get the 'collection_id' from query parameters
	collectionID := c.Query("collection_id")

	if collectionID == "" {
		// Respond with a 400 Bad Request if 'collection_id' is missing
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection_id is required"})
		return
	}

	// Handle the 'collection_id' here (e.g., fetch data, render a page, etc.)
	// For demonstration, we'll return a simple JSON response
	c.HTML(http.StatusOK, "folder.html", gin.H{"collection_id": collectionID})
}

func deletePhoto(c *gin.Context) {
	image_link := c.PostForm("image_link")
	image_id := c.PostForm("image_id")
	fmt.Println(image_id)

	// // Delete image from the database
	res, err := db.Exec(`DELETE FROM images WHERE id = $1`, image_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image from database", "image_link": image_link, "image_id": image_id})
		return
	}

	// Check if any rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve deletion status", "image_link": image_link, "image_id": image_id})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No image found with the provided ID", "image_link": image_link, "image_id": image_id})
		return
	}
	fmt.Println(rowsAffected)

	//Get working directory
	workingDir, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get working directory"})
		return
	}

	// Construct destination path
	destinationPath := filepath.Join("frontend", image_link)
	destination := filepath.Join(workingDir, destinationPath)

	// Delete the file from the filesystem
	deletionError := os.Remove(destination)
	if deletionError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image from filesystem", "image_link": image_link, "image_id": image_id})
		log.Println("Failed to delete file:", deletionError)
		return
	}

	log.Println("File successfully deleted:", destination)

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully", "image_link": image_link, "image_id": image_id})
}

func uploadImage(c *gin.Context) {

	collection_id := c.PostForm("collection_id")
	image, err := c.FormFile("image")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}
	imageName := image.Filename

	myDir, dir_err := os.Getwd()
	if dir_err != nil {
		log.Fatal(dir_err)
	}

	destinationPath := filepath.Join(myDir, "frontend/static/images")

	if err := c.SaveUploadedFile(image, filepath.Join(destinationPath, imageName)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image", "image_name": imageName})
		return
	}

	dbString := filepath.Join("/static/images", imageName)
	fmt.Println(dbString)

	_, db_err := db.Exec(`INSERT INTO images (image_link, collection_id) VALUES ($1 , $2)`, dbString, collection_id)
	if db_err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert image into database", "image_name": imageName})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully", "image_name": imageName})

}

func submitAppointment(c *gin.Context) {

	date := c.PostForm("date")
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	description := c.PostForm("description")
	first_name := c.PostForm("first_name")
	lastname := c.PostForm("last_name")
	full_name := first_name + " " + lastname

	fmt.Println(email, phone, description, first_name, lastname)

	_, err := db.Exec(`INSERT INTO applications (date, email, phone_number, name , status) VALUES ($1 , $2, $3, $4, $5)`, date, email, phone, full_name, "pending")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert application into database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "appointment submitted successfully"})
}

func handleGetDashboardData(c *gin.Context) {

	// Query for the number of pending requests
	var pendingCount int
	err := db.QueryRow("SELECT COUNT(*) FROM applications WHERE status = 'pending'").Scan(&pendingCount)
	if err != nil {
		log.Printf("Error fetching pending count: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending count"})
		return
	}

	// Query for the number of accepted requests
	var acceptedCount int
	err = db.QueryRow("SELECT COUNT(*) FROM applications WHERE status = 'accepted'").Scan(&acceptedCount)
	if err != nil {
		log.Printf("Error fetching accepted count: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accepted count"})
		return
	}
	// Get the current date
	currentDate := time.Now().Format("2006-01-02")

	// Query for the number of today's appointments
	var todayCount int
	err = db.QueryRow("SELECT COUNT(*) FROM applications WHERE date = $1", currentDate).Scan(&todayCount)
	if err != nil {
		log.Printf("Error fetching today's appointments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch today's appointments"})
		return
	}

	var applications []Application

	rows, err := db.Query("SELECT id, date, email, phone_number, name, status FROM applications")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query applications"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var app Application
		if err := rows.Scan(&app.ID, &app.Date, &app.Email, &app.Phone, &app.Name, &app.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan application"})
			return
		}
		applications = append(applications, app)
	}

	// Create the response structure
	statusCount := StatusCount{
		Pending:      pendingCount,
		Accepted:     acceptedCount,
		Disapproved:  0,
		TodayCount:   todayCount,
		Applications: applications,
	}

	// Return the result as JSON
	c.JSON(http.StatusOK, statusCount)

}
