how to build and run the entire project (it also builds all images before running):

- move to the project root directory (the directory where this README.md is located)
- run the following commands:

```bash
docker-compose up --build -d
```

to stop the project and remove all containers, networks, images and volumes created by `docker-compose up`:

```bash
docker-compose down --rmi all -v
```