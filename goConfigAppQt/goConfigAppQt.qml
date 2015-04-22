import QtQuick 2.3
import QtQuick.Controls 1.2
import QtQuick.Layouts 1.0
import QtQuick.Window 2.1
import QtQuick.Dialogs 1.1

ApplicationWindow {
    id: applicationWindow1
    title: qsTr("Android Push")
    width: 375
    height: 250

    MessageDialog {
        id: messageDialog
		objectName: "messageDialog"
        title: "May I have your attention please"
        text: "It's so cool that you are using Qt Quick."
        //Component.onCompleted: visible = true
    }
    
    ListModel {
        id: providersMdl
        objectName: "providersMdl"

        //## work around for goqml bug
        function myAppend(json) {
            providersMdl.append(JSON.parse(json))
        }
    }
	
	// ColumnLayout
    GroupBox {
        id: groupBox1
        anchors.margins: 5
        anchors.fill: parent
        title: qsTr("Found Providers ")

        ColumnLayout {
            id: columnLayout2
            //anchors.margins: 5
            anchors.fill: parent
            //spacing: 5

            TableView {
                id: providersView
                model: providersMdl

                Layout.fillHeight: true
                Layout.fillWidth: true

                frameVisible: true
                headerVisible: true
                sortIndicatorVisible: true
                alternatingRowColors: true

                TableViewColumn {
                    role: "name"
                    title: "Name"
                    width: 150
                }
                TableViewColumn {
                    role: "address"
                    title: "Address"
                    width: 150
                }
            }

            Button {
                id: queryButton
                objectName: "queryButton"
                text: "Query"
            }
        }
    }

}
