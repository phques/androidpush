
import java.io.*;
import java.net.*;


public class testNetJ {
	public static void main2(String args[]) {
		try {
			DatagramSocket socket = new DatagramSocket(4445);

			byte[] buf = new byte[1024*4];
			DatagramPacket packet = new DatagramPacket(buf, buf.length);
			socket.receive(packet);

			String received = new String(packet.getData(), 0, packet.getLength());
			System.out.println("recv: " + received);

			socket.close();
		}
		catch(Exception e) {
			System.out.println("Whoops! It didn't work!\n");
			System.out.println(e.getMessage());
		}
	}

	public static void main(String args[]) {
		try {
			Socket skt = new Socket();
			skt.connect(new InetSocketAddress("localhost", 8888), 10000);
			BufferedReader in = new BufferedReader(new
					InputStreamReader(skt.getInputStream()));
			System.out.print("Received string: '");

			while (!in.ready()) {}
			System.out.println(in.readLine()); // Read one line and output it

			System.out.print("'\n");
			in.close();

			skt.close();
		}
		catch(Exception e) {
			System.out.println("Whoops! It didn't work!\n" + e.getMessage());
		}
	}	  

}

