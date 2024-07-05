# Spy Cat Agency Managemnt System

This is a backend service for the SCA, a test task force.

# Launching

1. Git clone the repo
2. Install docker
3. Run `docker compose up`
4. Open Postman and interact with the endpoints

# Feautres

The endpoints are described in the postman file. Here is a brief list of features instead:

## Cats

A spy cat goes on missions and spyies on targets. The SCA knows the following about their agents:

- Name
- Breed
- Salary
- Years of experience

The SCA can hire, fire and give out raises to their cats.

## Missions

SCA can create a mission, add up to 3 targets and assign a cat and set the completion status. Once a cat is assigned the
mission cannot be deleted. Once the mission is complete no targets can be added to it.

## Targets

A target is part of a mission, always. It is defined by it's/their name and country. A target's notes can be updated
until it's marked as complete or the entire mission is marked as complete.

# Tech stack

- Gin (Go)
- PostgreSQL
- go-migrate for migrations
- sqlc for queries
- Migrations and SQL query code is part of the build/early run steps of the docker image
- Docker compose