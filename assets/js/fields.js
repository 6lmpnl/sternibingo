require("expose-loader?$!expose-loader?jQuery!jquery");

$(() => {
    var cells = $(".bingofield td");
    var max_cells = 0;

    var maxCount = cells.toArray()
            .map(( obj) => obj.getAttribute("count"))
            .filter(o => o != "")
            .reduce((agg, o) => Math.max(agg, parseInt(o)), 0);

    $.each(cells, function (i, obj) {
        var val = obj.getAttribute("count");
        if (val == "") return;
        val = parseFloat(val);
        console.log();
        $(obj).css("background-color", "hsl\(0,100%," + (80 - (val/maxCount) * 30) +  "%)")
    });
    //console.log("MaxCount: " + maxCount);
});