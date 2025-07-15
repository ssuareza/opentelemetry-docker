package spinwheel

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// Run handles the spinwheel request and returns
// a random number as a string.
func Run(w http.ResponseWriter, r *http.Request) {
	roll := 1 + rand.Intn(100)

	resp := strconv.Itoa(roll)
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %vn", err)
	}
}
