<map version="0.9.0">
<!-- To view this file, download free mind mapping software FreeMind from http://freemind.sourceforge.net -->
<node CREATED="1377099294490" ID="ID_394723472" MODIFIED="1377101709696" STYLE="fork" TEXT="marco polo">
<font NAME="SansSerif" SIZE="15"/>
<node CREATED="1377100373355" FOLDED="true" ID="ID_959423427" MODIFIED="1377101710618" POSITION="right" TEXT="marcoPolo / app connect">
<node CREATED="1377099309997" HGAP="23" ID="ID_1115893101" MODIFIED="1377100796904" STYLE="fork" TEXT="waits on UDP port for connections / requests" VSHIFT="-24">
<font NAME="SansSerif" SIZE="15"/>
<node CREATED="1377099957835" HGAP="26" ID="ID_781887355" MODIFIED="1377100410620" TEXT="could use range/list of ports&#xa;in case of failure to open port" VSHIFT="10">
<icon BUILTIN="help"/>
</node>
</node>
<node CREATED="1377100089551" ID="ID_1921755573" MODIFIED="1377100950213" TEXT="app &apos;connects&apos; to marcoPolo (UDP broadcast)&#xa;clientApp registers to:&#xa;&apos;kwez.org/androidPush/notif&apos;&#xa;  get answer from marcoPolo">
<node CREATED="1377101357557" ID="ID_175850506" MODIFIED="1377101387636" TEXT="marcoPolo sends back &quot;your appId&quot;"/>
<node CREATED="1377100161221" ID="ID_1622029243" MODIFIED="1377101420002" TEXT="app could then connect to &#xa;a fixed TCP socket" VSHIFT="16">
<icon BUILTIN="help"/>
<node CREATED="1377100482850" HGAP="19" ID="ID_1208356818" MODIFIED="1377100839358" TEXT="possible PingAlive" VSHIFT="25">
<icon BUILTIN="help"/>
<node CREATED="1377100725872" HGAP="21" ID="ID_535245107" MODIFIED="1377101429665" TEXT="unregister app when no answer" VSHIFT="5"/>
</node>
</node>
</node>
<node CREATED="1377101264226" ID="ID_852292473" MODIFIED="1377101447969" TEXT="When app closes:&#xa;&quot;unregister me&quot; (appId)">
<node CREATED="1377101324067" HGAP="25" ID="ID_924797808" MODIFIED="1377101394464" TEXT="could be done by just closing permanant TCP socket" VSHIFT="12">
<icon BUILTIN="help"/>
</node>
</node>
</node>
<node CREATED="1377100598364" FOLDED="true" HGAP="41" ID="ID_1798721036" MODIFIED="1377101672205" POSITION="right" TEXT="broadcast msgs" VSHIFT="25">
<node CREATED="1377100253654" ID="ID_1424147096" MODIFIED="1377101084745" TEXT="serverApp sends :&#xa;kwez.org/androidPush/notif&apos;" VSHIFT="5">
<icon BUILTIN="forward"/>
<node CREATED="1377101011489" HGAP="34" ID="ID_363667393" MODIFIED="1377101140964" TEXT="marcoPolo sends back answer to serverApp" VSHIFT="-6">
<icon BUILTIN="back"/>
</node>
<node CREATED="1377100670984" HGAP="36" ID="ID_316452477" MODIFIED="1377101109839" TEXT="marcoPolo sends to clientApp" VSHIFT="2">
<arrowlink DESTINATION="ID_1977219182" ENDARROW="Default" ENDINCLINATION="214;0;" ID="Arrow_ID_1083713839" STARTARROW="None" STARTINCLINATION="214;0;"/>
<icon BUILTIN="forward"/>
</node>
</node>
<node CREATED="1377099524331" ID="ID_1977219182" MODIFIED="1377101063660" STYLE="fork" TEXT="clientApp registered to:&#xa;&apos;kwez.org/androidPush/notif&apos;">
<font NAME="SansSerif" SIZE="15"/>
<node CREATED="1377100984398" HGAP="49" ID="ID_1901858870" MODIFIED="1377101153884" TEXT="possible answer" VSHIFT="16">
<arrowlink DESTINATION="ID_363667393" ENDARROW="Default" ENDINCLINATION="163;0;" ID="Arrow_ID_84087455" STARTARROW="None" STARTINCLINATION="163;0;"/>
<icon BUILTIN="forward"/>
<icon BUILTIN="help"/>
</node>
</node>
</node>
<node CREATED="1377099674415" FOLDED="true" ID="ID_558032627" MODIFIED="1377102071480" POSITION="right" STYLE="fork" TEXT="&apos;services&apos; discovery" VSHIFT="31">
<font NAME="SansSerif" SIZE="14"/>
<node CREATED="1377101566981" ID="ID_58866484" MODIFIED="1377101952782" TEXT="appY registered as service&#xa;&apos;kwez.org/androidPush/getPush&apos;&#xa;(on connect)"/>
<node CREATED="1377099684469" HGAP="24" ID="ID_448274881" MODIFIED="1377101666564" STYLE="fork" TEXT="appX queries for service&#xa;&apos;kwez.org/androidPush/getPush&apos;" VSHIFT="17">
<font NAME="SansSerif" SIZE="14"/>
</node>
<node CREATED="1377101959898" ID="ID_1955642783" MODIFIED="1377102003030" TEXT="marcoPolo answers appX with&#xa;&apos;appY&apos; @ tcp:host:port&#xa;or @ udp:host:port"/>
<node CREATED="1377102046425" ID="ID_1692447688" MODIFIED="1377102064457" TEXT="appX connects to appY ..."/>
</node>
<node CREATED="1377101715687" FOLDED="true" ID="ID_699904154" MODIFIED="1377101932822" POSITION="left" TEXT="messages">
<node CREATED="1377101722881" HGAP="45" ID="ID_1984794553" MODIFIED="1377101929463" TEXT="method/msg&apos;address&apos;" VSHIFT="-27">
<node CREATED="1377101733263" HGAP="22" ID="ID_1176621873" MODIFIED="1377101892602" TEXT="Domain&#xa;eg kwez.org" VSHIFT="5">
<icon BUILTIN="full-1"/>
</node>
<node CREATED="1377101754511" ID="ID_124424908" MODIFIED="1377101901955" TEXT="AppName&#xa;eg androidPush">
<icon BUILTIN="full-2"/>
</node>
<node CREATED="1377101767994" ID="ID_763250280" MODIFIED="1377101913513" TEXT="Method/msg etc&#xa;eg notifPushAvail">
<icon BUILTIN="full-3"/>
</node>
<node CREATED="1377101797478" ID="ID_1615771619" MODIFIED="1377101879985" TEXT="&apos;kwez.org.androidPush/notifPushAvail">
<icon BUILTIN="info"/>
</node>
</node>
</node>
</node>
</map>
