package handlers

import "log"

func (s server) logging(msg ...any) {
	log.Println(msg...)
}
