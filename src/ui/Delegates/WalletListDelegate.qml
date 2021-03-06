import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQuick.Controls.Material 2.12
import QtQuick.Layouts 1.12
import QtGraphicalEffects 1.12
import WalletsManager 1.0

// Resource imports
// import "qrc:/ui/src/ui/Dialogs"
import "../Dialogs/" // For quick UI development, switch back to resources when making a release

Item {
    id: root

    readonly property real delegateHeight: 30
    property bool emptyAddressVisible: true
    property bool expanded: expand
    // The following property is used to avoid a binding conflict with the `height` property.
    // Also avoids a bug with the animation when collapsing a wallet
    readonly property real finalViewHeight: expanded ? delegateHeight*(addressList.count) + 50 : 0


    width: walletList.width
    height: itemDelegateMainButton.height + (expanded ? finalViewHeight : 0)

    Behavior on height { NumberAnimation { duration: 250; easing.type: Easing.OutQuint } }

    ColumnLayout {
        id: delegateColumnLayout
        anchors.fill: parent

        ItemDelegate {
            id: itemDelegateMainButton
            Layout.fillWidth: true
            Layout.alignment: Qt.AlignTop
            font.bold: expanded

            RowLayout {
                id: delegateRowLayout
                anchors.fill: parent
                anchors.leftMargin: listWalletLeftMargin
                anchors.rightMargin: listWalletRightMargin
                spacing: listWalletSpacing

                Image {
                    id: status
                    source: statusIcon
                    sourceSize: "24x24"
                }

                Label {
                    id: labelWalletName
                    text: name // a role of the model
                    Layout.fillWidth: true
                }

                Item {
                    id: itemImageLockIcon

                    width: lockIcon.width
                    height: lockIcon.height

                    Image {
                        id: lockIcon
                        source: "qrc:/images/resources/images/icons/lock" + (encryptionEnabled ? "On" : "Off") + ".svg"
                        sourceSize: "24x24"
                    }

                    ColorOverlay {
                        anchors.fill: lockIcon
                        source: lockIcon
                        color: Material.theme === Material.Dark ? Material.foreground : "undefined"
                    }
                }

                Label {
                    id: labelSky
                    text: sky // a role of the model
                    color: Material.accent
                    horizontalAlignment: Text.AlignRight
                    Layout.preferredWidth: internalLabelsWidth
                }

                Label {
                    id: labelCoins
                    text: coinHours // a role of the model
                    horizontalAlignment: Text.AlignRight
                    Layout.preferredWidth: internalLabelsWidth
                }
            } // RowLayout

            onClicked: {

                expanded = !expanded
            }
        } // ItemDelegate

        ListView {
            id: addressList
            model: listAddresses
            implicitHeight: expanded ? delegateHeight*(addressList.count) + 50 : 0
            property alias parentRoot: root
            opacity: expanded ? 1.0 : 0.0
            clip: true
            interactive: false
            Layout.fillWidth: true
            Layout.alignment: Qt.AlignTop

            Behavior on implicitHeight { NumberAnimation { duration: 250; easing.type: Easing.OutQuint } }
            Behavior on opacity { NumberAnimation { duration: expanded ? 250 : 1000; easing.type: Easing.OutQuint } }

            delegate: WalletListAddressDelegate {
                width: walletList.width
                height: index == 0 ? delegateHeight + 20 : visible ? delegateHeight : 0

                onAddAddressesRequested: {
                    dialogAddAddresses.open()
                }
                onEditWalletRequested: {
                    dialogEditWallet.originalWalletName = name
                    dialogEditWallet.name = name
                    dialogEditWallet.open()
                }
            }

        } // ListView
    } // ColumnLayout

    DialogAddAddresses {
        id: dialogAddAddresses
        anchors.centerIn: Overlay.overlay

        modal: true
        focus: true

        onAccepted: {
            if (encryptionEnabled){
                dialogGetPasswordForAddAddresses.title = "Enter Password"
                dialogGetPasswordForAddAddresses.height = dialogGetPassword.height - 20
                dialogGetPasswordForAddAddresses.nAddress = spinValue
                dialogGetPasswordForAddAddresses.open()


            } else{
                walletManager.newWalletAddress(fileName, spinValue, "")
                listAddresses.loadModel(walletManager.getAddresses(fileName))
            }


        }
        onRejected: {
            console.log("Adding rejected")
        }
    } // DialogAddAddresses
    DialogGetPassword {
        id: dialogGetPasswordForAddAddresses
        anchors.centerIn: Overlay.overlay
        property int nAddress
        width: applicationWindow.width > 540 ? 540 - 120 : applicationWindow.width - 40
        height: applicationWindow.height > 570 ? 570 - 180 : applicationWindow.height - 40

        focus: true
        modal: true

        onAccepted: {
            walletManager.newWalletAddress(fileName, nAddress, dialogGetPasswordForAddAddresses.password)
            listAddresses.loadModel(walletManager.getAddresses(fileName))
        }
    }

    DialogGetPassword {
        id: dialogGetPassword

        anchors.centerIn: Overlay.overlay
        width: applicationWindow.width > 440 ? 440 - 40 : applicationWindow.width - 40
        height: applicationWindow.height > 320 ? 320 - 40 : applicationWindow.height - 40

        headerMessage: qsTr("<b>Warning:</b> for security reasons, it is not recommended to keep the wallets unencrypted. Caution is advised.")
        headerMessageColor: Material.color(Material.Red)
        focus: true
        modal: true
        onAccepted: {
            var isEncrypted = walletManager.decryptWallet(fileName, password)
            walletModel.editWallet(index, name, isEncrypted, sky, coinHours)
        }
    }

    DialogSetPassword {
        id: dialogSetPassword

        anchors.centerIn: Overlay.overlay
        width: applicationWindow.width > 540 ? 540 - 40 : applicationWindow.width - 40
        height: applicationWindow.height > implicitHeight + 40 ? implicitHeight : applicationWindow.height - 40

        headerMessage: qsTr("We suggest that you encrypt each one of your wallets with a password. If you forget your password, you can reset it with your seed. Make sure you have your seed saved somewhere safe before encrypting your wallet.")
        headerMessageColor: Material.color(Material.Red)
        focus: true
        modal: true

        onAccepted: {
            var isEncypted = walletManager.encryptWallet(fileName, password)
            walletModel.editWallet(index, name, isEncypted, sky, coinHours)
        }
    } // DialogSetPassword

    DialogEditWallet {
        id: dialogEditWallet
        anchors.centerIn: Overlay.overlay

        focus: true
        modal: true

        onAccepted: {
            var qwallet = walletManager.editWallet(fileName, name)
            walletModel.editWallet(index, qwallet.name, encryptionEnabled, qwallet.sky, qwallet.coinHours )
        }
    } // DialogEditWallet

    // Roles: address, addressSky, addressCoinHours
    // Use listModel.append( { "address": value, "addressSky": value, "addressCoinHours": value } )
    // Or implement the model in the backend (a more recommendable approach)
    AddressModel {

        id: listAddresses
        property Timer timer: Timer {
            id: addressModelTimer
            interval: 7000
            repeat: true
            running: true
            onTriggered: {
                listAddresses.updateModel(fileName);
            }
        }
    }
    Component.onCompleted: {
        listAddresses.updateModel(fileName);
    }
}
