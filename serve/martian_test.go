package serve

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	s1 := time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC)
	s2 := time.Date(2024, 1, 1, 1, 1, 20, 1, time.UTC)
	s3 := time.Date(2024, 1, 1, 1, 1, 35, 1, time.UTC)
	fmt.Println(s2.Sub(s1))
	fmt.Println(s2.Sub(s1) <= time.Second*30)
	fmt.Println(s3.Sub(s1))
	fmt.Println(s3.Sub(s1) <= time.Second*30)
}
