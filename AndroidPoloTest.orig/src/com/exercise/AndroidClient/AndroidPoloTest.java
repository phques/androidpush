package com.exercise.AndroidClient;

// AndroidPush project
// test marco/polo: client polo
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

import java.io.IOException;
import java.net.BindException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.net.SocketException;

import android.app.Activity;
import android.os.AsyncTask;
import android.os.Bundle;
import android.text.InputFilter.LengthFilter;
import android.util.Log;
import android.view.Menu;
import android.view.View;
import android.widget.EditText;


public class AndroidPoloTest extends Activity {

	EditText textOut;
	short udpPort = 4444;
	RecvMarcoAsyncTask recvAsyncTask;
	
	// results
	String recvdMsg = "";
	String remoteHost = "";
	SocketAddress remoteHostSockAddr;
	
	
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
			recvAsyncTask = new RecvMarcoAsyncTask();
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
