package utils

import (
	"context"
	"log"
	"time"

)

func RunBackground(label string, timeout time.Duration, fn func(ctx context.Context) error) {
	go func ()  {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[ERROR] Panic pada background task %q: %v", label, r)
			}
		}()
	
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		done := make(chan error, 1)
		go func ()  {
			done <- fn(ctx)
		} ()

		select	{
		case err := <-done:
			if err != nil {
				log.Printf("[WARN] background task %q gagal: %v", label, err)
			}
		case <- ctx.Done():
			log.Printf("[WARN] background task %q timeout setelah %s", label, timeout)
		}
		
	} ()
}