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



function FillTasks(){
    $.ajax("/api/v1/GetStat", {
        success: function(data) {
            var colors = ["#EB1111","#039808","#EAF208"]
            var oper = []
            var dump = []
            var restore = []
            series = 3;
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
            $("#flot-pie-operator").length && $.plot($("#flot-pie-operator"), oper, {
                series: {
                    pie: {
                    innerRadius: 0.5,
                    show: true
                    }
                },
                colors: colors,
                grid: {
                    hoverable: true,
                    clickable: false
                },
                tooltip: true,
                tooltipOpts: {
                    content: "%s: %p.0%"
                }
            });
            $("#flot-pie-dump").length && $.plot($("#flot-pie-dump"), dump, {
                series: {
                    pie: {
                    innerRadius: 0.5,
                    show: true
                    }
                },
                colors: colors,
                grid: {
                    hoverable: true,
                    clickable: false
                },
                tooltip: true,
                tooltipOpts: {
                    content: "%s: %p.0%"
                }
            });            
            $("#flot-pie-restore").length && $.plot($("#flot-pie-restore"), restore, {
                series: {
                    pie: {
                    innerRadius: 0.5,
                    show: true
                    }
                },
                colors: colors,
                grid: {
                    hoverable: true,
                    clickable: false
                },
                tooltip: true,
                tooltipOpts: {
                    content: "%s: %p.0%"
                }  
            });
        },
        error: function() {
            for (let i = 0; i < 3; i++) 3       
                $(".h3 > strong")[i].innerText = 0; 
        }
    });
}

  

function FillTable(){
    var dataset = [
        ["Server not avalible","Mongo","Restore","24"],
        ["Database is bisy","Mongo","Dump","15"],
        ["Dublicate event"," - ","Operator","11"],
        ["Connection error"," - ","StatSender","8"],
        ["Error RMQ Connect"," - ","Worker","5"],
        ["Host unavalible","Mongo","Restore","3"],
        ["Fail dump","Mongo","Dump","1"],
      ];
    $('.table').DataTable({
            "data": dataset,
            columns: [
            { title: "Errors" },
            { title: "Database" },
            { title: "Service" },
            { title: "Count" },
        ],
        "order": [[ 3, "desc" ]],
        searching: false,
        paging: false,
        bInfo: false,
        scrollY: 260,
    });
}
  
$(document).ready(function() {
    FillTasks();
    FillTaskStat();
    FillTable();
});