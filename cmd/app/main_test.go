package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	// tastCases := []string{
	// 	{"123"},
	// 	{},
	// 	{"232", "2232"},
	// 	{0, "2232"},
	// }

	// handler := http.HandlerFunc(refreshKeysJson)
	// for _, ts := range tastCases {
	// 	var keys sync.Map
	// 	keys.Store("arr", ts)

	// 	t.Run(ts.name, func(t *testing.T) {
	// 		rec := httptest.NewRecorder()
	// 		req, _ := http.NewRequest("GET", fmt.Springf("/keys"), nil)

	// 		handler.ServeHTTP(rec, req)
	// 		assert.Equel(t, tc.want, rec.Body.Bytes())
	// 	})

	// 	keys.Delete("arr")
	// }
}
