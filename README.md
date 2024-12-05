![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/svdf18/44e7725b61d78d612fa0ee53b3437c78/raw/go-coverage.json)

# Whoknows (Ukendt Gruppe)

This is the Ukendt Gruppe edit of Whoknows search engine project

## How to get started

1. Fork the repository

2. Check out main (`git checkout main`).

4. Branch out using feature branches (`git checkout -b feature/[your-branch-name]`) 

5. Work on your feature to the project


## Pull requests

Feel free to open a pull request using the pull request template:

1. Ensure it works locally 

2. Make pull request to the branch you want to merge into

3. Await approval

=======

### Run Database For Development (github.com/ukendt-gruppe/wiki_scraper)
```
DOCKER COMPOSE DB:
cd docker
docker compose -f docker-compose.dev.yml up

INTERACT WITH DB (In another terminal window):
docker exec -it <db_container_name> psql -U <db_user> -d <db_name>
```

### RUN APP in Development (github.com/ukendt-gruppe/whoKnows)
```
cd src/backend
go run main.go
```

### RUN APP in Development w/ Docker (WIP)
```
cd src
docker compose -f docker-compose.dev.yml up
```
