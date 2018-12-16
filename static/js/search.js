$(document).ready(function() {
    function load_data(query) {
        $.ajax({
            url: "/search",
            method: "POST",
            dataType: "html",
            data: { username: query },
            success: function(message) {
                // Extract data to table.
                var new_tbody = document.createElement("tbody");
                new_tbody.id = "result-data";
                user_list = JSON.parse(message);
                for (var i in user_list) {
                    var tr = document.createElement("tr");

                    var td = document.createElement("td");
                    td.innerHTML = i.toString();
                    tr.appendChild(td);

                    td = document.createElement("td");
                    td.innerHTML = user_list[i].username;
                    tr.appendChild(td);

                    td = document.createElement("td");
                    td.innerHTML = user_list[i].email;
                    tr.appendChild(td);

                    new_tbody.appendChild(tr);
                }

                var old_tbody = document.getElementById("result-data");
                old_tbody.parentNode.replaceChild(new_tbody, old_tbody);
            },
            error: function(jqXHR, textStatus, errorThrown) {
                console.log("Error while perform '/search' request");
            }
        });
    }

    $('#search-box').keyup(function(){
        var search = $(this).val();
        if(search != ''){
            load_data(search);
        } else{
            load_data();
        }
    });
});