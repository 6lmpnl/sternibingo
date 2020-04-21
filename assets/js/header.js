require("expose-loader?$!expose-loader?jQuery!jquery");

function isInClass(original, current, classname) {
    for (var i = original; i != current && i !== undefined; i = i.parentElement) {
        if (i.classList.contains(classname))
            return true;
    }
    return false;
}

$(() => {
    var bttns = $("header .bttn");

    bttns.on("click", function (e) {
        if (!isInClass(e.originalEvent.originalTarget, e.currentTarget, "wndw")) {
            var cl = e.currentTarget.classList;
            if (cl.contains("open")) {
                cl.remove("open");
            } else {
                for (var i = 0; i < bttns.length; i++) {
                    bttns[i].classList.remove("open");
                }
                cl.add("open");
            }
        }
    });

    $("body").on("click", function (e) {
        if (!isInClass(e.originalEvent.target, e.currentTarget, "bttn")
         || e.originalEvent.target.classList.contains("wndw")
         || e.originalEvent.target.parentElement.classList.contains("wndw")) {
            for (var i = 0; i < bttns.length; i++) {
                bttns[i].classList.remove("open");
            }
        }
    });
})
