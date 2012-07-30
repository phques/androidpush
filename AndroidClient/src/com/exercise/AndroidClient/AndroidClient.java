package com.exercise.AndroidClient;

import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.SocketException;

import org.json.JSONException;

import android.app.Activity;
import android.app.DownloadManager;
import android.content.Context;
import android.os.AsyncTask;
import android.os.Bundle;
import android.view.Menu;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

public class AndroidClient extends Activity {

	EditText textOut;
	TextView textIn;
	DownloadWebpageText downloader;
	DownloadManager downloadMgr;
	DatagramSocket socket;
	int cptThread = 0;

	
	/** Called when the activity is first created. */
	@Override
	public void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.main);

		textOut = (EditText)findViewById(R.id.textout);
		textIn = (TextView)findViewById(R.id.textin);
		
		Button buttonSend = (Button)findViewById(R.id.send);
		buttonSend.setOnClickListener(buttonSendOnClickListener);
		
		downloadMgr = (DownloadManager) getSystemService(Context.DOWNLOAD_SERVICE);
	}


	Button.OnClickListener buttonSendOnClickListener =
			new Button.OnClickListener() {
		public void onClick(View arg0) {
			if (socket != null) {
				if (downloader == null) {
					textOut.setText("waiting ...");
					downloader = new DownloadWebpageText(); 
					downloader.execute();
				}
				else {
					cptThread++;
					textOut.setText("already waiting ... " + cptThread);
				}
			}
			else {
				textOut.setText("no socket !!");
			}
		}};
	
	@Override
	public void onStart() {
		super.onStart();  // Always call the superclass method first

		try {
			// open a udp socket to receive commands
			socket = new DatagramSocket(null);	// unbound..
			socket.setReuseAddress(true);		// .. and reuse ! (so we can restart while debugging)
			socket.bind(new InetSocketAddress(4445)); // void return ! cant check if worked !!!
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
				
				recvFrom = new RecvFrom(downloadMgr, json, remoteHost);
				
				recvFrom.execute();
							
			} catch (IOException e) {
				e.printStackTrace();
				recvFrom = new RecvFrom(e.getMessage());
			} catch (JSONException e) {
				e.printStackTrace();
				recvFrom = new RecvFrom(e.getMessage());
			}

			return recvFrom;
		}
		

		// onPostExecute displays the results of the AsyncTask.
		//@Override
		protected void onPostExecute(RecvFrom recvFrom) {
			if (recvFrom != null && !isCancelled()) {
				String text = recvFrom.errMessage;
				
				if (text == null)
					text = "received " + recvFrom.json + "\nfrom : " + recvFrom.remoteHost;
				
				textOut.setText(text);
			}
			
			// we're done, clear out the ref to us in the parent class
			downloader = null;
		}
	}	    

}
