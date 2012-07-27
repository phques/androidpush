package com.exercise.AndroidClient;

import java.io.BufferedReader;
import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.UnknownHostException;

import android.app.Activity;
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
	DownloadWebpageText downloader = null;

	/** Called when the activity is first created. */
	@Override
	public void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.main);

		textOut = (EditText)findViewById(R.id.textout);
		Button buttonSend = (Button)findViewById(R.id.send);
		textIn = (TextView)findViewById(R.id.textin);
		buttonSend.setOnClickListener(buttonSendOnClickListener);
	}

	Button.OnClickListener buttonSendOnClickListener
	= new Button.OnClickListener(){

		//		@Override
		public void onClick(View arg0) {
			textOut.setText("waiting ...");
			new DownloadWebpageText().execute();
//			downloader = new DownloadWebpageText(); 
//			downloader.execute();
		}};

	@Override
	public boolean onCreateOptionsMenu(Menu menu) {
		getMenuInflater().inflate(R.menu.main, menu);
		return true;
	}
	
//	public void onStop() {
//		super.onStop();  // Always call the superclass method first
//		if (downloader != null)
//			downloader.cancel(true);
//		downloader = null;
//	}

	private class DownloadWebpageText extends AsyncTask<Void, Void, Long> {
		private String result;
		//@Override
		protected Long doInBackground(Void ... voids) {

//			Socket socket = null;
			BufferedReader in = null;
			DatagramSocket socket = null;
			
			try {
//				socket = new DatagramSocket(4445);
				socket = new DatagramSocket(null);	// unbound..
				socket.setReuseAddress(true);		// .. and reuse ! (so we can restart while debugging)
				socket.bind(new InetSocketAddress(4445)); // void return ! cant check if worked !!!
//				socket.setBroadcast(true);

				byte[] buf = new byte[1024*4];
				DatagramPacket packet = new DatagramPacket(buf, buf.length);
				socket.receive(packet);

				result = new String(packet.getData(), 0, packet.getLength());
				InetAddress addr = packet.getAddress();
				result += addr.toString();
								
				/*
				socket = new Socket();
				socket.connect(new InetSocketAddress("192.168.1.119", 8888), 0);
				
				in = new BufferedReader(new InputStreamReader(socket.getInputStream()));				
				result = in.readLine();
				*/
			} catch (UnknownHostException e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
				result = e.getMessage();
			} catch (IOException e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
				result = e.getMessage();
			} catch (Exception e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
				result = e.getMessage();
			}
			finally{
				if (socket != null){
					socket.close();
//					try {
//						socket.close();
//					} catch (IOException e) {
//						// TODO Auto-generated catch block
//						e.printStackTrace();
//					}
				}
				
				if (in != null){
					try {
						in.close();
					} catch (IOException e) {
						// TODO Auto-generated catch block
						e.printStackTrace();
					}
				}
			}
			return 0L;
		}
		
		// onPostExecute displays the results of the AsyncTask.
		//@Override
		protected void onPostExecute(Long tt) {
			textOut.setText(result);
//			downloader = null;
		}
	}	    

}
