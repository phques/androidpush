package com.exercise.AndroidClient;

import java.io.File;
import java.io.IOException;
import java.net.URI;
import java.util.HashMap;
import java.util.Map;

import org.json.JSONException;
import org.json.JSONObject;
import org.json.JSONTokener;

import android.net.Uri;
import android.os.Environment;

import android.app.DownloadManager;
import android.app.DownloadManager.Request;


public 	class RecvFrom {
	DownloadManager downloadMgr;
	Map<String,String> destDirTypes;
	
	String json;
	String errMessage;
	
	int remotePort;
	String remoteHost;
	String destinationType;  
	String subDir;  
	String filename;  
	

	RecvFrom(DownloadManager downloadMgr, String json, String remoteHost) throws JSONException {
		this.downloadMgr = downloadMgr;
		this.remoteHost = remoteHost;
		this.json = json;
		parseJson(json);
		// setup map for getExternalStoragePublicDirectory values 
		destDirTypes = new HashMap<String,String>();
		destDirTypes.put("music", Environment.DIRECTORY_MUSIC);
		destDirTypes.put("downloads", Environment.DIRECTORY_DOWNLOADS);
		destDirTypes.put("pictures", Environment.DIRECTORY_PICTURES);
		destDirTypes.put("movies", Environment.DIRECTORY_MOVIES);
	}
	
	RecvFrom(String errMessage) {
		this.errMessage = errMessage;
	}
	
	void execute() throws IOException {
		// get file destination
		File destFilePath = getDestFilepath();
		
		// need to convert java.net.URI to android.net.Uti ;-p 
		URI destFileURI = destFilePath.toURI();
		Uri destFileUri = Uri.parse(destFileURI.toString());
//		String uriStr = destFilePath.toString();
		
/*        Request request = new Request();
        
        request.setDestinationUri(destFileUri);
        request.setDescription("Push to Android from PC");
*/        
        try {
        	// ask download manager to download our file
        	// system service that runs in the background,
        	// will show status in notif bar, can be seen / stop etc in 'Downloads app'
//        	downloadMgr.enqueue(request);
        }
        catch (Exception ex) {        
        	ex.printStackTrace();
        }
		
	}
	

	private File getDestFilepath() throws IOException {
		
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

	/** */
	private void parseJson(String json) throws JSONException {			
		JSONObject object = (JSONObject) new JSONTokener(json).nextValue();
		remotePort = object.getInt("recvFromPort");
		filename = object.getString("file");  
		subDir = object.getString("subDir");  
		destinationType = object.getString("destinationType");  
	}
	
}
