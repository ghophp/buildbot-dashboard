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
            var builder = data[key];
            gridster.add_widget('<li class="new">'+key+'</li>', 2, 2);
        });
    });
});

//<li data-row="1" data-col="1" data-sizex="1" data-sizey="1"></li>