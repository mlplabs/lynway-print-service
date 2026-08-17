package http

import (
	"context"
	"net/http"
	"strings"
)

// Language упрощенный вариант определения запрашиваемого языка
// использует первый из списка
// не определяет вес
// https://developer.mozilla.org/ru/docs/Web/HTTP/Reference/Headers/Accept-Language
func Language(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		langVal := "ru-RU"
		al := r.Header.Get("Accept-Language")

		split := strings.Split(al, ",")
		if len(split) > 0 {
			switch split[0] {
			case "en-US", "en-GB", "en":
				langVal = "en-US"
			case "ru", "ru-RU":
				langVal = "ru-RU"
			}
		}
		ctx := context.WithValue(r.Context(), "language", langVal)
		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
