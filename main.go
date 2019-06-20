package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// Book is struct
type Book struct {
	gorm.Model
	Title string
	Price int
}

// Result is struct
type Result struct {
	Total int
}

func dbInit() {
	db, err := gorm.Open("sqlite3", "book.sqlite3")
	if err != nil {
		panic("You can't open DB (dbInit())")
	}
	defer db.Close()
	db.AutoMigrate(&Book{})
}

func dbInsert(title string, price int) {
	db, err := gorm.Open("sqlite3", "book.sqlite3")
	if err != nil {
		panic("You can't open DB (dbInsert())")
	}
	defer db.Close()
	db.Create(&Book{Title: title, Price: price})
}

func dbUpdate(id int, title string, price int) {
	db, err := gorm.Open("sqlite3", "book.sqlite3")
	if err != nil {
		panic("You can't open DB (dbUpdate())")
	}
	defer db.Close()

	var book Book
	db.First(&book, id)
	book.Title = title
	book.Price = price
	db.Save(&book)
}

func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "book.sqlite3")
	if err != nil {
		panic("You can't open DB (dbDelete())")
	}
	defer db.Close()

	var book Book
	db.First(&book, id)
	db.Unscoped().Delete(&book)
}

func dbGetAll() []Book {
	db, err := gorm.Open("sqlite3", "book.sqlite3")
	if err != nil {
		panic("You can't open DB (dbGetAll())")
	}
	defer db.Close()
	var books []Book
	db.Order("created_at desc").Find(&books)
	return books
}

func dbGetOne(id int) Book {
	db, err := gorm.Open("sqlite3", "book.sqlite3")
	if err != nil {
		panic("You can't open DB (dbGetOne())")
	}
	defer db.Close()
	var book Book
	db.First(&book, id)
	return book
}

func dbGetNum() int {
	db, err := gorm.Open("sqlite3", "book.sqlite3")
	if err != nil {
		panic("You can't open DB (dbGetNum())")
	}
	defer db.Close()
	var num int
	db.Table("books").Count(&num)
	return num
}

func dbGetPrice() Result {
	db, err := gorm.Open("sqlite3", "book.sqlite3")
	if err != nil {
		panic("You can't open DB (dbGetPrice())")
	}
	defer db.Close()
	var result Result
	db.Table("books").Select("sum(price) as total").Scan(&result)
	return result
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*.html")

	dbInit()

	// Index
	router.GET("/", func(c *gin.Context) {
		books := dbGetAll()
		num := dbGetNum()
		sumPrice := dbGetPrice()
		c.HTML(200, "index.html", gin.H{"books": books, "num": num, "sumPrice": sumPrice.Total})
	})

	// Create
	router.POST("/new", func(c *gin.Context) {
		title := c.PostForm("title")
		p := c.PostForm("price")
		price, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
		dbInsert(title, price)
		c.Redirect(302, "/")
	})

	// Edit
	router.GET("/edit/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		book := dbGetOne(id)
		c.HTML(200, "edit.html", gin.H{"book": book})
	})

	// Update
	router.POST("/update/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}

		title := c.PostForm("title")
		p := c.PostForm("price")
		price, errPrice := strconv.Atoi(p)
		if errPrice != nil {
			panic(errPrice)
		}

		dbUpdate(id, title, price)
		c.Redirect(302, "/")
	})

	// Delete
	router.POST("/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		dbDelete(id)
		c.Redirect(302, "/")
	})

	// delete_confirm
	router.GET("/delete_confirm/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		book := dbGetOne(id)
		c.HTML(200, "delete.html", gin.H{"book": book})
	})

	router.Run()
}
