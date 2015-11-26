# Consistent Hashing
The purpose of this program is to demonstrate how consistent hashing works and to implement a simple RESTful key-value cache data store

### Server Side Cache Data Store
A simple RESTful key-value data store with the following features:

  1. PUT http://localhost:3000/keys/{key_id}/{value}

    - Request:
  
      ```http
      PUT http://localhost:3000/keys/1/foobar
      ```
    - Response: 
  
      ```http
      HTTP 200
      ```
  2. GET http://localhost:3000/keys/{key_id} 

    - Request:
    
      ```http
      GET  http://localhost:3000/keys/1/foobar
      ```
    
    - Response:
      ```json
      {
      "key" : 1,
      "value" : "foobar"
      }
      ```
  3. GET http://localhost:3000/keys

    - Request:

      ```http
      GET  http://localhost:3000/keys
      ```
    - Response:
      ```json
        [
          {
            "key" : 1,
            "value" : "foobar"
          },
          {
            "key" : 2,
            "value" : "b"
          }
        ]
        ```


Three server instances using ports 3000, 3001, and 3002

  server1.go: http://localhost:3000

  server2.go: http://localhost:3001

  server3.go: http://localhost:3002

---
### Consistent Hashing on Client Side

A consistent hashing client in GO to support PUT and GET [/keys/{key_id}] operations. 

Hashing algorithm - crc32

https://golang.org/pkg/hash/crc32/

client.go
