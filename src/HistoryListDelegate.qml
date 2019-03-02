import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQuick.Controls.Material 2.12
import QtQuick.Layouts 1.12

ItemDelegate {
    id: root

    enum Type {
        Send,
        Receive
    }

    implicitWidth: parent.width
    implicitHeight: (columnLayoutMainContent.height < 78 ? 78 : columnLayoutMainContent.height) + rowLayoutRoot.anchors.topMargin + rowLayoutRoot.anchors.bottomMargin

    RowLayout {
        id: rowLayoutRoot
        anchors.fill: parent
        anchors.leftMargin: 20
        anchors.rightMargin: 20
        anchors.topMargin: 10
        anchors.bottomMargin: 12

        spacing: 20

        Image {
            source: "qrc:/images/send-blue.svg"
            sourceSize: "32x32"
            fillMode: Image.PreserveAspectFit
            mirror: type === HistoryListDelegate.Type.Receive
            Layout.alignment: Qt.AlignTop | Qt.AlignLeft
        }

        ColumnLayout {
            id: columnLayoutMainContent
            Layout.fillWidth: true
            Layout.alignment: Qt.AlignTop

            RowLayout {
                spacing: 20
                Layout.fillWidth: true

                Label {
                    font.bold: true
                    text: (type === HistoryListDelegate.Type.Receive ? qsTr("Received") : qsTr("Sent")) + " SKY"
                }

                Label {
                    Material.foreground: Material.Grey
                    text: date // model's role
                    font.pointSize: Qt.application.font.pointSize * 0.9
                }
            }

            ColumnLayout {
                RowLayout {
                    id: rowLayoutSent
                    visible: type === HistoryListDelegate.Type.Send
                    Image {
                        source: "qrc:/images/qr.svg"
                        sourceSize: "24x24"
                    }
                    Label {
                        text: sentAddress // model's role
                        font.family: "Code New Roman"
                        Layout.fillWidth: true
                    }
                }
                RowLayout {
                    id: rowLayoutReceive
                    Image {
                        source: "qrc:/images/qr.svg"
                        sourceSize: "24x24"
                    }
                    Label {
                        text: receivedAddress // model's role
                        font.family: "Code New Roman"
                        Layout.fillWidth: true
                    }
                }
            } // ColumnLayout (addresses)
        } // ColumnLayout (main content)

        Label {
            text: (type === HistoryListDelegate.Type.Receive ? "" : "-") + amount + " SKY" // model's role
            font.pointSize: Qt.application.font.pointSize * 1.25
            font.bold: true
            Layout.alignment: Qt.AlignTop | Qt.AlignRight
        }

    } // RowLayout (root)
}