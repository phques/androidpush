package com.exercise.AndroidClient;

import java.io.IOException;
import java.net.BindException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetSocketAddress;
import java.net.SocketException;
import org.json.JSONException;

import android.app.Activity;
import android.media.MediaScannerConnection;
import android.net.Uri;
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
	
	
	/** Called when the activity is first created. */
	@Override
	public void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.main);

		textOut = (EditText)findViewById(R.id.textout);
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

	}

	@Override
	protected void onPause() {
		super.onPause();
	}
	
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
				recvFrom = new RecvFrom("JSONException : " + e1.getMessage());
				cancel(false);
			}

			// do the download
			if (!isCancelled() && recvFrom.valid) { 
				// launch the background download 
				recvFrom.execute();
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
				if (recvFrom.getErrMessage() != null) {
					textOut.append(recvFrom.getErrMessage() + "\n");
				}
				else {
					textOut.append("done, file " + recvFrom.filename + "\n");
				}
			}
			
			// we're done
			cleanup();
		}
		
		@Override
		protected void onCancelled(Void q) {
			// we're done
			if (recvFrom != null && recvFrom.getErrMessage() != null) {
				textOut.append(recvFrom.getErrMessage() + "\n");
			}
			
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
			
			recvFrom = new RecvFrom(json, remoteHost);
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
