
$(() => {
    $("#privatecheckbox").change(function() {
        if ($(this).prop('checked')) {
            $("#hiddengroupfield").show()
        } else {
            $("#hiddengroupfield").hide()
        }
    })
})
