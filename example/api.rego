package example

allow {
	input.method == "GET"
	input.path = ["users", user]
	input.user == user
}

is_get {
	input.method == "GET"
}
