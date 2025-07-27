how to build and run the entire project (it also builds all images before running):

- move to the project root directory (the directory where this README.md is located)
- run the following commands to stop any running containers and remove existing images and volumes, then build and run the project:

```bash
docker-compose down --rmi all -v
docker-compose up --build -d
```

to only build the project and run it in detached mode:

```bash
docker-compose up --build -d
```

to only stop the project and remove all containers, networks, images and volumes created by `docker-compose up`:

```bash
docker-compose down --rmi all -v
```