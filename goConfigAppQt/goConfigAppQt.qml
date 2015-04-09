import QtQuick 2.3
import QtQuick.Controls 1.2
import QtQuick.Window 2.1
import QtQuick.Layouts 1.0

ApplicationWindow {
    id: applicationWindow1
    title: qsTr("Android Push")
    width: 640
    height: 480

    ColumnLayout {
        id: columnLayout1
        anchors.fill: parent

        GroupBox {
            id: groupBox1
            title: qsTr("Providers")
            Layout.fillWidth: true

            RowLayout {
                id: rowLayout
                anchors.fill: parent

                Button {
                    text: "Button"
                }
            }
        }
    }

}
