import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQuick.Controls.Material 2.12
import QtQuick.Layouts 1.12

Page {
    id: root

    readonly property real internalLabelsWidth: 50

    header: Label {
        text: qsTr("Wallets")
        font.pixelSize: Qt.application.font.pixelSize * 2
        padding: 10
    }

    footer: RowLayout {
        ItemDelegate {
            id: buttonAddWallet
            text: qsTr("Add wallet")
            Layout.fillWidth: true
            display: ItemDelegate.TextUnderIcon
        }
        ItemDelegate {
            id: buttonLoadWallet
            text: qsTr("Load wallet")
            Layout.fillWidth: true
            display: ItemDelegate.TextUnderIcon
        }
    }

    ScrollView {
        id: scrollItem
        anchors.fill: parent
        ListView {
            id: walletList
            anchors.fill: parent
            model: listModel
            delegate: delegateWalletList
        }
    }

    // Roles: name, encryptionEnabled, sky, coinHours
    // Use listModel.append( { "name": value, "encryptionEnabled": value, "sky": value, "coinHours": value } )
    // Or implement the model in the backend (a more recommendable approach)
    ListModel {
        id: listModel
    }

    Component {
        id: delegateWalletList

        ItemDelegate {
            width: walletList.width

            RowLayout {
                anchors.fill: parent
                anchors.leftMargin: 30
                anchors.rightMargin: 50
                spacing: 20

                // TODO: add an 'encryption-disabled' SVG icon
                Image {
                    id: lockIcon
                    source: "qrc:/images/security.svg"
                    sourceSize: "24x24"
                }

                Label {
                    id: labelWalletName
                    text: name
                    Layout.fillWidth: true
                }

                Label {
                    id: labelskycoins
                    text: sky
                    color: Material.accent
                    horizontalAlignment: Text.AlignRight
                    Layout.preferredWidth: internalLabelsWidth
                }

                Label {
                    id: labelCoins
                    text: coinHours
                    horizontalAlignment: Text.AlignRight
                    Layout.preferredWidth: internalLabelsWidth
                }
            }
        }
    }
}