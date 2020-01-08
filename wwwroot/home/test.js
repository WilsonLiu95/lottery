$(function () {
    var connection = $.connection('/echo');

    connection.received(function (data) {
        $('#messages').append('<li>' + JSON.parse(data)['Type'] + ":" + data + '</li>');
        Lottery.factory(JSON.parse(data));
    });

    connection.start();
     
});