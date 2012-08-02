package com.exercise.AndroidClient;

import java.io.IOException;
import java.io.PrintStream;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.Socket;
import java.net.SocketException;
import java.net.UnknownHostException;

import org.json.JSONException;

import android.app.Activity;
import android.app.DownloadManager;
import android.app.DownloadManager.Query;
import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.database.Cursor;
import android.os.AsyncTask;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.View;
import android.widget.EditText;


public class AndroidClient extends Activity {

	EditText textOut;
	short udpPort = 4444;
	RecvFileAsyncTask recvAsyncTask;
	DownloadManager downloadMgr;
	
	
	/** Called when the activity is first created. */
	@Override
	public void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.main);

		textOut = (EditText)findViewById(R.id.textout);
		
		downloadMgr = (DownloadManager) getSystemService(Context.DOWNLOAD_SERVICE);
	}

	// button action, opens the Downsloads app/view
	public void showDownload(View view) {
		Intent i = new Intent();
		i.setAction(DownloadManager.ACTION_VIEW_DOWNLOADS);
		startActivity(i);
	}

	// button action, start the asynch task
	public void onGo(View view) {
		if (recvAsyncTask == null) {
			recvAsyncTask = new RecvFileAsyncTask();
			recvAsyncTask.execute();
		} else {
			textOut.append("already waiting ...\n");
		}
	}
	
		 
	@Override
	public void onStart() {
		super.onStart();  // Always call the superclass method first
	}

	@Override
	public void onStop() {
		super.onStop();  // Always call the superclass method first

	}
	@Override
	public void onDestroy() {
	    super.onDestroy();  // Always call the superclass
	    
	    // probably not needed, onStop closes the socket, which should make the thread stop (?)
	    if (recvAsyncTask != null)
	    	recvAsyncTask.cancel(true);
	    
	    recvAsyncTask = null;
	}
	
	@Override
	protected void onResume() {
		super.onResume();

        // register for download mgr notifications
		registerReceiver(broadcastReceiver, 
        		new IntentFilter(DownloadManager.ACTION_DOWNLOAD_COMPLETE));
	}

	@Override
	protected void onPause() {
		super.onPause();
		unregisterReceiver(broadcastReceiver);
	}
	
	// receives download mgr notifications
    private BroadcastReceiver broadcastReceiver = new BroadcastReceiver() {
        @Override
        public void onReceive(Context context, Intent intent) {

        	String message = "Failed to understand download mgr notification !\n";
        	
            // sanity check, we register with intent ACTION_DOWNLOAD_COMPLETE
        	String action = intent.getAction();
            if (DownloadManager.ACTION_DOWNLOAD_COMPLETE.equals(action)) {
            	
                // query about the status of the download etc
            	long downloadId = intent.getLongExtra(DownloadManager.EXTRA_DOWNLOAD_ID, 0);
                Query query = new Query();
                query.setFilterById(downloadId);
                Cursor cursor = downloadMgr.query(query);
                
                if (cursor.moveToFirst()) {
                	int status = getColumnInt(cursor, DownloadManager.COLUMN_STATUS);
                	int reason = getColumnInt(cursor, DownloadManager.COLUMN_REASON);
                	
                    if (status == DownloadManager.STATUS_SUCCESSFUL) {
                    	String title = getColumnString(cursor, DownloadManager.COLUMN_TITLE); 
                    	message = "download succeeded : " + title + "\n";
                    }
                    else if (status == DownloadManager.STATUS_FAILED) {
                    	message = "download failed, reason = " + reason + "\n";
                    }
                    else {
                    	message = "download status = " + status + 
                    			  " reason = " + reason + "\n";
                    }
                }
            }

            textOut.append(message);
        }
        
        private int getColumnInt(Cursor cursor, String columnID) {
            int columnIndex = cursor.getColumnIndex(columnID);
            return cursor.getInt(columnIndex);
        }

        private String getColumnString(Cursor cursor, String columnID) {
            int columnIndex = cursor.getColumnIndex(columnID);
            return cursor.getString(columnIndex);
        }
};
	
	@Override
	public boolean onCreateOptionsMenu(Menu menu) {
		getMenuInflater().inflate(R.menu.main, menu);
		return true;
	}

	/*-----------------------------------------*/
	
	// AsyncTask that waits for a 'recvFrom' datagram (on Udp Socket) as a JSON object  
	// then connects to server on PC to send a ACK message,
	// then launches a download mgr download for the file
	private class RecvFileAsyncTask extends AsyncTask<Void, Void, Void> {

		private DatagramSocket udpSocket;
		private RecvFrom recvFrom;
		
		@Override
		protected void onPreExecute() {
			super.onPreExecute();
			
			openUdpSocket();
			
			if (udpSocket != null) {
				textOut.setText("waiting on " + udpSocket.getLocalPort() + "...\n");

				//##NB: can't read udp datagram msg 'recvFrm' here .. in main thread !
			}
		}

		//@Override
		protected Void doInBackground(Void ... voids) {
			
			// wait for a udp 'recvFrom' message
			try {
				receiveRecvFromObj();
				publishProgress(new Void[1]);
				
			} catch (IOException e1) {
				e1.printStackTrace();
				recvFrom = new RecvFrom(e1.getMessage());
				cancel(false);
			} catch (JSONException e1) {
				e1.printStackTrace();
				recvFrom = new RecvFrom(e1.getMessage());
				cancel(false);
			}

			if (!isCancelled()) 
			try {
				// do a 1st connection to server, sending an ACK string
				// after the server receives this it knows to wait for the download request
				Socket serverSocket = new Socket(recvFrom.remoteHost, recvFrom.remotePort);
				
				PrintStream outStream = new PrintStream(serverSocket.getOutputStream());
				outStream.print("ACK\n");
				
				outStream.close();
				serverSocket.close();
				
				// launch the background download 
				recvFrom.execute();
								
			} catch (UnknownHostException e) {
				e.printStackTrace();
				recvFrom.addErrMessage("\n" + e.getMessage());
			} catch (IOException e) {
				e.printStackTrace();
				recvFrom.addErrMessage("\n" + e.getMessage());
			}
			
			return null;
		}


		@Override
		protected void onProgressUpdate(Void... values) {
			super.onProgressUpdate(values);
			
			if (recvFrom.valid) { // always true here !? 
				textOut.append("received " + recvFrom.json + "\n");
				textOut.append("from : " + recvFrom.remoteHost + "\n");
			}
		}

		// onPostExecute displays the results of the AsyncTask.
		@Override
		protected void onPostExecute(Void q) {
			// isCancelled() always false here ?
			if (recvFrom != null && !isCancelled()) {
				if (recvFrom.getErrMessage() != null)
					textOut.append(recvFrom.getErrMessage() + "\n");
			}
			
			// we're done
			cleanup();
		}
		
		@Override
		protected void onCancelled(Void q) {
			// we're done
			cleanup();
		}
		
		private void openUdpSocket() {
			try {
				// open a udp socket to receive commands
//				socket = new DatagramSocket();
				udpSocket = new DatagramSocket(udpPort);
//				socket = new DatagramSocket(null);
//				socket.setReuseAddress(true);		// .. and reuse ! (so we can restart while debugging)
//				socket.bind(new InetSocketAddress(udpPort)); // void return ! cant check if worked !!!
				Log.d("", "isbound " + udpSocket.isBound());
				Log.d("", "isConnected " + udpSocket.isConnected());
				Log.d("", "on port " + udpSocket.getLocalPort());
			}
			catch (SocketException e) {
				e.printStackTrace();
				textOut.setText(e.getMessage());
			}		
		}		
	
		// receive a recvFrom object from udp datagram in a JSON object string 
		private void receiveRecvFromObj() throws IOException, JSONException {
			// Wait for a datagram to come in
			byte[] buf = new byte[1024*4];
			DatagramPacket packet = new DatagramPacket(buf, buf.length);
			udpSocket.receive(packet);

			// Create a RecvFrom with the received data
			String json = new String(packet.getData(), 0, packet.getLength());
			String remoteHost = packet.getAddress().toString();
			remoteHost = remoteHost.substring(1); // remove "/" at start of string !
			
			recvFrom = new RecvFrom(downloadMgr, json, remoteHost);
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
