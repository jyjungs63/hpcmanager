var objText=""
var jobid;
let hpcserver = "http://10.15.20.111:8001/";

showHPCStatus = () => {
    // callAjax("GET","http://10.15.20.111:8000/ajaxStatus", dat)
    webix.ajax().get(hpcserver + "ajaxStatus")
    .then(function(data) {
        $$("idNodetable").clearAll();
        $$("idNodetable").parse(data.text())
    })
    .fail(function(xhr) {
        console.error("Error:", xhr.responseText);
    });
}

var accordion1 = {
  rows: [ 
    { 
      cols: [
              {
                view: "button", value: "Refresh Job Status ", css: "webix_primary", margin: 30,
                  click: function() {
                    showHPCStatus();
                  }
              },
              { view: "button", value: "Job Cancel", css: "webix_danger", 
                  click: function() {
                    var tableValues = $$("idNodetable").data.pull;
                    console.log("Table values:", tableValues);
                  }
              },    
            ],
    },
    {
      id: "idNodetable",
      view: "datatable",
      select:"row",
      multiselect:true,
      scrollX: true,
      height: 300,
      columns: [
        {id:"checked", header:[{content:"masterCheckbox"}], template:"{common.checkbox()}", width:30 },
        {id:"Jobid", header:"ID",width:50},
        {id:"Name", header:"NAME",width:100},
        {id:"User", header:"USER",width:100},
        {id:"St", header:"ST",width:50},
        {id:"Time", header:"TIME",width:70},
        {id:"Nodes", header:"NODES",width:50},
        {id:"Nodelist", header:"NODELIST",width:150},
      ],
      scheme: {
        $init: function (obj) {
            obj.checked = obj.id % 2 ? 1 : 0;
        }
      },
      on : {
        onItemClick: function(id) {
          //$$("Jobid").disable();
          this.checked = 1;
          this.refresh();
          //this.editColumn("target");
          },
        onCheck: function (rowId, colId, state) {
          if (state === 1) { // If checkbox is checked
            var rowNumber = this.getIndexById(rowId); // Get the row number
            var rowData = this.getItem(rowId);
            this.select(rowId,true)
            // alert("Row data:", rowData);
            // alert("Row number:", rowNumber);
        }
        }
      },
      url: {
        $proxy:true,
        load: function (view, params) {
          showHPCStatus();
        }
      }  
    },

  ],
};

function logEvent(type, message, args){
  webix.message({ text:message, expire:2500 });
}

var form2 = {
  css: "item3",
  view: "form",
  elements: [
    {
      rows: [
        { view: "text", label: "Job name", value: "", id: "idJobName" ,on: {
          onFocus: function (event) {
            objText = event.config.id
            //$$("idSimProject").setValue($$("fm").getState().selectedItem[0].id);
            //alert($$("fm").getState().selectedItem[0].id);
          }
        }
        },
        { view: "text", label: "Script", value: "", id: "idScript" ,on: {
          onFocus: function (event) {
            objText = event.config.id
            $$("idScript").setValue($$("fm").getState().selectedItem[0].id);
            //alert($$("fm").getState().selectedItem[0].id);
          }
        }},
        { view: "button", value: "Make Sim File", css: "webix_primary",
          click: function() {
            var dat = {
              "Jobname": $$("idJobName").getValue(),
              "Runscript": $$("idScript").getValue() // idScritp   idScript
            }
            callAjax("http://10.15.20.111:8001/ajaxRunStar", dat)
          }
        },
        { cols: [
          { view: "text", 
            label: "Status", 
            inputAlign: "right", 
            value: "OK", 
            css: "status_color_green", 
            readonly:true, 
            labelWidth: 300, },     
        ]
      }
      ],
    },
  ],
};

var form3 = {
  view: "form",
  elements: [
    {
      rows: [
        { view: "text", label: "Job name", value: "" , id: "idJob" ,on: {
          onFocus: function (event) {
            //objText = event.config.id
            //alert($$("fm").getState().selectedItem[0].id);
          }
        }},
        { view: "text", label: "Nodes", value: "6", id: "idNodes" ,on: {
          onFocus: function (event) {
            objText = event.config.id
            //$$("idScript").setValue($$("fm").getState().selectedItem[0].id);
            //alert($$("fm").getState().selectedItem[0].id);
          }
        }},
        { view: "text", label: "Script", value: "" , id: "idScript2" ,on: {
          onFocus: function (event) {
            objText = event.config.id
            let sct = $$("fm").getState().selectedItem[0].id;
            $$("idScript2").setValue(sct.replace("/",""));
          }
        }},
        { view: "text", label: "Sim File", value: "" , id: "idSimFile" ,on: {
          onFocus: function (event) {
            objText = event.config.id
            $$("idSimFile").setValue($$("fm").getState().selectedItem[0].id);
          }
        }},        
        { view: "button", value: "HPC RUN", css: "webix_primary",
          click: function() {
            var dat = {
              "Jobname":   $$("idJob").getValue(),
              "Runscript": $$("idScript2").getValue(),
              "Nodes":     $$("idNodes").getValue()
            }
            callAjax("PUT", "http://10.15.20.111:8000/ajaxRunStar", dat)
          }
        }, 
        { cols: [
          { view: "text", 
            label: "Status", 
            inputAlign: "right", 
            value: "Running", 
            css: "status_color_red", 
            readonly: true,
            labelWidth: 300,   
          },           
        ]
      }
      ],
    },
  ],
};

