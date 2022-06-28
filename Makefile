NAME=ml


all:	
	go build -o $(NAME)
	./$(NAME) -term=$(ITEM)

run:	    
	go build -o $(NAME)
	./$(NAME) -term=$(ITEM)

build:
	go build -o $(NAME)

init:
	go mod init
	go mod tidy


