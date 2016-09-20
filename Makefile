.PHONY: encrypt

encrypt:
	go run main.go encrypt --key 4f2f62e1-41e0-49e2-8da4-3a7ec511f498 --value "Hello World" --name "testKey"

encryptFile:
	go run main.go encrypt --key 4f2f62e1-41e0-49e2-8da4-3a7ec511f498 --file "main.go" --name "testKey"
