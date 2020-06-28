new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        avatar: null, // avatar address used for grabbing an avatar
        username: null, // Our username
        joined: false // True if avatar and username have been filled in
    },
    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            self.chatContent += '<div class="chip">'
                    + '<img src="' + self.avatarURL(msg.avatar) + '">' // Avatar
                    + msg.username
                + '</div>'
                + emojione.toImage(msg.message) + '<br/>'; // Parse emojis

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        });
    },
    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        avatar: this.avatar,
                        username: this.username,
                        message: $('<p>').html(this.newMsg).text() // Strip out html
                    }
                ));
                this.newMsg = ''; // Reset newMsg
            }
        },
        join: function () {
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            };
            this.ws.send(
                JSON.stringify({
                        avatar: this.avatar,
                        username: this.username,
                        message: ""
                })
            );
            this.avatar = $('<p>').html(this.avatar).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
        },   
        avatarURL: function(avatar) {
             if (avatar.match(/\.(jpeg|jpg|gif|png)$/) != null) {
                return avatar;
            } else {
                return "avatar.jpg"
            }
        }
    }
});