import QtQuick 2.3
import QtQuick.Controls 1.2
import QtQuick.Window 2.1
import QtQuick.Layouts 1.0

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
	
    ColumnLayout {
        id: columnLayout1
        anchors.fill: parent

        GroupBox {
            id: groupBox1
            anchors.rightMargin: 5
            anchors.leftMargin: 5
            anchors.bottomMargin: 5
            anchors.topMargin: 5
            anchors.fill: parent
            title: qsTr("Found Providers")
//            Layout.fillWidth: true

            ColumnLayout {
                id: columnLayout2
                anchors.fill: parent
                spacing: 15
                //                anchors.right: parent.right
                //                anchors.left: parent.left
                Layout.fillWidth: true
//                anchors.fill: parent

                TableView {
                    id: providersView
                    anchors.bottomMargin: 5
                    anchors.right: parent.right
                    anchors.rightMargin: 5
                    anchors.left: parent.left
                    anchors.leftMargin: 5
                    anchors.top: parent.top
                    anchors.topMargin: 5
                    model: providersMdl

                    //anchors.bottom: parent.bottom
                    anchors.bottom: queryButton.top

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

}
