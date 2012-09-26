#ifndef MARCO_H
#define MARCO_H

#include <string>
#include <boost/asio.hpp>

class MarcoPolo
{
    public:
        MarcoPolo(boost::asio::io_service& ioService, unsigned short poloPort=4444);
        virtual ~MarcoPolo();

        bool marco();
        bool polo();

    private:
        std::string recv(boost::asio::ip::udp::socket& sock);

    private:
        boost::asio::io_service& ioService;
        boost::asio::ip::udp::endpoint response_endpoint;
        unsigned short poloPort;
};


#endif // MARCO_H
