package test

test_get_allowed {
	allow with input as {"method": "GET", "path": ["users", "john"]}
}

test_get_denied {
	not allow with input as {"method": "GET", "path": ["users", "alice"]}
}
