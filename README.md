# Wilhelmiina
## Info
### What is this?
This is supposed to be an open source clone of Wilma, the student management system used by many schools in Finland.
### But why?
I was dissatisfied with Wilma and I needed a programming project.
### How is this supposed to be better?
I'm trying to write this as really stable, also the finished product should offer real time updates for you through websockets. It also is programmed in Golang so it should be pretty fast.

This also should have an open REST API for those who are interested in programming optional applications and helper tools

## Connecting MongoDB database
In your MongoDB cluster you should find a connection string of the form\
`mongodb+srv://<username>:<password>@<cluster>.mongodb.net/<dbname>?retryWrites=true&w=majority`\
Store the required information in the <fields> in environment variables in the same order\
1. `WILHELMIINA_SERVER_USERNAME      = <username>`
2. `WILHELMIINA_SERVER_PASSWORD      = <password>`
3. `WILHELMIINA_SERVER_CLUSTER_NAME  = <cluster>`
4. `WILHELMIINA_SERVER_DATABASE_NAME = <dbname>`

## Running
1. Download golang
2. Clone repo
3. `cd wilhelmiina/src`
4. `go run .`

## Testing
`cd ./src`
### Test all files:
`go test ./database ./user ./schedule ./messages -count=1`

### For verbose test output:
`go test ./database ./user ./schedule ./messages -count=1 -v`

## Contributing
Just add a pull request lol.
