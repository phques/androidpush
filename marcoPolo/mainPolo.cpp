// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
// -- mainPolo.cpp --
#include "stdafx.h"
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

    MarcoPolo marcoPolo(io_service, "testMarcoPolo");

    unsigned short poloListenTcpPort = 0;
    std::cout << marcoPolo.polo(poloListenTcpPort) << std::endl;
  }
  catch (std::exception& e)
  {
    std::cerr << "Exception: " << e.what() << "\n";
  }

  return 0;
}
