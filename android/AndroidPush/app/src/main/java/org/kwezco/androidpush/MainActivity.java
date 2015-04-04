package org.kwezco.androidpush;

import android.content.Context;
import android.os.Environment;
import android.support.v7.app.ActionBarActivity;
import android.os.Bundle;
import android.view.Menu;
import android.view.MenuItem;

import java.io.File;

import go.Go;
import go.goInterface.GoInterface;

import static android.os.Environment.getExternalStoragePublicDirectory;


public class MainActivity extends ActionBarActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        // start http & mppq server
        try {
            // init Go & go lib
            initGoLib();
            GoInterface.Start();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    protected String getDir(String dirName) {
        File dir = getExternalStoragePublicDirectory(dirName);
        return dir.getPath();
    }

    protected void initGoLib() throws Exception {
        // 1st, init the Go runtime
        Context context = getApplicationContext();
        Go.init(context);

        // now init our go interface / lib
        // create/populate an init param
        GoInterface.InitParam initParam = GoInterface.NewInitParam();

        File filesDir = getFilesDir();
        initParam.setAppFilesDir(filesDir.getPath());
        initParam.setDCIM(getDir(Environment.DIRECTORY_DCIM));
        initParam.setDocuments(getDir(Environment.DIRECTORY_DOCUMENTS));
        initParam.setDownloads(getDir(Environment.DIRECTORY_DOWNLOADS));
        initParam.setMovies(getDir(Environment.DIRECTORY_MOVIES));
        initParam.setMusic(getDir(Environment.DIRECTORY_MUSIC));
        initParam.setPictures(getDir(Environment.DIRECTORY_PICTURES));
        // call Init
        GoInterface.Init(initParam);
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_main, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();

        //noinspection SimplifiableIfStatement
        if (id == R.id.action_settings) {
            return true;
        }

        return super.onOptionsItemSelected(item);
    }
}
