package encoding

import (
	"log"

	jsoniter "github.com/json-iterator/go"

	_ "github.com/anderseknert/roast/internal/encoding"
)

// JSON returns the fastest jsoniter configuration
// It is preferred using this function instead of jsoniter.ConfigFastest directly
// as there as the init function needs to be called to register the custom types,
// which will happen automatically on import.
func JSON() jsoniter.API {
	return jsoniter.ConfigFastest
}

// JSONRoundTrip convert any value to JSON and back again.
func JSONRoundTrip(from any, to any) error {
	bs, err := jsoniter.ConfigFastest.Marshal(from)
	if err != nil {
		return err
	}

	return jsoniter.ConfigFastest.Unmarshal(bs, to)
}

// MustJSONRoundTrip convert any value to JSON and back again, exit on failure.
func MustJSONRoundTrip(from any, to any) {
	if err := JSONRoundTrip(from, to); err != nil {
		log.Fatal(err)
	}
}
