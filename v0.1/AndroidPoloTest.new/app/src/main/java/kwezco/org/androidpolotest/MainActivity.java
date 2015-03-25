package kwezco.org.androidpolotest;

import android.os.AsyncTask;
import android.support.v7.app.ActionBarActivity;
import android.support.v7.app.ActionBar;
import android.support.v4.app.Fragment;
import android.os.Bundle;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.view.ViewGroup;
import android.os.Build;
import android.widget.EditText;

import java.io.IOException;
import java.net.BindException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.net.SocketException;

public class MainActivity extends ActionBarActivity {

    EditText textOut;
    short udpPort = 4444;
    RecvMarcoAsyncTask recvAsyncTask;

    // results
    String recvdMsg = "";
    String remoteHost = "";
    SocketAddress remoteHostSockAddr;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        if (savedInstanceState == null) {
            getSupportFragmentManager().beginTransaction()
                    .add(R.id.container, new PlaceholderFragment())
                    .commit();
        }
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

    // button action, start the asynch task
    public void onGo(View view) {
        textOut = (EditText)findViewById(R.id.textout);

        if (recvAsyncTask == null) {
            recvAsyncTask = new RecvMarcoAsyncTask();
            recvAsyncTask.execute();
        } else {
            textOut.append("already waiting ...\n");
        }
    }

    /**
     * A placeholder fragment containing a simple view.
     */
    public static class PlaceholderFragment extends Fragment {

        public PlaceholderFragment() {
        }

        @Override
        public View onCreateView(LayoutInflater inflater, ViewGroup container,
                                 Bundle savedInstanceState) {
            View rootView = inflater.inflate(R.layout.fragment_main, container, false);

//            textOut = (EditText)findViewById(R.id.textout);

            return rootView;
        }
    }


	/*-----------------------------------------*/

    // AsyncTask that waits for a 'recvFrom' datagram (on Udp Socket) as a JSON object
    // then connects to server on PC to send a ACK message,
    // then downloads the file
    private class RecvMarcoAsyncTask extends AsyncTask<Void, String, Void> {

        private DatagramSocket udpSocket;
        private String whichMarcoMsg;

        @Override
        protected void onPreExecute() {
            super.onPreExecute();

            openUdpSocket();

            if (udpSocket != null) {
                textOut.setText("waiting on " + udpSocket.getLocalPort() + "...\n");

                //##NB: can't read udp datagram msg here .. in main thread !
            }
        }

        //@Override
        protected Void doInBackground(Void ... voids) {

            // wait for a udp 'marco' message
            try {
                while ( !receiveMarco() )
                    ;

                String msg = "received " + recvdMsg + "\n" +
                        "from " + remoteHost + "\n";
                publishProgress(msg);

                sendPoloAnwser();

            } catch (IOException e1) {
                e1.printStackTrace();
                cancel(false);
            }

            return null;
        }


        @Override
        protected void onProgressUpdate(String... values) {
            super.onProgressUpdate(values);

            textOut.append(values[0] + "\n");
        }

        // onPostExecute displays the results of the AsyncTask.
        @Override
        protected void onPostExecute(Void q) {
            // isCancelled() always false here ?
            if (!isCancelled()) {
                textOut.append("done\n");
            }

            // we're done
            cleanup();
        }

        @Override
        protected void onCancelled(Void q) {
            // we're done
            textOut.append("onCancelled\n");

            cleanup();
        }

        private void openUdpSocket() {
            try {
                // open a udp socket to receive commands
//				udpSocket = new DatagramSocket(udpPort);
                udpSocket = new DatagramSocket(null);
                udpSocket.setReuseAddress(true);		// .. and reuse ! (so we can restart while debugging)
                udpSocket.bind(new InetSocketAddress(udpPort)); // void return ! cant check if worked !!!
                Log.d("", "isbound " + udpSocket.isBound());
                Log.d("", "isConnected " + udpSocket.isConnected());
                Log.d("", "on port " + udpSocket.getLocalPort());
            }
            catch (BindException e) {
                e.printStackTrace();
                textOut.setText(e.getMessage());
            }
            catch (SocketException e) {
                e.printStackTrace();
                textOut.setText(e.getMessage());
            }
        }

        // receive
        private boolean receiveMarco() throws IOException {
            // Wait for a datagram to come in
            byte[] buf = new byte[1024*4];
            DatagramPacket packet = new DatagramPacket(buf, buf.length);
            udpSocket.receive(packet);

            // Get the received message & address of the sender
            //recvdMsg = new String(packet.getData(), 0, packet.getLength());
            recvdMsg = new String(packet.getData(), 0, packet.getLength());
            remoteHost = packet.getAddress().toString();
            remoteHostSockAddr = packet.getSocketAddress();

            // check that it is the marco msg we expect : "marco|testMarcoPolo"
            String[] parts = recvdMsg.split("\\|");
            if (parts.length == 2 && parts[0].equals("marco")) {
                whichMarcoMsg = parts[1];
                if (!whichMarcoMsg.equals("testMarcoPolo")) {
                    publishProgress("Not the right type of marco: " + whichMarcoMsg);
                    return false;
                }

                return true;
            }

            // did not receive marco
            publishProgress("Received unknow message: " + recvdMsg.substring(0, Math.min(32, recvdMsg.length()-1)));
            return false;
        }

        private void sendPoloAnwser() throws IOException {
            int ourPort = 1234;
            String sendMsg = "polo" + "|" + whichMarcoMsg + "|" + ourPort;

            DatagramSocket sendTo = new DatagramSocket();
            DatagramPacket pack = new DatagramPacket(sendMsg.getBytes(), sendMsg.length(), remoteHostSockAddr);
            sendTo.send(pack);
        }

        private void cleanup() {
            if (udpSocket != null){
                udpSocket.close();
                udpSocket = null;
            }

            // we're done, clear out the ref to us in the parent class
            recvAsyncTask = null;
        }

    }

}
