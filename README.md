# crudnotes
## to run a server on port 8080

        go run .

# Examples

### Create/Update

    curl -X PUT -d "I like writing notes to myself" -H "Authorization: cole" http://localhost:8080/note/myfirstnote.txt

### Read

    curl -H "Authorization: cole" http://localhost:8080/note/myfirstnote.txt

### Delete

    curl -H "Authorization: cole" -X DELETE http://localhost:8080/note/myfirstnote.txt