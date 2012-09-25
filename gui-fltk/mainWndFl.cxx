// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// generated by Fast Light User Interface Designer (fluid) version 1.0300

#include "mainWndFl.h"

MainWindowFl::MainWindowFl() {
  { window = new Fl_Double_Window(255, 375, "Android Push");
    window->user_data((void*)(this));
    window->align(Fl_Align(FL_ALIGN_CLIP|FL_ALIGN_INSIDE));
    { Fl_Group* o = new Fl_Group(15, 25, 224, 110, "Destination");
      o->box(FL_ENGRAVED_FRAME);
      o->labeltype(FL_EMBOSSED_LABEL);
      { Fl_Box* o = new Fl_Box(20, 31, 55, 19, "Root");
        o->labelfont(2);
        o->align(Fl_Align(FL_ALIGN_LEFT|FL_ALIGN_INSIDE));
      } // Fl_Box* o
      { rootDestChoice = new Fl_Choice(24, 49, 207, 25);
        rootDestChoice->down_box(FL_BORDER_BOX);
        rootDestChoice->labeltype(FL_NO_LABEL);
      } // Fl_Choice* rootDestChoice
      { Fl_Box* o = new Fl_Box(20, 78, 75, 19, "Subdir");
        o->labelfont(2);
        o->align(Fl_Align(FL_ALIGN_LEFT|FL_ALIGN_INSIDE));
      } // Fl_Box* o
      { subdirInput = new Fl_Input(24, 98, 207, 25);
      } // Fl_Input* subdirInput
      { Fl_Box* o = new Fl_Box(105, 30, 95, 15, "foresize");
        o->hide();
        Fl_Group::current()->resizable(o);
      } // Fl_Box* o
      o->end();
    } // Fl_Group* o
    { Fl_Group* o = new Fl_Group(15, 159, 224, 172, "Local");
      o->box(FL_ENGRAVED_FRAME);
      o->labeltype(FL_EMBOSSED_LABEL);
      { Fl_Box* o = new Fl_Box(20, 167, 54, 19, "Root");
        o->labelfont(2);
        o->align(Fl_Align(FL_ALIGN_LEFT|FL_ALIGN_INSIDE));
      } // Fl_Box* o
      { localRootInput = new Fl_Input(24, 187, 207, 25);
      } // Fl_Input* localRootInput
      { Fl_Box* o = new Fl_Box(20, 219, 54, 19, "Files");
        o->labelfont(2);
        o->align(Fl_Align(FL_ALIGN_LEFT|FL_ALIGN_INSIDE));
      } // Fl_Box* o
      { filesBrowser = new Fl_Browser(24, 239, 207, 85);
        Fl_Group::current()->resizable(filesBrowser);
      } // Fl_Browser* filesBrowser
      { Fl_Box* o = new Fl_Box(125, 165, 95, 15, "foresize");
        o->hide();
      } // Fl_Box* o
      o->end();
      Fl_Group::current()->resizable(o);
    } // Fl_Group* o
    { Fl_Group* o = new Fl_Group(10, 339, 238, 29);
      { Fl_Button* o = new Fl_Button(145, 339, 96, 28, "Push  @-2<-");
        o->labeltype(FL_EMBOSSED_LABEL);
        o->labelfont(3);
        o->labelsize(15);
        o->labelcolor((Fl_Color)27);
      } // Fl_Button* o
      { Fl_Box* o = new Fl_Box(50, 345, 50, 15, "foresize");
        o->hide();
        Fl_Group::current()->resizable(o);
      } // Fl_Box* o
      o->end();
    } // Fl_Group* o
    window->size_range(window->w(), window->h());
    window->end();
  } // Fl_Double_Window* window
}

MainWindowFl::~MainWindowFl() {
  puts("~mainWndFl");
  if (window)
    delete window;
  window = 0;
}

void MainWindowFl::show(int argc, char ** argv) {
  window->show(argc, argv);
}

void MainWindowFl::hide() {
  window->hide();
}
