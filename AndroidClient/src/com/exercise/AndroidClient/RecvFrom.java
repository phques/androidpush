package com.exercise.AndroidClient;

import java.io.File;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

import org.json.JSONException;
import org.json.JSONObject;
import org.json.JSONTokener;

import android.app.DownloadManager;
import android.app.DownloadManager.Request;
import android.net.Uri;
import android.os.Environment;


public 	class RecvFrom {
	DownloadManager downloadMgr;
	Map<String,String> destDirTypes;
	
	String json;
	Boolean valid = false;
	private String errMessage;
	
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
	
	public void execute() throws IOException {
		// get source URL
		Uri srcUri = getSrcFileUri();
		
		// get file destination
		Uri destFileUri = getDestFileUri();
		
		// build the request we will pass to download manager 
        Request request = new Request(srcUri);
        request.setDestinationUri(destFileUri)
        	.setDescription("Push to Android from PC")
        	.setTitle(filename)
        	//.setAllowedOverMetered(false) //api 16
        	.setAllowedNetworkTypes(Request.NETWORK_WIFI) // only through wifi !
        	.setNotificationVisibility(Request.VISIBILITY_VISIBLE)
        	.allowScanningByMediaScanner();
    
    	// ask download manager to download our file.
    	// system service that runs in the background,
    	// will show status in notif bar, can be seen / stop etc in the 'Downloads' app
    	long downloadId = downloadMgr.enqueue(request);
	}

	private Uri getSrcFileUri() {
		// get source URL 
		Uri.Builder builder = new Uri.Builder();
		builder.scheme("http")
			.encodedPath("//" + remoteHost + ":" + remotePort)
			.appendPath(filename);
		
		return builder.build();
		
//		String srcUrl = "http:/" + remoteHost + ":" + remotePort;
//		srcUrl = srcUrl + "/" + filename;
//		Uri srcUri = Uri.parse(srcUrl);
	}
	

	private Uri getDestFileUri() throws IOException {
		
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
		File path = new File(dir, filename);
		return Uri.fromFile(path);
	}

	/** extract data from the jons object string */
	private void parseJson(String json) throws JSONException {			
		JSONObject object = (JSONObject) new JSONTokener(json).nextValue();
		remotePort = object.getInt("recvFromPort");
		filename = object.getString("file");  
		subDir = object.getString("subDir");  
		destinationType = object.getString("destDirType");  
	}
	
}
