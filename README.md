A bug in go's multipart handling?

When a file of size 4031 is uploaded, it triggers an unexpected EOF on the 
reading side of the request.
