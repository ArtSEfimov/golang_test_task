package main

import (
	"bytes"
	"fmt"
	"go_text_task/pkg/db"
	"math/rand"
	"time"
)

func main() {
	var buf bytes.Buffer
	dbManager := db.NewManager()

	for range 20 {
		buf.Write(randomBytes(20))
		dbManager.Create(buf.Bytes())
		buf.Reset()
	}

	fmt.Println(db.DataIndexes)

	r1, _ := dbManager.Read(1)
	fmt.Println(string(r1))
	r2, _ := dbManager.Read(2)
	fmt.Println(string(r2))

	buf.Write(randomBytes(20))
	dbManager.Update(1, buf.Bytes())
	buf.Reset()

	r, _ := dbManager.Read(1)
	fmt.Println(string(r))

	fmt.Println(db.DataIndexes)

}
func randomBytes(length int) []byte {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
		"!@#$%^&*()-_=+[]{}|;:',.<>?/`~"

	result := make([]byte, length)
	for i := range result {
		time.Sleep(2 * time.Millisecond)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		result[i] = charset[r.Intn(len(charset))]
	}
	return result
}
