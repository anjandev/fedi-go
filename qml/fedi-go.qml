import QtQuick 2.0

Item {

    Rectangle {
        id: rectangle1
        x: 0
        y: 0
        width: 640
        height: 480
        color: "#282c37"
    }

    TextInput {
        id: instanceURL
        x: 8
        y: 196
        width: 624
        height: 89
        color: "#9baec8"
        text: qsTr("Click and enter your Instance's URL")
        horizontalAlignment: Text.AlignHCenter
        font.pixelSize: 32
    }

    Item {
        id: item1
        x: 351
        y: 420
        width: 281
        height: 52

        Rectangle {
            id: rectangle
            width: 281
            height: 44
            color: "#9baec8"
            radius: 6

            TextEdit {
                id: textEdit
                x: 0
                y: -1
                width: 278
                height: 32
                color: "#282c37"
                text: qsTr("Authenticate Me")
                anchors.verticalCenter: parent.verticalCenter
                anchors.horizontalCenter: parent.horizontalCenter
                scale: 1
                transformOrigin: Item.Center
                horizontalAlignment: Text.AlignHCenter
                font.pixelSize: 26
            }
        }
    }


}

/*##^## Designer {
    D{i:0;autoSize:true;height:480;width:640}D{i:11;anchors_height:44;anchors_width:278;anchors_x:0;anchors_y:-1}
}
 ##^##*/
