package runtime

func DefaultFunctions() Functions {
	return Functions{
		"random_string": func(length float64) (string, error) {
			return "random-generated-string", nil
		},
		"random_ipv4_in": func(mask string) (string, error) {
			return "random-generated-ip4", nil
		},
		"random_mac": func() (string, error) {
			return "random-generated-mac-address", nil
		},
		"uuid": func() (string, error) {
			return "random-generated-uuid", nil
		},
	}
}
