package com.exercise.AndroidClient;

//import java.io.BufferedReader;
import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.SocketException;
import java.net.UnknownHostException;

import org.json.JSONException;
import org.json.JSONObject;
import org.json.JSONTokener;

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
	DatagramSocket socket = null;
	int cptThread = 0;

	class RecvFrom {

		int port;
		String destinationType;  
		String subDir;  
		String file;  
		
		RecvFrom(String json) throws JSONException {
			parseJson(json);
		}
		
		private void parseJson(String json) throws JSONException {			
			JSONObject object = (JSONObject) new JSONTokener(json).nextValue();
			port = object.getInt("recvFromPort");
			file = object.getString("file");  
			subDir = object.getString("subDir");  
			destinationType = object.getString("destinationType");  
		}
		
	}
	
	/** Called when the activity is first created. */
	@Override
	public void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.main);

		textOut = (EditText)findViewById(R.id.textout);
		textIn = (TextView)findViewById(R.id.textin);
		
		Button buttonSend = (Button)findViewById(R.id.send);
		buttonSend.setOnClickListener(buttonSendOnClickListener);
	}


	Button.OnClickListener buttonSendOnClickListener =
			new Button.OnClickListener() {
		public void onClick(View arg0) {
			if (socket != null) {
				if (downloader == null) {
					textOut.setText("waiting ...");
					//						new DownloadWebpageText().execute();
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
	
	
	private class DownloadWebpageText extends AsyncTask<Void, Void, String> {

		//@Override
		protected String doInBackground(Void ... voids) {

			String result=null;
//			Socket socket = null;
//			BufferedReader in = null;
			
			try {
				byte[] buf = new byte[1024*4];
				DatagramPacket packet = new DatagramPacket(buf, buf.length);
				socket.receive(packet);

				String json = new String(packet.getData(), 0, packet.getLength());
				RecvFrom recvFrom = new RecvFrom(json);
				
				InetAddress addr = packet.getAddress();
				result = "received " + json + "\nfrom : " + addr.toString();
				
				/*
				socket = new Socket();
				socket.connect(new InetSocketAddress("192.168.1.119", 8888), 0);
				
				in = new BufferedReader(new InputStreamReader(socket.getInputStream()));				
				result = in.readLine();
				*/
			/*} catch (UnknownHostException e) {
				e.printStackTrace();
				result = e.getMessage();
			*/} catch (IOException e) {
				e.printStackTrace();
				result = e.getMessage();
			} /*catch (Exception e) {
				e.printStackTrace();
				result = e.getMessage();
			}*/ catch (JSONException e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
				result = e.getMessage();
			}
			/*finally{
				if (socket != null){
					socket.close();

					try {
						socket.close();
					} catch (IOException e) {
						e.printStackTrace();
					}
				}
				
				if (in != null){
					try {
						in.close();
					} catch (IOException e) {
						e.printStackTrace();
					}
				}
			}*/
			return result;
		}
		

		// onPostExecute displays the results of the AsyncTask.
		//@Override
		protected void onPostExecute(String result) {
			if (result != null)
				textOut.setText(result);
			
			downloader = null;
		}
	}	    

}
