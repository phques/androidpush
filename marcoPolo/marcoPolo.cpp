// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
//--marcoPolo.cpp --
#include "stdafx.h"
#include <boost/algorithm/string/classification.hpp>
#include "marcoPolo.h"

//#include <boost/algorithm/string/classification.hpp>

using namespace boost::asio;
using namespace boost::algorithm;
using namespace std;

using boost::asio::ip::udp;

enum { max_length = 1024 };
const char* MSG_SEPAR = "|";

MarcoPolo::MarcoPolo(io_service& ioService, string poloName, unsigned short poloPort)
    : ioService(ioService), poloName(poloName), poloPort(poloPort)
{
    //ctor
}

MarcoPolo::~MarcoPolo()
{
    //dtor
}

string MarcoPolo::marcoMsg()
{
    stringstream msg;
    msg << "marco" << MSG_SEPAR << poloName;

    return msg.str();
}

string MarcoPolo::poloMsg(unsigned short poloTcpPort)
{
    stringstream msg;
    msg << "polo" << MSG_SEPAR << poloName << MSG_SEPAR << poloTcpPort;

    return msg.str();
}


bool MarcoPolo::marco()
{
    //--- send --
    udp::socket sock(ioService, udp::endpoint(udp::v4(), 0));
    socket_base::broadcast option(true);
    sock.set_option(option);

    udp::endpoint remoteEndpt(ip::address_v4::broadcast(), poloPort);

    string marcoMsg = this->marcoMsg();

    //##debug
    cout << "Broadcasting " << marcoMsg << " on " << sock.local_endpoint() << " to port " << poloPort << "\n";

    sock.send_to(boost::asio::buffer(marcoMsg), remoteEndpt);

    //-- recv --
    string reply = recv(sock);

    //##debug
    cout << "Reply is: " << reply << "\n";
    cout << "from " << response_endpoint << endl;

    // decode reply !
    // split by separator MSG_SEPAR
    bool ret = false;
    vector<string> responseParts;
    split(responseParts, reply, is_any_of(MSG_SEPAR), token_compress_on);
    if (responseParts.size() == 3)
    {
        if (responseParts[0] == "polo" && responseParts[1] == poloName)
        {
            poloResponsePort_ = responseParts[2];
            ret = true;
        }
    }

    return ret;
}


bool MarcoPolo::polo(unsigned short poloListenTcpPort)
{
    // -- recv --
    udp::socket sock(ioService, udp::endpoint(udp::v4(), poloPort));

    string msg = recv(sock);

    //## debug
    cout << "received  : " << msg << "\n";
    cout << "from " << response_endpoint << endl;

    //## PQ TODO parse !
    bool ok = (msg == marcoMsg());
    if (ok) {
        // -- send --
        string poloMsg = this->poloMsg(poloListenTcpPort);
        sock.send_to(boost::asio::buffer(poloMsg), response_endpoint);
    }
    else {
        cout << "NOT a MARCO request" << endl;
    }

    return ok;
}

string MarcoPolo::recv(udp::socket& sock)
{
    char data[max_length+1];
    size_t length = sock.receive_from(
        boost::asio::buffer(data, max_length), response_endpoint);

    data[length] = 0;
    string reply(data);

    return reply;
}
