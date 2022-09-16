package main

import (
	"fmt"
	db "min-dms/database"
	"time"
)

func main() {
	fmt.Println(time.Now().Weekday())
	db.InitDb()
}
