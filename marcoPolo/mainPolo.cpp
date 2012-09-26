//#include <cstdlib>
//#include <cstring>
//#include <iostream>
#include <boost/asio.hpp>
#include "marcoPolo.h"

using namespace boost::asio;

using boost::asio::ip::udp;



int main(int argc, char* argv[])
{
  try
  {
//    if (argc != 3)
//    {
//      std::cerr << "Usage: blocking_udp_echo_client <host> <port>\n";
//      return 1;
//    }

    io_service io_service;

    MarcoPolo marcoPolo(io_service);

    std::cout << marcoPolo.polo() << std::endl;
  }
  catch (std::exception& e)
  {
    std::cerr << "Exception: " << e.what() << "\n";
  }

  return 0;
}
