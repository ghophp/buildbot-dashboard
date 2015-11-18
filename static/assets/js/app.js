$(function(){

    // more sizes can be implemented in the future
    var currentSizeDimension = [140, 140];
    if (genericSize == 'small') {
        currentSizeDimension = [100, 100];
    }

    $(".gridster ul").gridster({
        widget_base_dimensions: currentSizeDimension,
        widget_margins: [10, 10],
        autogrow_cols: true,
        resize: {
            enabled: false
        }
    });

    var widgets = {};
    var gridster = $(".gridster ul").gridster().data('gridster');

    $.get("/builders", function(data) {
        _.each(_.keys(data), function(key, i){
            var builder = data[key]; builder.id = key;

            gridster.add_widget('<li class="new"><div id="'+key+'" class="widget-wrapper"></div></li>', 2, 2);

            var widget = ReactDOM.render(
                React.createElement(BuildWidget, { builder: builder }),
                document.getElementById(key));

            widgets[key] = widget;
        });

        var loc = window.location, new_uri;
        if (loc.protocol === "https:") {
            new_uri = "wss:";
        } else {
            new_uri = "ws:";
        }
        new_uri += "//" + loc.host;
        new_uri += loc.pathname + "ws";

        ws = new WebSocket(new_uri);
        ws.onmessage = function(e) {

            var builder = $.parseJSON(event.data);
            if (widgets[builder.id]) {
                widgets[builder.id].updateBuilder(builder);
            }
        };
    });
});