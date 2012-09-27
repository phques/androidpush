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
}

MarcoPolo::~MarcoPolo()
{
}


//--------------

// Send a 'marco' message (broadcast udp) looking for a 'polo' answer
// this lets us find a remote app/pier with minimum config/parameters
bool MarcoPolo::marco()
{
    //##TODO: loop on whole thing (re-send & recv, if no answer from polo
    //      ie- do recv w. timeout, re-send marco & get response etc...
    //          will need to use async_send_to for this !
    bool foundPolo = false;

    while (!foundPolo)
    {

        //--- send 'marco', broadcast UDP --
        udp::socket sock(ioService, udp::endpoint(udp::v4(), 0));
        socket_base::broadcast option(true);
        sock.set_option(option);

        udp::endpoint remoteEndpt(ip::address_v4::broadcast(), poloPort);

        string marcoMsg = this->marcoMsg();

        //##debug
        cout << "Broadcasting " << marcoMsg << " on " << sock.local_endpoint() << " to port " << poloPort << "\n";

        sock.send_to(boost::asio::buffer(marcoMsg), remoteEndpt);

        // ## TODO: loop if not the correct response (could be a different polo or just some other app)

        //-- recv polo response --
        string reply = recv(sock, response_endpoint);

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
                //##TODO: check that all digits string / convert to unsigned short
                poloResponsePort_ = responseParts[2];
                ret = true;
            }
        }
    }

    return foundPolo;
}


//--------------

// wait for a 'marco' message saking for our 'polo'
// send back 'polo' with our tcp listen socket port
bool MarcoPolo::polo(unsigned short poloListenTcpPort)
{
    // -- recv --
    udp::socket sock(ioService, udp::endpoint(udp::v4(), poloPort));

    string msg = recv(sock, response_endpoint);

    //## debug
    cout << "received  : " << msg << "\n";
    cout << "from " << response_endpoint << endl;

    // check if a marco asking for us!
    bool ok = (msg == marcoMsg());
    if (ok) {
        // -- send --
        string poloMsg = this->poloMsg(poloListenTcpPort);
        sock.send_to(boost::asio::buffer(poloMsg), response_endpoint);
    }
    else {
        //## debug
        cout << "NOT a MARCO request for us" << endl;
    }

    return ok;
}

//--------------

// Receive UDP datagram from socket sock, return as a string
// response_endpoint will hold hte endpoint of the 'caller'
string MarcoPolo::recv(udp::socket& sock, udp::endpoint& response_endpoint)
{
    char data[max_length+1];
    size_t length = sock.receive_from(
        boost::asio::buffer(data, max_length), response_endpoint);

    data[length] = 0;
    string reply(data);

    return reply;
}

// the 'marco' text message to send
string MarcoPolo::marcoMsg()
{
    stringstream msg;
    msg << "marco" << MSG_SEPAR << poloName;

    return msg.str();
}

// the 'polo' reply text message to send
string MarcoPolo::poloMsg(unsigned short poloTcpPort)
{
    stringstream msg;
    msg << "polo" << MSG_SEPAR << poloName << MSG_SEPAR << poloTcpPort;

    return msg.str();
}
