function FillTaskStat(){
    $.ajax("/api/v1/GetTasksStatus", {
        success: function(data) {
            $(".h3 > strong")[0].innerText = data.Count; 
            $(".h3 > strong")[1].innerText = data.Active; 
            $(".h3 > strong")[2].innerText = data.Inactive;
            $(".h3 > strong")[3].innerText = 0;
        },
        error: function() {
            for (let i = 0; i < 3; i++) 3       
                $(".h3 > strong")[i].innerText = 0; 
        }
     });
}
function labelFormatter(label, series) {
    return "<div style='font-size:8pt; text-align:center; padding:2px; color:white;'>" + label + "<br/>" +Math.round(series.percent) + "% (" +series.data[0][1] +")</div>";
}

function FillTasks(){
    $.ajax("/api/v1/GetStat", {
        success: function(data) {
            var colors = ["#EB1111","#039808"]
            var oper = []
            var dump = []
            var restore = []
            var empty = [{"label":"NoData", "data": 1}]
            series = 2;
            var pieSettings = {
                series: {
                    pie: {
                        radius: 1,
                        label: {
                            show: true,
                            radius: 0.5,
                            formatter: labelFormatter,
                        },
                        show: true,
                    }
                },
                colors: colors,
                legend: {show: false}
            }
            var emptySettings = Object.assign({}, pieSettings);
            emptySettings.colors = ["#d3dfde"]
            Object.keys(data.Operator).forEach(function(key, i ){
                oper[i] = {
                    label: key,
                    data: data.Operator[key]
                }
            });
            Object.keys(data.Dump).forEach(function(key, i ){
                dump[i] = {
                    label: key,
                    data: data.Dump[key]
                }
            });
            Object.keys(data.Restore).forEach(function(key, i ){
                restore[i] = {
                    label: key,
                    data: data.Restore[key]
                }
            });
            
            operator = $("#flot-pie-operator").length && $.plot(
                $("#flot-pie-operator"),oper, pieSettings
            );
            
            if (isNaN(operator.getData()[0].percent)){
                $.plot($("#flot-pie-operator"), empty, emptySettings);
            }

            dumper = $("#flot-pie-dump").length && $.plot(
                $("#flot-pie-dump"), dump, pieSettings
            );
            if (isNaN(dumper.getData()[0].percent)){
                $.plot($("#flot-pie-dump"), empty, emptySettings);
            }
            
            restorer = $("#flot-pie-restore").length && $.plot(
                $("#flot-pie-restore"), restore, pieSettings
            );
            if (isNaN(restorer.getData()[0].percent)){
                $.plot($("#flot-pie-restore"), empty, emptySettings);
            }
        },
    });
}

  

function FillTable(){
    $.ajax("/api/v1/GetErrors", {
        success: function(data) {
            $('.table').DataTable({
                    "data": data,
                    columns: [
                    { title: "Error" },
                    { title: "Service" },
                    { title: "Count" },
                ],
                "order": [[ 2, "desc" ]],
                searching: false,
                paging: false,
                bInfo: false,
                scrollY: 275,
            });
        }
    });
}
  
$(document).ready(function() {
    FillTasks();
    FillTaskStat();
    FillTable();
});