# DISYS-Handin-3

## How to run the program
1. Open two terminal windows: one from the server directory and one from the client directory

2. In the server terminal run:
```
go run server.go
```

3. In the client terminal run:
```
go run client.go -name Dan
```
`Dan` can be substituted with another name. The server and client should now be connected

4. To add additional clients just open a new terminal window and run the same command as in step 3

5. To disconnect a client write `quit` in said client's terminal