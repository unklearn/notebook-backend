# Unklearn Notebook Backeend

Based on a React + Golang [prototype](https://unklearn.gatsbyjs.io/prototype_two/), this repo hosts the backend code required to author a notebook that can talk to containers via docker runtime. The backend can be found [here](https://github.com/unklearn/notebook-fe)

## Development setup

Make sure you have go installed on your system. Also make sure that you have Docker-CE installed. If Docker is run as root, the backend must be run as root as well.

Run `go get` to clean up and install dependencies.

Run `go run .` to run main program from project root.
