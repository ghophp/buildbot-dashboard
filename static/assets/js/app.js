$(function(){

    var GRIDSTER_HASHED_KEY = hashedUrl + 'gridster';

    // more sizes can be implemented in the future
    var currentSizeDimension = [140, 140];
    if (genericSize == 'small') {
        currentSizeDimension = [100, 100];
    }

    var gridster = $(".gridster ul").gridster({
        widget_base_dimensions: currentSizeDimension,
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
    
    $.get("/builders", function(data) {
        _.each(_.keys(data), function(key, i){
            var builder = data[key]; builder.id = key;
            if (!displayEmptyBuilder && (!builder.cachedBuilds || builder.cachedBuilds.length <= 0)) {
                return;
            }

            var positions = null;

            if (localStorage && localStorage.getItem(GRIDSTER_HASHED_KEY)) {
                positions = $.parseJSON(localStorage.getItem(GRIDSTER_HASHED_KEY));
            }

            var html = '<li class="new"><div id="'+key+'" class="widget-wrapper"></div></li>';
            var current = _.find(positions, function(i) {
                return key == i['id'];
            });

            if (current) {
                gridster.add_widget(html, current.size_x, current.size_y, current.col, current.row);
            } else {
                gridster.add_widget(html, 2, 2);    
            }

            var widget = ReactDOM.render(
                React.createElement(BuildWidget, { builder: builder }),
                document.getElementById(key));

            if (localStorage && localStorage.getItem(hashedUrl + key)) {
                widget.updateBuilder($.parseJSON(localStorage.getItem(hashedUrl + key)));
            }

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
            
            if (localStorage) {
                localStorage.setItem(hashedUrl + builder.id, event.data);
            }

            if (widgets[builder.id]) {
                widgets[builder.id].updateBuilder(builder);
            }
        };
    });
});