![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/svdf18/44e7725b61d78d612fa0ee53b3437c78/raw/go-coverage.json)

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

# Whoknows (Ukendt Gruppe)

This is the Ukendt Gruppe edit of Whoknows search engine project

## How to get started

1. Fork the repository to your own account.

2. Clone the repository to your local machine.

3. Check out the branch you are interested in (e.g. `git checkout <branch_name>`).

4. Branch out

5. Work on your feature to the project


## Pull requests

Feel free to open a pull request:

1. Ensure it works locally 

2. Make pull request to Development branch

3. Await approval

=======
# whoKnows
