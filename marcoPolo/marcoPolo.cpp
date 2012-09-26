// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
//--marcoPolo.cpp --
#include "marcoPolo.h"

using namespace boost::asio;

using boost::asio::ip::udp;

enum { max_length = 1024 };

MarcoPolo::MarcoPolo(io_service& ioService, unsigned short poloPort)
    : ioService(ioService), poloPort(poloPort)
{
    //ctor
}

MarcoPolo::~MarcoPolo()
{
    //dtor
}

bool MarcoPolo::marco()
{
    //--- send --
    udp::socket sock(ioService, udp::endpoint(udp::v4(), 0));
    socket_base::broadcast option(true);
    sock.set_option(option);

    udp::endpoint remoteEndpt(ip::address_v4::broadcast(), poloPort);

    //##debug
    std::cout << "Sending marco on " << sock.local_endpoint() << "\n";

    std::string msg = "marco";
    sock.send_to(boost::asio::buffer(msg), remoteEndpt);

    //-- recv --
    std::string reply = recv(sock);

    //##debug
    std::cout << "Reply is: " << reply << "\n";
    std::cout << "from " << response_endpoint << std::endl;

    return reply == "polo";
}


bool MarcoPolo::polo()
{
    // -- recv --
    udp::socket sock(ioService, udp::endpoint(udp::v4(), poloPort));

    std::string reply = recv(sock);

    //## debug
    std::cout << "received  : " << reply << "\n";
    std::cout << "from " << response_endpoint << std::endl;

    bool ok = (reply == "marco");
    if (ok) {
        // -- send --
        std::string polo = "polo";
        sock.send_to(boost::asio::buffer(polo), response_endpoint);
    }
    else {
        std::cout << "NOT a MARCO request" << std::endl;
    }

    return ok;
}

std::string MarcoPolo::recv(udp::socket& sock)
{
    char data[max_length+1];
    size_t length = sock.receive_from(
        boost::asio::buffer(data, max_length), response_endpoint);

    data[length] = 0;
    std::string reply(data);

    return reply;
}
