package petclinic.authz

default allow = false

allow {
	input.method = "GET"
	input.path = ["pets", name]
	allowed[pet]
	pet.name = name
}

# True if either one the conditions below match.

# READ: Returns true if the pet owner matches the input user.
allowed[pet] {
	pet = data.pets[_]
	pet.owner = input.subject.user
}

# READ: Returns true if the veterinarian matches the input user and
# clinic matches the user location.
allowed[pet] {
	pet = data.pets[_]
	pet.veterinarian = input.subject.user
	pet.clinic = input.subject.location
}
