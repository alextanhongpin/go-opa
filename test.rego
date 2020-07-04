package test

allow {
	input.method == "GET"
	input.path = ["users", user]
	user == "john"
}
