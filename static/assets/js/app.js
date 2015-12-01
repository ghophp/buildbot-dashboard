$(function(){

    var GRIDSTER_HASHED_KEY = hashedUrl + 'gridster';

    var gridster = $(".gridster ul").gridster({
        widget_base_dimensions: [140, 140],
        widget_margins: [10, 10],
        autogrow_cols: true,
        serialize_params: function ($w, wgd) {
            return {
                id: $w.find('.widget-wrapper').attr('id'),
                col: wgd.col,
                row: wgd.row,
                size_x: wgd.size_x,
                size_y: wgd.size_y
            };
        },
        draggable: {
            stop: function() {
                if (localStorage) {
                    localStorage.setItem(GRIDSTER_HASHED_KEY, JSON.stringify(gridster.serialize()));
                }
            }
        },
        resize: {
            enabled: false
        }
    }).data('gridster');

    var widgets = {};
    var addWidget = function(builder) {
        var positions = null;

        if (localStorage && localStorage.getItem(GRIDSTER_HASHED_KEY)) {
            positions = $.parseJSON(localStorage.getItem(GRIDSTER_HASHED_KEY));
        }

        var html = '<li class="new"><div id="'+builder.id+'" class="widget-wrapper"></div></li>';
        var current = _.find(positions, function(i) {
            return builder.id == i['id'];
        });

        if (current) {
            gridster.add_widget(html, current.size_x, current.size_y, current.col, current.row);
        } else {
            gridster.add_widget(html, 2, 2);    
        }

        var widget = ReactDOM.render(
            React.createElement(BuildWidget, { builder: builder }),
            document.getElementById(builder.id));

        if (localStorage && localStorage.getItem(hashedUrl + builder.id)) {
            widget.updateBuilder($.parseJSON(localStorage.getItem(hashedUrl + builder.id)));
        }

        widgets[builder.id] = widget;
    };
    
    $.get("/builders", function(data) {
        _.each(_.keys(data), function(key, i){
            var builder = data[key]; 
            builder.id = key;
            
            addWidget(builder);
        });

        var loc = window.location, new_uri;
        if (loc.protocol === "https:") {
            new_uri = "wss:";
        } else {
            new_uri = "ws:";
        }
        new_uri += "//" + loc.host;
        new_uri += loc.pathname + "ws";

        ws = new ReconnectingWebSocket(new_uri);
        ws.onmessage = function(e) {
            var message = $.parseJSON(e.data);
            var decoded = null;
            var builder = null;

            if (message.text) {
                decoded = Base64.decode(message.text);
                if (decoded) {
                    builder = $.parseJSON(decoded);
                }
            }

            if (builder) {
                if (localStorage) {
                    localStorage.setItem(hashedUrl + builder.id, decoded);
                }

                if (widgets[builder.id]) {
                    widgets[builder.id].updateBuilder(builder);
                } else if (builder.id) {
                    addWidget(builder);
                }
            }
        };
    });
});