var form4 = {
  view: "form",
  elements: [
    {
      rows: [ 
        {
        cols: [
          { view: "button", value: "Get User", 
          click: function() {

            var data = {jobid: String(jobid)}

            $.ajax({
              url:  "http://10.15.10.148:9022/getUser",         // URL of the API endpoint
              type: "GET",         // HTTP method PUT With Data Get no Data
              //data: JSON.stringify(data),
              dataType: "json",      // Type of data expected from the server
              success: function (response) {
                 alert("HPC에서 Jobid " +  "작업 종료중 입니다")
              },
              error: function (xhr, status, error) {
                  // Error callback
                  console.error("Error:", error);
                  // Handle errors here
              }
            });
          }},
          { view: "button", value: "Free Surface half"},
          { view: "button", value: "Pressure Stern Proj."},
        ],
       },
       {
        cols: [
            { view: "button", value: "Pressure Surface Hull"},
            { view: "button", value: "Pressure Surface half"},
            { view: "button", value: "Pressure Stern Proj."},
          ],
       },
       {
        cols: [
          { view: "button", value: "Stream Stern"},
          { view: "button", value: "Stream Stern"},
          { view: "button", value: "Wake Contuor"},
        ]
       }
      ],
    },
  ],
};

let fileserver = new URLSearchParams(location.search)
let fileurl = fileserver.get("server")
if (fileurl == undefined)
  fileurl = "http://10.15.20.111:3200/";

  let data = {user: "leadship", fileurl: fileurl}
  //let data = {user: "leadship", fileurl: fileurl}

  $.ajax({
    url:  "http://10.15.20.111:8001/ajaxrunFileserver",         // URL of the API endpoint
    type: "PUT",         // HTTP method PUT With Data Get no Data
    data: JSON.stringify(data),
    dataType: "json",      // Type of data expected from the server
    async: false,
    success: function (response) {
      let job = JSON.parse(response) ;
      jobid = response
      webixready();
    },
    error: function (xhr, status, error) {
        // Error callback
        console.error("Error:", error);
        // Handle errors here
    }
  });

function webixready () {
  webix.ready(function () {
    webix.ui({
      cols: [
        {
          gravity: 0.7,
          css: "parentView",
          rows: [
            {
              view: "template",
              data: {
                columns: "home",
                src: "lib/images/icon/home.png",
                class: "home",
              },
              template: function (obj) {
                return ('<a id="idManager" href="#" <span><img src="' + obj.src + '" class="' + obj.class + '" /> </span> </a> ');
              },
              height: 29,
              css: "root",
            },
            {
              id: "fm",
              view: "filemanager",
              url: "http://10.15.20.111:3200/",
              //url: "https://docs.webix.com/filemanager-backend/",
              //mode: "tables",
              //override: new Map([[fileManager.views.topbar, CustomTopBar]]),
              preview: {
                  active: true,
                  type: "jpg",
                  width: 400,
                  height: 400
              },
              on: {
                onInit: app => {
                  const state = app.getState();
                  const evs = app.getState().$changes;
                  //alert ($$("fm").getState().path)
                  if (  $$("idSimProject") != undefined)
                    state.$observe("idSimProject", v => $$("idSimProject").setValue(v));
                  //state.$observe("selectedItem", v => {
                    //$$("idSimProject").setValue(v.map(a => a.id).join(", "))
                   //});
                  //if (objText == "idParam")
                  //state.$observe("selectedItem", v => {
                   // $$("idSimProject").setValues({ idSimProject: v.map(a => a.id).join(", ") })
                  //});
                },
              },
            },
          ],
        },
        {
          id: "id1",
          gravity: 0.3,
          view: "scrollview",
          scroll: "y",
          body: {
            rows: [
              {
                rows: [
                  {
                    cols: [
                      {
                        rows: [
                          {
                            rows: [
                              {
                                multi: true,
                                view: "accordion",
                                height: 0,
                                css: "acc_menulist",
                                rows: [
                                  {
                                    view: "accordionitem",
                                    header: "HPC Monitoring",
                                    body: accordion1,
                                  },
                                  {
                                    view: "accordionitem",
                                    header: "Make SIM File",
                                    body: form2,
                                  },
                                  {
                                    view: "accordionitem",
                                    header: "Star-CCM+ Run",
                                    body: form3,
                                  },
                                  {
                                    view: "button",
                                    css: "webix_primary",
                                    value: "Stop Filemanager",
                                    click: function() {

                                      var data = {jobid: String(jobid)}
                                      $.ajax({
                                        url:  "http://10.15.20.111:8001/ajaxStopfileserver",         // URL of the API endpoint
                                        type: "PUT",         // HTTP method PUT With Data Get no Data
                                        data: JSON.stringify(data),
                                        dataType: "json",      // Type of data expected from the server
                                        success: function (response) {
                                           alert("HPC에서 Jobid " +  "작업 종료중 입니다")
                                        },
                                        error: function (xhr, status, error) {
                                            // Error callback
                                            console.error("Error:", error);
                                            // Handle errors here
                                        }
                                      });
                                    }


                                  },
                                  {
                                    view: "accordionitem",
                                    header: "OutPut Display",
                                    body: form4,
                                  },
                                ],
                              },
                            ],
                          },
                        ],
                      },
                    ],
                  },
                ],
              },
            ],
          },
        },
      ],
    });
  });
} 

function callAjax(fncType, fucname, data) {
  $.ajax({
    url:  fucname,         // URL of the API endpoint
    type: fncType,         // HTTP method PUT With Data Get no Data
    data: JSON.stringify(data),
    dataType: "json",      // Type of data expected from the server
    success: function (response) {
        if ( fucname !="http://10.15.20.111:8000/ajaxRunStar" )
          $$("idNodetable").parse(response)
        else
          alert("HPC에서 Jobid " + response[0]['Jobid'] + "실행중입니다")
    },
    error: function (xhr, status, error) {
        // Error callback
        console.error("Error:", error);
        // Handle errors here
    }
  });
}

window.onbeforeunload = function() {


  

}