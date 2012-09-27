// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
//--marcoPolo.h --
#ifndef MARCO_H
#define MARCO_H

#include <string>
#include <boost/asio.hpp>

class MarcoPolo
{
    public:
        MarcoPolo(boost::asio::io_service& ioService, std::string poloName, unsigned short poloPort=4444);
        virtual ~MarcoPolo();

        bool marco();
        bool polo(unsigned short poloListenTcpPort);

        std::string poloResponsePort() { return poloResponsePort_; }
        boost::asio::ip::udp::endpoint poloEndpoint() { return response_endpoint; }

    private:
        std::string recv(boost::asio::ip::udp::socket& sock, boost::asio::ip::udp::endpoint& response_endpoint);
        std::string marcoMsg();
        std::string poloMsg(unsigned short poloListenTcpPort);

    private:
        boost::asio::io_service& ioService;
        std::string poloName;
        unsigned short poloPort;

        boost::asio::ip::udp::endpoint response_endpoint;
        std::string poloResponsePort_;
};


#endif // MARCO_H
