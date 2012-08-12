package com.exercise.AndroidClient;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.io.PrintStream;
import java.net.Socket;
import java.net.UnknownHostException;
import java.util.HashMap;
import java.util.Map;

import org.json.JSONException;
import org.json.JSONObject;
import org.json.JSONTokener;

import android.media.MediaScannerConnection;
import android.os.Environment;
import android.util.Log;
import android.net.Uri;

public 	class RecvFrom {
	Map<String,String> destDirTypes;
	
	String json;
	Boolean valid = false;
	private String errMessage;
	
	int remotePort;
	int pushId = 0;
	long fileLength = 0;
	String remoteHost;
	String destinationType;  
	String subDir;  
	String filename;  
	
	final int BUFFER_SIZE = 1024*32;
	Socket serverSocket=null;

	RecvFrom(String json, String remoteHost) throws JSONException {
		this.remoteHost = remoteHost;
		this.json = json;
		parseJson(json);
		// setup map for getExternalStoragePublicDirectory values 
		destDirTypes = new HashMap<String,String>();
		destDirTypes.put("music", Environment.DIRECTORY_MUSIC);
		destDirTypes.put("downloads", Environment.DIRECTORY_DOWNLOADS);
		destDirTypes.put("pictures", Environment.DIRECTORY_PICTURES);
		destDirTypes.put("movies", Environment.DIRECTORY_MOVIES);
		valid = true;
	}
	
	RecvFrom(String errMessage) {
		this.errMessage = errMessage;
		valid = false;
	}

	public String getErrMessage() {
		return errMessage;
	}

	public void setErrMessage(String errMessage) {
		this.errMessage = errMessage;
	}

	public void addErrMessage(String errMessage) {
		if (this.errMessage == null)
			this.errMessage = errMessage;
		else
			this.errMessage += errMessage;
	}
	
	public void execute() {
		
		try {
			// do a 1st connection to server, sending an ACK string
			// after the server receives this it knows to wait for the download request
			serverSocket = new Socket(remoteHost, remotePort);
			
			PrintStream outStream = new PrintStream(serverSocket.getOutputStream());
			outStream.print("ACK\n");
			outStream.print("pushId" + pushId + "\n");
			
			//
			receiveFile();
							
		} catch (UnknownHostException e) {
			e.printStackTrace();
			addErrMessage("\n" + e.getMessage());
		} catch (IOException e) {
			e.printStackTrace();
			addErrMessage("\n" + e.getMessage());
		}
		finally {
			if (serverSocket != null) {
				try {
					serverSocket.close();
				} catch (IOException e) {
					e.printStackTrace();
				}
				serverSocket = null;
			}
		}
		
	}

	private void receiveFile() throws IOException {

		// get file destination
		File destFilePath = getDestFilePath();

		// write to file from socket
		OutputStream out = null;
		BufferedInputStream in = null;
		//InputStream in = null;
		byte[] buffer = new byte [BUFFER_SIZE];
		long totalRead = 0;
		
		try {
			// open file
			out = new BufferedOutputStream(
					new FileOutputStream(destFilePath));
			
			// open socket stream
			in = new BufferedInputStream(serverSocket.getInputStream());
			//in = serverSocket.getInputStream();
			
			// loop read/write
			while (totalRead < fileLength){
				int nbRead = in.read(buffer, 0, buffer.length);
				if (nbRead <= 0) // ooops
					break;
				out.write(buffer, 0, nbRead);
				totalRead += nbRead;
			}
			
		}
		catch (Exception e) {
			e.printStackTrace();
			addErrMessage("\n" + e.getMessage());
		}
		finally {
			if (out != null) {
				out.close();
			}
			if (in != null){
				in.close();
			}

			// check that we got the whole file
			if (totalRead == fileLength) {
		        // Tell the media scanner about the new file so that it is
		        // immediately available to the user.
//nb: cant call mediascanner here, needs Context
/*		        MediaScannerConnection.scanFile(this,
		                new String[] { recvFrom.getFilePath() }, null,
		                new MediaScannerConnection.OnScanCompletedListener() {
		            public void onScanCompleted(String path, Uri uri) {
		                Log.i("ExternalStorage", "Scanned " + path + ":");
		                Log.i("ExternalStorage", "-> uri=" + uri);
		            }
		        });
*/			}
			else {
				// delete the file, it is invalid
				destFilePath.delete();
				addErrMessage("download stopped before full file length, deleting file");
			}
		}
		
	}

	private File getDestFilePath() throws IOException {
		
		// map "music" to Environment.DIRECTORY_MUSIC 
		String destDir = destDirTypes.get(destinationType);
		if (destDir == null)
			throw new IllegalArgumentException("Invalid destination dir type : " + destinationType);
		
		// get the public 'standard directory', ie Download, Music etc
		File dir = Environment.getExternalStoragePublicDirectory(destDir);
		
		// add any subdir asked by caller & create dir struct upto the subdir if required
		dir = new File(dir, subDir);
		if (!dir.exists() || !dir.isDirectory()) {
			if (!dir.mkdirs())
				throw new IOException("Failed to create directories : " + dir.toString());
		}

		// the full path to the file
		return new File(dir, filename);
	}

	/** extract data from the json object string */
	private void parseJson(String json) throws JSONException {			
		JSONObject object = (JSONObject) new JSONTokener(json).nextValue();
		remotePort = object.getInt("recvFromPort");
		filename = object.getString("file");  
		subDir = object.getString("subDir");  
		destinationType = object.getString("destDirType");
		pushId = object.getInt("pushId");
		fileLength = object.getLong("fileLength");
	}
	
}
