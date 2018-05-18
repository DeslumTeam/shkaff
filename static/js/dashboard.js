function FillTasks(){
    var da1 = [],
    series = 3;
    chart_names = ['New', 'Success', 'Fail'];
    for (var i = 0; i < series; i++) {
    da1[i] = {
        label: chart_names[i],
        data: Math.floor(Math.random() * 100) + 1
    }
    }
    var colors = ["#039808","#EAF208","#EB1111"]
    $("#flot-pie-operator").length && $.plot($("#flot-pie-operator"), da1, {
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

    $("#flot-pie-dump").length && $.plot($("#flot-pie-dump"), da1, {
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

    $("#flot-pie-restore").length && $.plot($("#flot-pie-restore"), da1, {
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
}

  function FillTaskStat(){
    a = {"Count": Math.floor(Math.random() * 100) + 1,
         "Active":Math.floor(Math.random() * 100) + 1,
         "Inactive":Math.floor(Math.random() * 100) + 1,
         "Storage":Math.floor(Math.random() * 100) + 1 +"Gb"
        }
    $(".h3 > strong")[0].innerText = a.Count; 
    $(".h3 > strong")[1].innerText = a.Active; 
    $(".h3 > strong")[2].innerText = a.Inactive;
    $(".h3 > strong")[3].innerText = a.Storage; 
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