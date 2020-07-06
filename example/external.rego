package example

default admin = false

admin {
	response := http.send({
		"method": "GET",
		# "url": "http://localhost:8080/admins"
		# NOTE: If you are running this in container...
		"url": "http://host.docker.internal:8080/admins"
	})
	response.status_code == 200
	response.body[_].name == input.name
}
