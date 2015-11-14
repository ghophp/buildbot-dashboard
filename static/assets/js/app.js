$(function(){

    $(".gridster ul").gridster({
        widget_base_dimensions: [140, 140],
        widget_margins: [10, 10],
        autogrow_cols: true,
        resize: {
            enabled: true
        }
    });

    var gridster = $(".gridster ul").gridster().data('gridster');

    $.get("/builders", function(data) {
        _.each(_.keys(data), function(key, i){
            var builder = data[key]; builder.id = key;

            gridster.add_widget('<li class="new"><div id="'+key+'" class="widget-wrapper"></div></li>', 2, 2);

            ReactDOM.render(
                React.createElement(BuildWidget, { builder: builder }),
                document.getElementById(key));
        });
    });
});