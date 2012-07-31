package com.exercise.AndroidClient;

import java.io.IOException;
import java.io.OutputStream;
import java.io.PrintStream;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.ServerSocket;
import java.net.Socket;
import java.net.SocketException;

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
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

public class AndroidClient extends Activity {

	EditText textOut;
	DownloadWebpageText downloader;
	DownloadManager downloadMgr;
	DatagramSocket socket;
	short udpPort = 4444;
	
	/** Called when the activity is first created. */
	@Override
	public void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.main);

		textOut = (EditText)findViewById(R.id.textout);
		
		Button buttonSend = (Button)findViewById(R.id.send);
		buttonSend.setOnClickListener(buttonSendOnClickListener);
		
		downloadMgr = (DownloadManager) getSystemService(Context.DOWNLOAD_SERVICE);
	}

	// opens the Downsloads app/view
	public void showDownload(View view) {
		Intent i = new Intent();
		i.setAction(DownloadManager.ACTION_VIEW_DOWNLOADS);
		startActivity(i);
	}


	private void openUdpSocket() {
		try {
			// open a udp socket to receive commands
//			socket = new DatagramSocket();
			socket = new DatagramSocket(udpPort);
//			socket = new DatagramSocket(null);
//			socket.setReuseAddress(true);		// .. and reuse ! (so we can restart while debugging)
//			socket.bind(new InetSocketAddress(udpPort)); // void return ! cant check if worked !!!
			Log.d("", "isbound " + socket.isBound());
			Log.d("", "isConnected " + socket.isConnected());
			Log.d("", "on port " + socket.getLocalPort());
		}
		catch (SocketException e) {
			e.printStackTrace();
			textOut.setText(e.getMessage());

			if (socket != null){
				socket.close();
				socket = null;
			}
		}		
	}
	
	Button.OnClickListener buttonSendOnClickListener = new Button.OnClickListener() {
		public void onClick(View arg0) {
			if (socket != null) {
				if (downloader == null) {
					textOut.setText("waiting on " + socket.getLocalPort() + "...\n");
					downloader = new DownloadWebpageText();
					downloader.execute();
				} else {
					textOut.append("already waiting ...\n");
				}
			} else {
				textOut.append("no socket !!\n");
			}
		}
	};
		 
	@Override
	public void onStart() {
		super.onStart();  // Always call the superclass method first

		openUdpSocket(); 
	}

	
	@Override
	public void onStop() {
		super.onStop();  // Always call the superclass method first
		
		if (socket != null) {
			socket.close();
			socket = null;
		}
	}
	
	@Override
	public void onDestroy() {
	    super.onDestroy();  // Always call the superclass
	    
	    // probably not needed, onStop closes the socket, which should make the thread stop (?)
	    if (downloader != null)
	    	downloader.cancel(true);
	    
	    downloader = null;
	}

	
	@Override
	protected void onResume() {
		super.onResume();

        registerReceiver(broadcastReceiver, 
        		new IntentFilter(DownloadManager.ACTION_DOWNLOAD_COMPLETE));
	}

	@Override
	protected void onPause() {
		super.onPause();
		unregisterReceiver(broadcastReceiver);
	}
	
	
    private BroadcastReceiver broadcastReceiver = new BroadcastReceiver() {
        @Override
        public void onReceive(Context context, Intent intent) {

        	String message = "Failed to understand download mgr notification !\n";
        	
            // sanity check, we register with intent ACTION_DOWNLOAD_COMPLETE
        	String action = intent.getAction();
            if (DownloadManager.ACTION_DOWNLOAD_COMPLETE.equals(action)) {
            	
                long downloadId = intent.getLongExtra(DownloadManager.EXTRA_DOWNLOAD_ID, 0);
                Query query = new Query();
                query.setFilterById(downloadId);
                Cursor cursor = downloadMgr.query(query);
                
                if (cursor.moveToFirst()) {
                	int status = getColumn(cursor, DownloadManager.COLUMN_STATUS);
                	int reason = getColumn(cursor, DownloadManager.COLUMN_REASON);
                    if (status == DownloadManager.STATUS_SUCCESSFUL) {
                    	message = "download succeeded\n";
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
        
        private int getColumn(Cursor cursor, String columnID) {
            int columnIndex = cursor.getColumnIndex(columnID);
            return cursor.getInt(columnIndex);
        }
    };
	
	@Override
	public boolean onCreateOptionsMenu(Menu menu) {
		getMenuInflater().inflate(R.menu.main, menu);
		return true;
	}

	/*-----------------------------------------*/
	
	
	private class DownloadWebpageText extends AsyncTask<Void, Void, RecvFrom> {

		//@Override
		protected RecvFrom doInBackground(Void ... voids) {

			RecvFrom recvFrom = null;
			try {
				// Wait for a datagram to come in
				byte[] buf = new byte[1024*4];
				DatagramPacket packet = new DatagramPacket(buf, buf.length);
				socket.receive(packet);

				// Create a RecvFrom with the received data
				String json = new String(packet.getData(), 0, packet.getLength());
				String remoteHost = packet.getAddress().toString();
				remoteHost = remoteHost.substring(1); // remove "/" at start of string !
				
				recvFrom = new RecvFrom(downloadMgr, json, remoteHost);
				
				// do a 1st connection to server, sending an ACK string
				// after the server receives this it knows to wait for the download request
				Socket serverSocket = new Socket(recvFrom.remoteHost, recvFrom.remotePort);
				PrintStream outStream = new PrintStream(serverSocket.getOutputStream());
				outStream.print("ACK\n");
				outStream.close();
			
				recvFrom.execute();
							
			} catch (IOException e) {
				e.printStackTrace();
				recvFrom = new RecvFrom(e.getMessage());
			} catch (JSONException e) {
				e.printStackTrace();
				recvFrom = new RecvFrom(e.getMessage());
			} catch (Exception e) {
				e.printStackTrace();
				recvFrom = new RecvFrom(e.getMessage());
			}

			return recvFrom;
		}
		

		// onPostExecute displays the results of the AsyncTask.
		//@Override
		protected void onPostExecute(RecvFrom recvFrom) {
			if (recvFrom != null && !isCancelled()) {
				if (recvFrom.valid) { 
					textOut.append("received " + recvFrom.json + "\n");
					textOut.append("from : " + recvFrom.remoteHost + "\n");
				}
				
				if (recvFrom.errMessage != null)
					textOut.append(recvFrom.errMessage + "\n");
			}
			
			// we're done, clear out the ref to us in the parent class
			downloader = null;
		}
	}	    

}
