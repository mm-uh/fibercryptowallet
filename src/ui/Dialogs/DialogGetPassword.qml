import QtQuick 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls 2.12

// Resource imports
// import "qrc:/ui/src/ui/"
import "../" // For quick UI development, switch back to resources when making a release

Dialog {
    id: dialogGetPassword

    property alias headerMessage: labelHeaderMessage.text
    property alias headerMessageColor: labelHeaderMessage.color
    property alias password: passwordRequester.text

    function clear() {
        passwordRequester.clear()
        standardButton(Dialog.Ok).enabled = password
    }

    onAboutToShow: {
        clear()
        passwordRequester.forceTextFocus()
    }

    title: qsTr("Password requested")
    standardButtons: Dialog.Ok | Dialog.Cancel

    Flickable {
        id: flickable
        anchors.fill: parent
        contentHeight: columnLayoutRoot.height
        clip: true

        ColumnLayout {
            id: columnLayoutRoot
            width: parent.width
            spacing: 10

            Label {
                id: labelHeaderMessage

                Layout.fillWidth: true
                wrapMode: Text.WordWrap
                visible: text
            }

            PasswordRequester {
                id: passwordRequester

                Layout.fillWidth: true

                onPasswordForgotten: {
                    Qt.openUrlExternally("http://skycoin.com/")
                }
                onTextChanged: {
                    dialogGetPassword.standardButton(Dialog.Ok).enabled = text !== ""
		            password = text
                }
            }
        } // ColumnLayout (root)
    } // Flickable
}
