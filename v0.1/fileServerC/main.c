// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>

#include <sys/stat.h>

void error(const char *msg)
{
    perror(msg);
    exit(1);
}

int main(int argc, char *argv[])
{
    FILE* file = 0;
    int sockfd, newsockfd, portno, n;
    struct sockaddr_in serv_addr, cli_addr;
    size_t filen, total=0;
    socklen_t clilen;
    char buffer[1024*16+1];
    char* filename = 0;

    if (argc < 2)
    {
        puts("params: filename [port]");
        return -1;
    }

    filename = argv[1];
    printf("file : %s\n", filename);

    {
        struct stat _stat = {0};
        stat(filename, &_stat);
        printf("size=%ld\n", _stat.st_size);
    }
    file = fopen(filename, "rb");
    if (file == NULL)
        error("failed to open file");

    if (argc == 3)
        portno = atoi(argv[2]);
    else
        portno = 8888;

    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd < 0)
        error("ERROR opening socket");
    bzero((char *) &serv_addr, sizeof(serv_addr));

    serv_addr.sin_family = AF_INET;
    serv_addr.sin_addr.s_addr = INADDR_ANY;
    serv_addr.sin_port = htons(portno);
    if (bind(sockfd, (struct sockaddr *) &serv_addr,
             sizeof(serv_addr)) < 0)
        error("ERROR on binding");
    listen(sockfd,5);
    clilen = sizeof(cli_addr);
    newsockfd = accept(sockfd,
                       (struct sockaddr *) &cli_addr,
                       &clilen);
    if (newsockfd < 0)
        error("ERROR on accept");

    do
    {

        bzero(buffer,1024*16);
        filen=fread(buffer, 1, 1024*16, file);
        n = write(newsockfd, buffer, filen);
        if (n < 0) error("ERROR writing to socket");
        if (filen < 1024*16) printf("\nfilen %d\n", filen);
        total += filen;
        printf("%d\r", total);
    }
    while (filen > 0);

    close(newsockfd);
    close(sockfd);
    fclose(file);

    puts("");

    return 0;
}
