RDS_Alma_Tools
=======

RDS_Alma_Tools is small server built with the Golang Echo framework that automates tasks using the Alma API for the Resource Description Services Department at UO Libraries. Currently, an automated withdraw process is under development.

### Usage

Endpoints:
  - the export endpoint returns a list of item records with a subset of fields for a given set.
    - `/withdraw/export/:id`
      - param: alma set id
      - returns a tsv
  - the upload form endpoint allows user to upload a previously exported list of item records and choose a withdraw option. Once the user submits the form, the withdraw process is launched.
    - `/withdraw/set.html`
      - returns the report name that will be generated as the process runs.

### Local development

- clone repo
- set up the .env. See the docker-compose for current variables needed by the system to run all of the supported processes.

Running directly on local system
- install go
- run `go get <package>` for the packages in go.mod OR wait for the system to tell you what needs to be installed in the next step
- `go run main.go`

Docker
- in one terminal: `docker-compose up`
- in another terminal: `docker-compose exec server bash`
- then: `go run main.go`

NOTE for local development: connecting to the UO ArchivesSpace API requires being on the VPN.

### Staging:
Podman (required) 
Staging instance must be set up to use traefik.
- `podman compose -f docker-compose.staging.yml up`
- exec into container; `go run main.go`

