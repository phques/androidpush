<map version="0.9.0">
<!-- To view this file, download free mind mapping software FreeMind from http://freemind.sourceforge.net -->
<node CREATED="1377099294490" ID="ID_394723472" MODIFIED="1377132702936" STYLE="fork" TEXT="marco polo&#xa;simple">
<font NAME="SansSerif" SIZE="15"/>
<node CREATED="1377100373355" FOLDED="true" HGAP="16" ID="ID_959423427" MODIFIED="1377133913210" POSITION="right" TEXT="marcoPolo / app connect" VSHIFT="-24">
<node CREATED="1377099309997" FOLDED="true" HGAP="23" ID="ID_1115893101" MODIFIED="1377133881874" STYLE="fork" TEXT="waits on UDP port for connections / requests" VSHIFT="-24">
<font NAME="SansSerif" SIZE="15"/>
<node CREATED="1377099957835" HGAP="26" ID="ID_781887355" MODIFIED="1377100410620" TEXT="could use range/list of ports&#xa;in case of failure to open port" VSHIFT="10">
<icon BUILTIN="help"/>
</node>
</node>
<node CREATED="1377100089551" FOLDED="true" ID="ID_1921755573" MODIFIED="1377133912246" TEXT="app &apos;connects&apos; to marcoPolo (UDP broadcast)">
<node CREATED="1377101357557" ID="ID_175850506" MODIFIED="1377101387636" TEXT="marcoPolo sends back &quot;your appId&quot;"/>
<node CREATED="1377100161221" FOLDED="true" ID="ID_1622029243" MODIFIED="1377133911554" TEXT="app then connects to &#xa;a fixed TCP socket" VSHIFT="16">
<node CREATED="1377100482850" HGAP="37" ID="ID_1208356818" MODIFIED="1377133906000" TEXT="marcoPol sends PingAlive" VSHIFT="7"/>
<node CREATED="1377100725872" HGAP="40" ID="ID_535245107" MODIFIED="1377132401749" TEXT="unregister app when no answer" VSHIFT="11"/>
</node>
</node>
<node CREATED="1377101264226" FOLDED="true" ID="ID_852292473" MODIFIED="1377133802816" TEXT="When app closes:&#xa;&quot;unregister me&quot; (appId)">
<node CREATED="1377101324067" HGAP="25" ID="ID_924797808" MODIFIED="1377131047094" TEXT="could be done by just closing permanent TCP socket" VSHIFT="12">
<icon BUILTIN="help"/>
</node>
</node>
</node>
<node CREATED="1377099674415" FOLDED="true" HGAP="34" ID="ID_558032627" MODIFIED="1377133809750" POSITION="right" STYLE="fork" TEXT="msgs discovery" VSHIFT="4">
<font NAME="SansSerif" SIZE="14"/>
<node CREATED="1377101566981" ID="ID_58866484" MODIFIED="1377132274128" TEXT="appY registers msg&#xa;&apos;kwez.org/androidPush/getPush&apos;"/>
<node CREATED="1377099684469" HGAP="24" ID="ID_448274881" MODIFIED="1377132292065" STYLE="fork" TEXT="appX queries for msg&#xa;&apos;kwez.org/androidPush/getPush&apos;" VSHIFT="17">
<font NAME="SansSerif" SIZE="14"/>
</node>
<node CREATED="1377101959898" FOLDED="true" ID="ID_1955642783" MODIFIED="1377133808766" TEXT="marcoPolo answers appX with&#xa;host:port">
<node CREATED="1377102046425" ID="ID_1692447688" MODIFIED="1377132351451" TEXT="appX talks to appY ..."/>
<node CREATED="1377132538773" ID="ID_877205127" MODIFIED="1377132569163" TEXT="protocol for this msg must&#xa;be decided between apps"/>
</node>
</node>
<node CREATED="1377132575487" FOLDED="true" HGAP="38" ID="ID_778880487" MODIFIED="1377133811630" POSITION="right" TEXT="msg broadcast" VSHIFT="23">
<node CREATED="1377132581168" ID="ID_112861534" MODIFIED="1377133019628" TEXT="serverApp requests broadcast&#xa;kwez.org/androidPush/transferUpdt&#xa;payload=&quot;..&quot;"/>
<node CREATED="1377132591823" ID="ID_146670919" MODIFIED="1377133504453" TEXT="marcoPolo broadcasts msg &#xa;to all registered clientApps"/>
<node CREATED="1377133512054" ID="ID_297170486" MODIFIED="1377133519462" TEXT="(no retry / ack)"/>
</node>
<node CREATED="1377101715687" ID="ID_699904154" MODIFIED="1377130078247" POSITION="left" TEXT="messages">
<node CREATED="1377130119924" ID="ID_835544906" MODIFIED="1377187477271" TEXT="&quot;marco.polo:{JSON object msg}&quot;&#xa;marcoPolo msg/cmd">
<node CREATED="1377130359609" HGAP="84" ID="ID_1497430216" MODIFIED="1377186361220" TEXT="version" VSHIFT="7">
<node CREATED="1377130366619" ID="ID_1482147845" MODIFIED="1377130369433" TEXT="major"/>
<node CREATED="1377130370145" ID="ID_1652153778" MODIFIED="1377130372436" TEXT="minor"/>
</node>
<node CREATED="1377132065856" HGAP="80" ID="ID_170084488" MODIFIED="1377186360032" TEXT="action" VSHIFT="46">
<node CREATED="1377132657041" ID="ID_1212759630" MODIFIED="1377132671799" TEXT="register.app"/>
<node CREATED="1377132661197" ID="ID_847473981" MODIFIED="1377132674698" TEXT="unregister.app"/>
<node CREATED="1377132068450" ID="ID_1325827875" MODIFIED="1377132677085" TEXT="register.msg"/>
<node CREATED="1377132113524" ID="ID_1571369487" MODIFIED="1377132679866" TEXT="unregister.msg"/>
<node CREATED="1377132115928" ID="ID_1697067950" MODIFIED="1377132682716" TEXT="query.msg"/>
<node CREATED="1377132235378" ID="ID_841304758" MODIFIED="1377133113642" TEXT="broadcast.msg"/>
</node>
<node CREATED="1377101722881" FOLDED="true" HGAP="47" ID="ID_1984794553" MODIFIED="1377186359220" TEXT="msg name" VSHIFT="26">
<node CREATED="1377101733263" HGAP="22" ID="ID_1176621873" MODIFIED="1377101892602" TEXT="Domain&#xa;eg kwez.org" VSHIFT="5">
<icon BUILTIN="full-1"/>
</node>
<node CREATED="1377101754511" ID="ID_124424908" MODIFIED="1377101901955" TEXT="AppName&#xa;eg androidPush">
<icon BUILTIN="full-2"/>
</node>
<node CREATED="1377101767994" ID="ID_763250280" MODIFIED="1377101913513" TEXT="Method/msg etc&#xa;eg notifPushAvail">
<icon BUILTIN="full-3"/>
</node>
<node CREATED="1377101797478" HGAP="27" ID="ID_1615771619" MODIFIED="1377132873116" TEXT="&apos;kwez.org/androidPush/notifPushAvail&apos;" VSHIFT="19">
<icon BUILTIN="info"/>
</node>
</node>
<node CREATED="1377132835247" HGAP="55" ID="ID_1428439205" MODIFIED="1377133413741" TEXT="opt msg payload (string)&#xa;for broadcast msgs" VSHIFT="25"/>
</node>
</node>
</node>
</map>
