import QtQuick 2.3
import QtQuick.Controls 1.2
import QtQuick.Layouts 1.0
import QtQuick.Window 2.1

ApplicationWindow {
    id: applicationWindow1
    title: qsTr("Android Push")
    width: 375
    height: 250

    ListModel {
        id: providersMdl
        objectName: "providersMdl"

        //## work around for goqml bug
        function myAppend(json) {
            providersMdl.append(JSON.parse(json))
        }
    }

	// Column + anchors
    GroupBox {
        id: groupBox1
        anchors.margins: 5
        anchors.fill: parent
        title: qsTr("Found Providers")

        Column {
            id: columnLayout2
            anchors.fill: parent

            TableView {
                id: providersView
                model: providersMdl

                anchors.right: parent.right
                anchors.left: parent.left
                anchors.top: parent.top
				
                // anchor bottom of table to queryButton
                // (which is anchored to bottom of our parent)
                anchors.bottom: queryButton.top
                anchors.bottomMargin: 5

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
                anchors.bottom: parent.bottom
            }
        }
    }
}
