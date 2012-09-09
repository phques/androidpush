
#include "mainWndFl.h"

int main (int argc, char ** argv)
{
    MainWindowFl* window = new MainWindowFl();
    window->show(argc, argv);
    return(Fl::run());
}

